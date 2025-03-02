package geoip

import (
	"TOomaAh/emby_exporter_go/internal/conf"
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
	Close()
}

func InitGeoIPDatabase(cfg *conf.Config, logger logger.Interface) (GeoIP, error) {
	switch cfg.Options.GeoIP {
	case true:
		file := os.Getenv("GEOIP_DB")
		if file == "" {
			file = "geoip.mmdb"
		}

		if _, err := os.Stat(file); os.IsNotExist(err) {
			if cfg.Options.GeoIPOptions.AccountId == "" || cfg.Options.GeoIPOptions.LicenceKey == "" {
				logger.Error("GeoIP database is not found and no account ID or licence key is provided")
				return newNoGeoIP(), err
			}
			geoIP, err := newGeoIP(file, cfg.Options.GeoIPOptions.AccountId, cfg.Options.GeoIPOptions.LicenceKey)

			if err != nil {
				return newNoGeoIP(), nil
			}

			return geoIP, nil
		}

		geoIP, err := newGeoIP(file, cfg.Options.GeoIPOptions.AccountId, cfg.Options.GeoIPOptions.LicenceKey)

		if err != nil {
			logger.Error("Error while opening GeoIP database: %s", err)
			return nil, err
		}

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
