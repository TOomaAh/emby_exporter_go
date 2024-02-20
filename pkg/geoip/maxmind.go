package geoip

import (
	"TOomaAh/emby_exporter_go/pkg/logger"
	"log"
	"net"
	"os"

	"github.com/oschwald/geoip2-golang"
)

type GeoIPDatabase struct {
	reader *geoip2.Reader
	logger logger.Interface
}

func newGeoIP(file string) *GeoIPDatabase {
	var db *geoip2.Reader
	var err error

	db, err = geoip2.Open(file)

	if err != nil {
		log.Fatalf("Error while opening GeoIP database: %s", err)
		os.Exit(-1)
	}

	return &GeoIPDatabase{
		reader: db,
	}
}

func (g *GeoIPDatabase) SetLogger(logger logger.Interface) {
	g.logger = logger
}

func (g *GeoIPDatabase) GetCountryCode(ip string) string {
	record, err := g.reader.Country(net.ParseIP(ip))

	if err != nil {
		g.logger.Error("Error while getting country code: %s", err)
		return ""
	}

	return record.Country.IsoCode
}

func (g *GeoIPDatabase) GetCountryName(ip string) string {
	record, err := g.reader.Country(net.ParseIP(ip))

	if err != nil {
		g.logger.Error("Error while getting country name: %s", err)
		return ""
	}

	return record.Country.Names["en"]
}

func (g *GeoIPDatabase) GetCity(ip string) string {
	record, err := g.reader.City(net.ParseIP(ip))

	if err != nil {
		g.logger.Error("Error while getting city: %s", err)
		return ""
	}

	return record.City.Names["en"]
}

func (g *GeoIPDatabase) GetContinent(ip string) string {
	record, err := g.reader.City(net.ParseIP(ip))

	if err != nil {
		g.logger.Error("Error while getting continent: %s", err)
		return ""
	}

	return record.Continent.Names["en"]
}

func (g *GeoIPDatabase) GetLocation(ip string) (float64, float64) {
	record, err := g.reader.City(net.ParseIP(ip))

	if err != nil {
		g.logger.Error("Error while getting location: %s", err)
		return 0, 0
	}

	return record.Location.Latitude, record.Location.Longitude
}

func (g *GeoIPDatabase) GetPostalCode(ip string) string {
	record, err := g.reader.City(net.ParseIP(ip))

	if err != nil {
		g.logger.Error("Error while getting postal code: %s", err)
		return ""
	}

	return record.Postal.Code
}

func (g *GeoIPDatabase) GetRegion(ip string) string {
	record, err := g.reader.City(net.ParseIP(ip))

	if err != nil {
		g.logger.Error("Error while getting region: %s", err)
		return ""
	}

	if len(record.Subdivisions) == 0 {
		return ""
	}
	return record.Subdivisions[0].Names["en"]
}

func (g *GeoIPDatabase) Close() {
	g.reader.Close()
}
