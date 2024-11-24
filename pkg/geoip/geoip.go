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

func InitGeoIPDatabase(enable bool, logger logger.Interface) (GeoIP, error) {
	switch enable {
	case true:
		file := os.Getenv("GEOIP_DB")
		if file == "" {
			file = "geoip.mmdb"
		}

		if _, err := os.Stat(file); os.IsNotExist(err) {
			logger.Error("GeoIP database file does not exist: %s, disable geoip", file)
			return newNoGeoIP(), nil
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
