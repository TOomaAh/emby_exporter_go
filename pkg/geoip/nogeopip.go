package geoip

import "TOomaAh/emby_exporter_go/pkg/logger"

type NoGeoIpDatabase struct{}

func newNoGeoIP() *NoGeoIpDatabase {
	return &NoGeoIpDatabase{}
}

func (g *NoGeoIpDatabase) SetLogger(logger logger.Interface) {
}

func (g *NoGeoIpDatabase) GetCountryCode(ip string) string {
	return ""
}

func (g *NoGeoIpDatabase) GetCountryName(ip string) string {
	return ""
}

func (g *NoGeoIpDatabase) GetCity(ip string) string {
	return ""
}

func (g *NoGeoIpDatabase) GetContinent(ip string) string {
	return ""
}

func (g *NoGeoIpDatabase) GetLocation(ip string) (float64, float64) {
	return 0, 0
}

func (g *NoGeoIpDatabase) GetPostalCode(ip string) string {
	return ""
}

func (g *NoGeoIpDatabase) GetRegion(ip string) string {
	return ""
}

func (g *NoGeoIpDatabase) Close() {
}
