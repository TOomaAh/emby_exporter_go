package geoip

import (
	"net"

	"github.com/oschwald/geoip2-golang"
)

type GeoIPDatabase struct {
	Reader *geoip2.Reader
}

var (
	geoIP *GeoIPDatabase
)

func newGeoIP(file string) *GeoIPDatabase {
	db, err := geoip2.Open(file)

	if err != nil {
		panic(err)
	}

	return &GeoIPDatabase{
		Reader: db,
	}
}

func GetGeoIPDatabase() *GeoIPDatabase {
	if geoIP == nil {
		geoIP = newGeoIP("./geoip.mmdb")
	}

	return geoIP
}

func (g *GeoIPDatabase) GetCountryCode(ip string) string {
	record, err := g.Reader.Country(net.ParseIP(ip))

	if err != nil {
		panic(err)
	}

	return record.Country.IsoCode
}

func (g *GeoIPDatabase) GetCountryName(ip string) string {
	record, err := g.Reader.Country(net.ParseIP(ip))

	if err != nil {
		panic(err)
	}

	return record.Country.Names["en"]
}

func (g *GeoIPDatabase) GetCity(ip string) string {
	record, err := g.Reader.City(net.ParseIP(ip))

	if err != nil {
		panic(err)
	}

	return record.City.Names["en"]
}

func (g *GeoIPDatabase) GetContinent(ip string) string {
	record, err := g.Reader.City(net.ParseIP(ip))

	if err != nil {
		panic(err)
	}

	return record.Continent.Names["en"]
}

func (g *GeoIPDatabase) GetLocation(ip string) (float64, float64) {
	record, err := g.Reader.City(net.ParseIP(ip))

	if err != nil {
		panic(err)
	}

	return record.Location.Latitude, record.Location.Longitude
}

func (g *GeoIPDatabase) GetPostalCode(ip string) string {
	record, err := g.Reader.City(net.ParseIP(ip))

	if err != nil {
		panic(err)
	}

	return record.Postal.Code
}

func (g *GeoIPDatabase) GetRegion(ip string) string {
	record, err := g.Reader.City(net.ParseIP(ip))

	if err != nil {
		panic(err)
	}

	if len(record.Subdivisions) == 0 {
		return ""
	}
	return record.Subdivisions[0].Names["en"]
}
