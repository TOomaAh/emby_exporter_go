package geoip

import "TOomaAh/emby_exporter_go/pkg/logger"

// NoGeoIP implements the GeoIP interface with no-op functions.
// It's used when GeoIP functionality is disabled or unavailable.
type NoGeoIP struct {
	logger logger.Interface
}

// NewNoGeoIP creates a new NoGeoIP instance.
func NewNoGeoIP() *NoGeoIP {
	return &NoGeoIP{
		logger: logger.New("info"),
	}
}

// GetCountryCode returns an empty string and nil error for NoGeoIP.
func (n *NoGeoIP) GetCountryCode(ip string) (string, error) {
	return "", nil
}

// GetCountryName returns an empty string and nil error for NoGeoIP.
func (n *NoGeoIP) GetCountryName(ip string) (string, error) {
	return "", nil
}

// GetCity returns an empty string and nil error for NoGeoIP.
func (n *NoGeoIP) GetCity(ip string) (string, error) {
	return "", nil
}

// GetContinent returns an empty string and nil error for NoGeoIP.
func (n *NoGeoIP) GetContinent(ip string) (string, error) {
	return "", nil
}

// GetLocation returns zeros and nil error for NoGeoIP.
func (n *NoGeoIP) GetLocation(ip string) (float64, float64, error) {
	return 0, 0, nil
}

// GetPostalCode returns an empty string and nil error for NoGeoIP.
func (n *NoGeoIP) GetPostalCode(ip string) (string, error) {
	return "", nil
}

// GetRegion returns an empty string and nil error for NoGeoIP.
func (n *NoGeoIP) GetRegion(ip string) (string, error) {
	return "", nil
}

// Close is a no-op for NoGeoIP.
func (n *NoGeoIP) Close() error {
	return nil
}
