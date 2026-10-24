package geoip

// NoGeoIP implements the GeoIP interface with no-op functions.
// It's used when GeoIP functionality is disabled or unavailable.
type NoGeoIP struct{}

// NewNoGeoIP creates a new NoGeoIP instance.
func NewNoGeoIP() *NoGeoIP {
	return &NoGeoIP{}
}

func (n *NoGeoIP) Lookup(ip string) (*Record, error) {
	return nil, ErrGeoIPDisabled
}

// Close is a no-op for NoGeoIP.
func (n *NoGeoIP) Close() error {
	return nil
}
