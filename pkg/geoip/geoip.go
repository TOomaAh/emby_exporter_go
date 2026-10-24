package geoip

import (
	"TOomaAh/emby_exporter_go/internal/conf"
	"TOomaAh/emby_exporter_go/pkg/logger"
	"errors"
	"fmt"
	"os"
	"sync"
)

var (
	// ErrDatabaseNotFound is returned when the GeoIP database file cannot be found.
	ErrDatabaseNotFound = errors.New("GeoIP database file not found")

	// ErrDatabaseNotInitialized is returned when the GeoIP database has not been initialized.
	ErrDatabaseNotInitialized = errors.New("GeoIP database is not initialized")

	// ErrGeoIPDisabled is returned when GeoIP functionality is disabled in the configuration.
	ErrGeoIPDisabled = errors.New("GeoIP functionality is disabled")

	// geoIPInstance holds the singleton instance of the GeoIP database.
	geoIPInstance GeoIP

	// geoIPMutex protects access to the geoIPInstance.
	geoIPMutex sync.RWMutex
)

// GeoIP defines the interface for geolocation services.
type GeoIP interface {
	// Lookup returns a Record containing all available geolocation information for the given IP.
	Lookup(ip string) (*Record, error)

	// Close closes the GeoIP database.
	Close() error
}

type Record struct {
	Latitude    float64
	Longitude   float64
	City        string
	Region      string
	CountryCode string
}

func LocalRecord() *Record {
	return &localRecord
}

func DisabledRecord() *Record {
	return &disabledRecord
}

// InitGeoIPDatabase initializes the GeoIP database based on configuration.
// It returns a GeoIP interface that can be used for geolocation lookups.
func InitGeoIPDatabase(cfg *conf.Config, l logger.Interface) (GeoIP, error) {
	geoIPMutex.Lock()
	defer geoIPMutex.Unlock()

	// If GeoIP is disabled in configuration, return a no-op implementation
	if !cfg.Options.GeoIP {
		l.Info("GeoIP functionality is disabled in configuration")
		noGeo := NewNoGeoIP()
		geoIPInstance = noGeo
		return noGeo, nil
	}

	// Determine the database file path from environment or default
	dbFile := os.Getenv("GEOIP_DB")
	if dbFile == "" {
		dbFile = "geoip.mmdb"
		l.Info("Using default GeoIP database path: %s", dbFile)
	} else {
		l.Info("Using GeoIP database path from environment: %s", dbFile)
	}

	// Check if the database file exists
	_, err := os.Stat(dbFile)
	fileExists := !os.IsNotExist(err)

	accountID := cfg.Options.GeoIPOptions.AccountId
	licenseKey := cfg.Options.GeoIPOptions.LicenseKey

	// If file doesn't exist and we have no credentials, we can't proceed with GeoIP functionality
	if !fileExists && (accountID == "" || licenseKey == "") {
		l.Error("GeoIP database file not found and no credentials provided to download it")
		noGeo := NewNoGeoIP()
		geoIPInstance = noGeo
		return noGeo, fmt.Errorf("%w: %s", ErrDatabaseNotFound, dbFile)
	}

	// Try to create a GeoIP manager
	geoIPManager, err := NewGeoIPManager(dbFile, accountID, licenseKey, l)
	if err != nil {
		l.Error("Failed to initialize GeoIP database: %s", err)
		noGeo := NewNoGeoIP()
		geoIPInstance = noGeo
		return noGeo, fmt.Errorf("failed to initialize GeoIP database: %w", err)
	}

	l.Info("GeoIP database initialized successfully")
	geoIPInstance = geoIPManager
	return geoIPManager, nil
}

// GetGeoIPDatabase returns the initialized GeoIP database instance.
// It returns an error if the database has not been initialized.
func GetGeoIPDatabase() (GeoIP, error) {
	geoIPMutex.RLock()
	defer geoIPMutex.RUnlock()

	if geoIPInstance == nil {
		return nil, ErrDatabaseNotInitialized
	}

	return geoIPInstance, nil
}
