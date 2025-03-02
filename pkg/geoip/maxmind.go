package geoip

import (
	"TOomaAh/emby_exporter_go/pkg/logger"
	"TOomaAh/emby_exporter_go/pkg/request"
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
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

var (
	MAXMIND_PATH     = "/geoip/databases/GeoLite2-City/download?suffix=tar.gz"
	MAXMIND_SHA_PATH = "/geoip/databases/GeoLite2-City/download?suffix=tar.gz.sha256"
)

// GeoIPManager gère le chargement et la mise à jour de la base de données.
type GeoIPManager struct {
	mu      *sync.Mutex
	db      *geoip2.Reader
	path    string
	updater geoIpUpdater
	logger  logger.Interface
}

type geoIpUpdater interface {
	needUpdate() (bool, error)
	update(geoIpManager *GeoIPManager) error
}

type noAuthentUpdater struct{}

type authentUpdater struct {
	client *maxmindClient
}

type maxmindClient struct {
	accountId      string
	licenceKey     string
	baseUrl        string
	client         *http.Client
	requestManager *request.RequestManager
}

func newMaxmindClient(licenceKey, accountId string) request.Client {
	m := &maxmindClient{
		baseUrl:    "https://download.maxmind.com",
		accountId:  accountId,
		licenceKey: licenceKey,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	m.requestManager = request.NewRequestManager(m)
	return m
}

func (n *noAuthentUpdater) needUpdate() (bool, error) {
	return false, nil
}

func (n *noAuthentUpdater) update(geoIpManager *GeoIPManager) error {
	return nil
}

// NewGeoIPManager ouvre la base de données et stocke son timestamp de modification.
func newGeoIP(file string, accountId string, licenceKey string) (*GeoIPManager, error) {
	var updater geoIpUpdater
	mu := &sync.Mutex{}

	geopIpManager := &GeoIPManager{
		path:    file,
		updater: updater,
		mu:      mu,
	}

	if accountId == "" || licenceKey == "" {
		updater = &noAuthentUpdater{}
	} else {
		l := logger.New("info")

		updater = &authentUpdater{
			client: newMaxmindClient(licenceKey, accountId).(*maxmindClient),
		}

		runDatabaseUpdater(updater, geopIpManager, l)

		go func() {
			ticker := time.NewTicker(24 * time.Hour)

			defer ticker.Stop()
			l.Info("GeoIP database updater started")
			for range ticker.C {
				runDatabaseUpdater(updater, geopIpManager, l)
			}
		}()

	}

	db, err := geoip2.Open(file)
	if err != nil {
		return nil, err
	}

	geopIpManager.db = db

	return geopIpManager, nil
}

func runDatabaseUpdater(updater geoIpUpdater, geoIPManager *GeoIPManager, l logger.Interface) {
	b, err := updater.needUpdate()
	if err == nil && b {
		err := updater.update(geoIPManager)
		if err != nil {
			l.Error("Error while updating geoip database: %s", err)
		} else {
			l.Info("GeoIP database updated")
		}
	}
}

func (g *GeoIPManager) GetCountryCode(ip string) string {
	g.mu.Lock()
	defer g.mu.Unlock()
	record, err := g.db.Country(net.ParseIP(ip))

	if err != nil {
		g.logger.Error("Error while getting country code: %s", err)
		return ""
	}

	return record.Country.IsoCode
}

func (g *GeoIPManager) GetCountryName(ip string) string {
	g.mu.Lock()
	defer g.mu.Unlock()
	record, err := g.db.Country(net.ParseIP(ip))

	if err != nil {
		g.logger.Error("Error while getting country name: %s", err)
		return ""
	}

	return record.Country.Names["en"]
}

func (g *GeoIPManager) GetCity(ip string) string {
	g.mu.Lock()
	defer g.mu.Unlock()
	record, err := g.db.City(net.ParseIP(ip))

	if err != nil {
		g.logger.Error("Error while getting city: %s", err)
		return ""
	}

	return record.City.Names["en"]
}

func (g *GeoIPManager) GetContinent(ip string) string {
	g.mu.Lock()
	defer g.mu.Unlock()
	record, err := g.db.City(net.ParseIP(ip))

	if err != nil {
		g.logger.Error("Error while getting continent: %s", err)
		return ""
	}

	return record.Continent.Names["en"]
}

func (g *GeoIPManager) GetLocation(ip string) (float64, float64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	record, err := g.db.City(net.ParseIP(ip))

	if err != nil {
		g.logger.Error("Error while getting location: %s", err)
		return 0, 0
	}

	return record.Location.Latitude, record.Location.Longitude
}

func (g *GeoIPManager) GetPostalCode(ip string) string {
	g.mu.Lock()
	defer g.mu.Unlock()
	record, err := g.db.City(net.ParseIP(ip))

	if err != nil {
		g.logger.Error("Error while getting postal code: %s", err)
		return ""
	}

	return record.Postal.Code
}

func (g *GeoIPManager) GetRegion(ip string) string {
	g.mu.Lock()
	defer g.mu.Unlock()
	record, err := g.db.City(net.ParseIP(ip))

	if err != nil {
		g.logger.Error("Error while getting region: %s", err)
		return ""
	}

	if len(record.Subdivisions) == 0 {
		return ""
	}
	return record.Subdivisions[0].Names["en"]
}

func (g *GeoIPManager) Close() {
	g.db.Close()
}

func (m *maxmindClient) ApplyAuthentication(r *http.Request) error {
	// basic auth
	r.SetBasicAuth(m.accountId, m.licenceKey)
	return nil
}

func (m *maxmindClient) SetHeaders(headers http.Header) {
	headers.Set("User-Agent", "emby_exporter_go")
}

func (m *maxmindClient) GetBaseURL() *url.URL {
	u, _ := url.Parse(m.baseUrl)
	return u
}

func (m *maxmindClient) GetClient() *http.Client {
	return m.client
}

func (m *authentUpdater) needUpdate() (bool, error) {
	req, err := m.client.requestManager.NewRequest(http.MethodGet, MAXMIND_SHA_PATH, nil)

	if err != nil {
		return false, err
	}

	// string io writer
	var sha256 = new(bytes.Buffer)
	err = m.client.requestManager.DoFile(req, sha256)

	if err != nil {
		return false, err
	}

	// read sha256
	// split sha256
	sha256str := strings.Split(sha256.String(), " ")[0]

	// check if sha256 is the same
	if sha256str == "" {
		return false, errors.New("sha256 is empty")
	}

	var sha256File *os.File

	// check if sha256 file exists
	if _, err := os.Stat("geoip_update.mmdb.sha256"); os.IsNotExist(err) {
		f, err := os.Create("geoip_update.mmdb.sha256")

		if err != nil {
			return false, err
		}

		_, err = f.Write([]byte(sha256str))

		if err != nil {
			return false, err
		}
		return true, nil
	} else {
		sha256File, err = os.Open("geoip_update.mmdb.sha256")
		if err != nil {
			return false, err
		}
	}

	defer sha256File.Close()

	// read sha256 from file
	// split sha256
	currentSha256, err := os.ReadFile("geoip_update.mmdb.sha256")

	if err != nil {
		return false, err
	}

	if string(currentSha256) == "" {
		return true, nil
	}

	// read sha256 from file
	return !(sha256str == string(currentSha256)), nil

}

func (m *authentUpdater) update(geoIpManager *GeoIPManager) error {
	req, err := m.client.requestManager.NewRequest(http.MethodGet, MAXMIND_PATH, nil)
	l := logger.New("info")
	defer func() {
		err = os.Remove("geoip_update.mmdb.tar.gz")
		if err != nil {
			l.Error("Error while removing file: %s", err)
		}
		os.Remove("geoip_update.mmdb")

	}()

	if err != nil {
		return err
	}

	// create output file
	out, err := os.Create("geoip_update.mmdb.tar.gz")

	if err != nil {
		return err
	}

	err = m.client.requestManager.DoFile(req, out)

	if err != nil {
		return err
	}

	decompressTarGz()

	// remove tar.gz file
	err = os.Remove("geoip_update.mmdb.tar.gz")

	if _, err := os.Stat("geoip_update.mmdb"); os.IsNotExist(err) {
		return errors.New("file does not exist")
	}

	// remove actual file if exists
	if _, err := os.Stat(geoIpManager.path); !os.IsNotExist(err) {
		geoIpManager.mu.Lock()
		defer geoIpManager.mu.Unlock()

		if geoIpManager.db != nil {
			geoIpManager.Close()
		}

		err = os.Remove(geoIpManager.path)
		if err != nil {
			geoIpManager.db, _ = geoip2.Open(geoIpManager.path)
			return errors.New("cannot remove file")
		}
		err = os.Rename("geoip_update.mmdb", geoIpManager.path)
		if err != nil {
			geoIpManager.db, _ = geoip2.Open(geoIpManager.path)
			return errors.New("cannot remove file")
		}

		db, err := geoip2.Open(geoIpManager.path)
		if err != nil {
			return err
		}

		geoIpManager.db = db
	}

	// rename file

	return nil

}

func decompressTarGz() error {
	l := logger.New("info")

	// check if file exists
	if _, err := os.Stat("geoip_update.mmdb.tar.gz"); os.IsNotExist(err) {
		return errors.New("file does not exist")
	}

	// open file
	file, err := os.Open("geoip_update.mmdb.tar.gz")

	if err != nil {
		return err
	}

	gzr, err := gzip.NewReader(file)

	if err != nil {
		return err
	}

	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {

		header, err := tr.Next()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}

		l.Info("Decompressing %s", header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			continue
		case tar.TypeReg:
			if filepath.Ext(header.Name) == ".mmdb" {

				out, err := os.OpenFile("geoip_update.mmdb", os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
				if err != nil {
					return err
				}

				if _, err := io.Copy(out, tr); err != nil {
					return err
				}

				out.Close()
			}
		}

	}

}
