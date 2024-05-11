package geoip

import (
	"TOomaAh/emby_exporter_go/pkg/logger"
	"log"
	"os"
)

var (
	geoIP GeoIP
)

type GeoIP interface {
	GetCountryCode(ip string) string
	GetCountryName(ip string) string
	GetCity(ip string) string
	GetContinent(ip string) string
	GetLocation(ip string) (float64, float64)
	GetPostalCode(ip string) string
	GetRegion(ip string) string
	SetLogger(logger logger.Interface)
	Close()
}

func InitGeoIPDatabase(enable bool) (GeoIP, error) {

	switch enable {
	case true:
		file := os.Getenv("GEOIP_DB")

		if file == "" {
			file = "geoip.mmdb"
		}

		geoIP = newGeoIP(file)
		return geoIP, nil
	default:
		geoIP = newNoGeoIP()
		return geoIP, nil
	}
}

func GetGeoIPDatabase() GeoIP {
	if geoIP == nil {
		log.Fatal("GeoIP database is not initialized")
	}
	return geoIP
}
