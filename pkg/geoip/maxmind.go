package geoip

import (
	"TOomaAh/emby_exporter_go/pkg/logger"
	"TOomaAh/emby_exporter_go/pkg/request"
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/oschwald/geoip2-golang"
)

const (
	maxmindBaseURL      = "https://download.maxmind.com"
	maxmindPath         = "/geoip/databases/GeoLite2-City/download?suffix=tar.gz"
	maxmindShaPath      = "/geoip/databases/GeoLite2-City/download?suffix=tar.gz.sha256"
	httpClientTimeout   = 30 * time.Second // Increased timeout for download
	updateCheckInterval = 24 * time.Hour
	userAgent           = "emby_exporter_go"
	shaFileName         = "geoip_update.mmdb.sha256" // Consistent file name for the SHA256 checksum
)

var (
	ErrSHA256Empty      = errors.New("SHA256 checksum is empty")
	ErrFileDoesNotExist = errors.New("file does not exist")
	ErrCannotRemoveFile = errors.New("cannot remove file")
	ErrCannotRenameFile = errors.New("cannot rename file")
	ErrInvalidIPAddress = errors.New("invalid IP address")
)

// GeoIPManager manages the loading and updating of the GeoIP database.
type GeoIPManager struct {
	mu      sync.RWMutex // Changed to RWMutex for better concurrency
	db      *geoip2.Reader
	path    string
	updater GeoIPUpdater
	logger  logger.Interface
	ctx     context.Context
	cancel  context.CancelFunc
}

// GeoIPUpdater defines the interface for updating the GeoIP database.
type GeoIPUpdater interface {
	NeedUpdate() (bool, error)
	Update(geoIPManager *GeoIPManager) error
	ReadCurrentSHA() string
}

// NoAuthUpdater implements GeoIPUpdater for scenarios without authentication.
type NoAuthUpdater struct {
	Logger logger.Interface
}

// AuthUpdater implements GeoIPUpdater for authenticated updates.
type AuthUpdater struct {
	client     *MaxmindClient
	currentSHA string
	logger     logger.Interface
}

// MaxmindClient handles API communication with MaxMind.
type MaxmindClient struct {
	accountID      string
	licenceKey     string
	baseURL        string
	client         *http.Client
	requestManager *request.RequestManager
}

// NewGeoIPManager creates a new GeoIP manager with the specified database file and credentials.
func NewGeoIPManager(file, accountID, licenceKey string, logLevel string) (*GeoIPManager, error) {
	l := logger.New(logLevel)

	// Create context with cancellation for clean shutdown
	ctx, cancel := context.WithCancel(context.Background())

	geoIPManager := &GeoIPManager{
		path:   file,
		logger: l,
		ctx:    ctx,
		cancel: cancel,
	}

	// Create the directory for the GeoIP database if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory for GeoIP database: %w", err)
	}

	var updater GeoIPUpdater

	if accountID == "" || licenceKey == "" {
		l.Info("No authentication provided for GeoIP database updates")
		updater = &NoAuthUpdater{
			Logger: l,
		}
	} else {
		l.Info("Authentication provided for GeoIP database updates")
		client := NewMaxmindClient(licenceKey, accountID)
		updater = &AuthUpdater{
			client: client.(*MaxmindClient),
			logger: l,
		}

		// Initialize the current SHA
		updater.(*AuthUpdater).currentSHA = updater.ReadCurrentSHA()

		// Run initial update if needed
		runDatabaseUpdater(updater, geoIPManager, l)

		// Start update goroutine
		go func() {
			ticker := time.NewTicker(updateCheckInterval)
			defer ticker.Stop()

			l.Info("GeoIP database updater started")
			for {
				select {
				case <-ticker.C:
					runDatabaseUpdater(updater, geoIPManager, l)
				case <-ctx.Done():
					l.Info("GeoIP database updater stopped")
					return
				}
			}
		}()
	}

	geoIPManager.updater = updater

	// Try to open the database, creating it if it doesn't exist
	if _, err := os.Stat(file); os.IsNotExist(err) {
		l.Info("GeoIP database does not exist, creating...")
		if accountID != "" && licenceKey != "" {
			if err := updater.Update(geoIPManager); err != nil {
				l.Error("Failed to create initial GeoIP database: %s", err)
				// Continue, maybe a default database will be used
			}
		}
	}

	// Open the database
	db, err := geoip2.Open(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open GeoIP database: %w", err)
	}

	geoIPManager.db = db

	return geoIPManager, nil
}

// Close closes the GeoIP database and stops the updater.
func (g *GeoIPManager) Close() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Cancel the context to stop the updater goroutine
	if g.cancel != nil {
		g.cancel()
	}

	// Close the database
	if g.db != nil {
		return g.db.Close()
	}

	return nil
}

// GetCountryCode returns the ISO code for the country of the given IP.
func (g *GeoIPManager) GetCountryCode(ipStr string) (string, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", ErrInvalidIPAddress
	}

	g.mu.RLock()
	defer g.mu.RUnlock()

	record, err := g.db.Country(ip)
	if err != nil {
		g.logger.Error("Error getting country code: %s", err)
		return "", err
	}

	return record.Country.IsoCode, nil
}

// GetCountryName returns the English name for the country of the given IP.
func (g *GeoIPManager) GetCountryName(ipStr string) (string, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", ErrInvalidIPAddress
	}

	g.mu.RLock()
	defer g.mu.RUnlock()

	record, err := g.db.Country(ip)
	if err != nil {
		g.logger.Error("Error getting country name: %s", err)
		return "", err
	}

	return record.Country.Names["en"], nil
}

// GetCity returns the city name for the given IP.
func (g *GeoIPManager) GetCity(ipStr string) (string, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", ErrInvalidIPAddress
	}

	g.mu.RLock()
	defer g.mu.RUnlock()

	record, err := g.db.City(ip)
	if err != nil {
		g.logger.Error("Error getting city: %s", err)
		return "", err
	}

	return record.City.Names["en"], nil
}

// GetContinent returns the continent name for the given IP.
func (g *GeoIPManager) GetContinent(ipStr string) (string, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", ErrInvalidIPAddress
	}

	g.mu.RLock()
	defer g.mu.RUnlock()

	record, err := g.db.City(ip)
	if err != nil {
		g.logger.Error("Error getting continent: %s", err)
		return "", err
	}

	return record.Continent.Names["en"], nil
}

// GetLocation returns the latitude and longitude for the given IP.
func (g *GeoIPManager) GetLocation(ipStr string) (float64, float64, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return 0, 0, ErrInvalidIPAddress
	}

	g.mu.RLock()
	defer g.mu.RUnlock()

	record, err := g.db.City(ip)
	if err != nil {
		g.logger.Error("Error getting location: %s", err)
		return 0, 0, err
	}

	return record.Location.Latitude, record.Location.Longitude, nil
}

// GetPostalCode returns the postal code for the given IP.
func (g *GeoIPManager) GetPostalCode(ipStr string) (string, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", ErrInvalidIPAddress
	}

	g.mu.RLock()
	defer g.mu.RUnlock()

	record, err := g.db.City(ip)
	if err != nil {
		g.logger.Error("Error getting postal code: %s", err)
		return "", err
	}

	return record.Postal.Code, nil
}

// GetRegion returns the region name for the given IP.
func (g *GeoIPManager) GetRegion(ipStr string) (string, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", ErrInvalidIPAddress
	}

	g.mu.RLock()
	defer g.mu.RUnlock()

	record, err := g.db.City(ip)
	if err != nil {
		g.logger.Error("Error getting region: %s", err)
		return "", err
	}

	if len(record.Subdivisions) == 0 {
		return "", nil
	}
	return record.Subdivisions[0].Names["en"], nil
}

// NewMaxmindClient creates a new client for MaxMind API.
func NewMaxmindClient(licenceKey, accountID string) request.Client {
	m := &MaxmindClient{
		baseURL:    maxmindBaseURL,
		accountID:  accountID,
		licenceKey: licenceKey,
		client: &http.Client{
			Timeout: httpClientTimeout,
		},
	}

	m.requestManager = request.NewRequestManager(m)
	return m
}

// ApplyAuthentication adds authentication to the request.
func (m *MaxmindClient) ApplyAuthentication(r *http.Request) error {
	r.SetBasicAuth(m.accountID, m.licenceKey)
	return nil
}

// SetHeaders sets the headers for the request.
func (m *MaxmindClient) SetHeaders(headers http.Header) {
	headers.Set("User-Agent", userAgent)
}

// GetBaseURL returns the base URL for the MaxMind API.
func (m *MaxmindClient) GetBaseURL() *url.URL {
	u, _ := url.Parse(m.baseURL)
	return u
}

// GetClient returns the HTTP client.
func (m *MaxmindClient) GetClient() *http.Client {
	return m.client
}

// NeedUpdate checks if the GeoIP database needs to be updated.
func (n *NoAuthUpdater) NeedUpdate() (bool, error) {
	n.Logger.Info("No authentication provided for GeoIP database updates, skipping update check")
	return false, nil
}

// Update updates the GeoIP database (no-op for NoAuthUpdater).
func (n *NoAuthUpdater) Update(geoIPManager *GeoIPManager) error {
	n.Logger.Info("No authentication provided for GeoIP database updates, skipping update")
	return nil
}

// ReadCurrentSHA returns an empty string for NoAuthUpdater.
func (n *NoAuthUpdater) ReadCurrentSHA() string {
	return ""
}

// ReadCurrentSHA reads the current SHA256 checksum of the GeoIP database.
func (m *AuthUpdater) ReadCurrentSHA() string {
	// First check if we have already cached the SHA
	if m.currentSHA != "" {
		return m.currentSHA
	}

	shaFile := shaFileName
	if _, err := os.Stat(shaFile); os.IsNotExist(err) {
		return ""
	}

	currentSha256, err := os.ReadFile(shaFile)
	if err != nil {
		m.logger.Error("Error reading SHA file: %s", err)
		return ""
	}

	// Cache the SHA value
	m.currentSHA = string(currentSha256)

	return m.currentSHA
}

// NeedUpdate checks if the GeoIP database needs to be updated.
func (m *AuthUpdater) NeedUpdate() (bool, error) {
	req, err := m.client.requestManager.NewRequest(http.MethodGet, maxmindShaPath, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	var sha256 = new(bytes.Buffer)
	err = m.client.requestManager.DoFile(req, sha256)
	if err != nil {
		return false, fmt.Errorf("failed to download SHA256: %w", err)
	}

	// Parse the SHA256 from the response
	sha256Content := sha256.String()
	// Switched to Debug level
	m.logger.Debug("Downloaded SHA256: '%s'", sha256Content)

	// Extract just the hash part (first field, ignoring filename)
	parts := strings.Fields(sha256Content)
	if len(parts) == 0 {
		return false, ErrSHA256Empty
	}
	sha256str := parts[0]

	if sha256str == "" {
		return false, ErrSHA256Empty
	}

	currentSHA256 := m.ReadCurrentSHA()
	// Switched to Debug level
	m.logger.Debug("Current SHA256: '%s'", currentSHA256)

	if currentSHA256 == "" {
		return true, nil
	}

	// Trim any whitespace or newlines for accurate comparison
	sha256str = strings.TrimSpace(sha256str)
	currentSHA256 = strings.TrimSpace(currentSHA256)

	// Extract just the hash part from current SHA (in case it contains filename)
	currentParts := strings.Fields(currentSHA256)
	if len(currentParts) > 0 {
		currentSHA256 = currentParts[0]
	}

	// Compare the current and new SHA256
	needsUpdate := sha256str != currentSHA256
	if needsUpdate {
		m.logger.Info("GeoIP database update available")
		// Switched to Debug level
		m.logger.Debug("SHA256 changed: '%s' vs '%s'", sha256str, currentSHA256)
	} else {
		// Switched to Debug level
		m.logger.Debug("GeoIP database is up to date (SHA256: '%s')", sha256str)
	}

	return needsUpdate, nil
}

// Update updates the GeoIP database.
func (m *AuthUpdater) Update(geoIPManager *GeoIPManager) error {
	m.logger.Info("Updating GeoIP database...")

	// Create a request to download the database
	req, err := m.client.requestManager.NewRequest(http.MethodGet, maxmindPath, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Create temporary files
	tmpDir, err := os.MkdirTemp("", "geoip_update")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}

	// We'll manually clean up the temp directory at the end to ensure all file handles are closed first
	// This helps avoid Windows file locking issues

	tarGzPath := filepath.Join(tmpDir, "geoip.tar.gz")
	mmdbPath := filepath.Join(tmpDir, "GeoLite2-City.mmdb")
	// No longer need shaPath as we handle the SHA file differently

	// Download the compressed database
	out, err := os.Create(tarGzPath)
	if err != nil {
		return fmt.Errorf("failed to create tarball file: %w", err)
	}
	defer out.Close()

	err = m.client.requestManager.DoFile(req, out)
	if err != nil {
		return fmt.Errorf("failed to download database: %w", err)
	}

	// Reset file pointer to beginning of file
	if _, err := out.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek in tarball: %w", err)
	}

	// Extract the database
	if err := decompressTarGz(out, mmdbPath, m.logger); err != nil {
		return fmt.Errorf("failed to decompress database: %w", err)
	}

	// Verify the extracted file exists
	if _, err := os.Stat(mmdbPath); os.IsNotExist(err) {
		// Clean up temp directory since we're returning early
		os.RemoveAll(tmpDir)
		return ErrFileDoesNotExist
	}

	// Get the SHA256 of the downloaded file - we'll handle this after successful DB update
	// Instead of trying to save the file immediately, we'll download it again later

	// Acquire exclusive lock for database update
	geoIPManager.mu.Lock()
	defer geoIPManager.mu.Unlock()

	// Close the current database if it exists
	if geoIPManager.db != nil {
		if err := geoIPManager.db.Close(); err != nil {
			m.logger.Error("Failed to close database: %s", err)
			// Continue with update anyway
		}
		geoIPManager.db = nil
	}

	// Replace the database file
	if err := os.Rename(mmdbPath, geoIPManager.path); err != nil {
		// Try to reopen the old database
		db, reopenErr := geoip2.Open(geoIPManager.path)
		if reopenErr == nil {
			geoIPManager.db = db
		}
		return fmt.Errorf("failed to replace database file: %w", err)
	}

	// Save the SHA256 - no shaFile to close as we're taking a different approach now

	// Get the SHA content we just downloaded to save it properly
	req, err = m.client.requestManager.NewRequest(http.MethodGet, maxmindShaPath, nil)
	if err != nil {
		m.logger.Error("Failed to create SHA request: %s", err)
	} else {
		var shaBuffer = new(bytes.Buffer)
		if err := m.client.requestManager.DoFile(req, shaBuffer); err != nil {
			m.logger.Error("Failed to download SHA: %s", err)
		} else {
			// Parse the SHA256 properly to store only the hash part
			shaContent := shaBuffer.String()
			parts := strings.Fields(shaContent)
			if len(parts) > 0 {
				shaHash := parts[0]

				// Remove existing SHA file if it exists
				shaFilePath := shaFileName
				if _, err := os.Stat(shaFilePath); err == nil {
					if err := os.Remove(shaFilePath); err != nil {
						m.logger.Error("Failed to remove existing SHA file: %s", err)
					}
				}

				// Write only the hash part to the SHA file
				if err := os.WriteFile(shaFilePath, []byte(shaHash), 0644); err != nil {
					m.logger.Error("Failed to save SHA file: %s", err)
				}
			} else {
				m.logger.Error("Failed to parse SHA content: %s", shaContent)
			}
		}
	}

	// Reopen the database
	db, err := geoip2.Open(geoIPManager.path)
	if err != nil {
		return fmt.Errorf("failed to open updated database: %w", err)
	}

	geoIPManager.db = db
	m.currentSHA = m.ReadCurrentSHA()

	m.logger.Info("GeoIP database updated successfully")

	// Clean up the temporary directory now that all operations are complete
	// This is placed at the end of the function to ensure all file handles are closed
	os.RemoveAll(tmpDir)

	return nil
}

// runDatabaseUpdater checks if an update is needed and performs it if necessary.
func runDatabaseUpdater(updater GeoIPUpdater, geoIPManager *GeoIPManager, l logger.Interface) {
	needsUpdate, err := updater.NeedUpdate()
	if err != nil {
		l.Error("Error checking for GeoIP database updates: %s", err)
		return
	}

	if needsUpdate {
		if err := updater.Update(geoIPManager); err != nil {
			l.Error("Error updating GeoIP database: %s", err)
		} else {
			l.Info("GeoIP database updated successfully")
		}
	}
}

// decompressTarGz extracts the GeoIP database from a tar.gz file.
func decompressTarGz(file *os.File, outputPath string, l logger.Interface) error {
	l.Info("Decompressing GeoIP database archive...")

	// Reset file pointer to beginning of file
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek in tarball: %w", err)
	}

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	fileFound := false
	var extractErr error

	for {
		header, err := tr.Next()
		switch {
		case err == io.EOF:
			if !fileFound {
				return fmt.Errorf("no .mmdb file found in archive")
			}
			return extractErr
		case err != nil:
			return fmt.Errorf("error reading tar: %w", err)
		case header == nil:
			continue
		}

		// Only extract .mmdb files
		if header.Typeflag == tar.TypeReg && filepath.Ext(header.Name) == ".mmdb" {
			fileFound = true
			l.Info("Found GeoIP database file in archive: %s", header.Name)

			outputDir := filepath.Dir(outputPath)
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				extractErr = fmt.Errorf("failed to create output directory: %w", err)
				continue
			}

			// We use a separate function to ensure the file is closed before returning
			extractErr = extractFile(tr, outputPath, header.Mode, l)
			if extractErr == nil {
				return nil
			}
		}
	}
}

// extractFile extracts a file from a tar reader to the specified output path.
// It ensures proper file handle closing in all cases.
func extractFile(tr *tar.Reader, outputPath string, mode int64, l logger.Interface) error {
	// Open output file
	out, err := os.OpenFile(outputPath, os.O_CREATE|os.O_RDWR, os.FileMode(mode))
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}

	// Always close the file when done
	defer func() {
		closeErr := out.Close()
		if closeErr != nil {
			l.Error("Failed to close output file: %s", closeErr)
		}
	}()

	// Copy content to output file
	n, err := io.Copy(out, tr)
	if err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}
	// Switched to Debug level
	l.Debug("Wrote %d bytes to %s", n, outputPath)

	// Explicitly sync the file to disk
	if err := out.Sync(); err != nil {
		return fmt.Errorf("failed to sync output file: %w", err)
	}

	return nil
}
