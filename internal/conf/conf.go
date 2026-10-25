package conf

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// Server holds the Emby server configuration.
type Server struct {
	Hostname string `yaml:"url"`
	Port     string `yaml:"port"`
	Token    string `yaml:"token"`
	UserID   string `yaml:"userID"`
}

type GeoIPOptions struct {
	AccountId  string `yaml:"accountId"`
	LicenseKey string `yaml:"licenseKey"`
}

type Options struct {
	RetryInterval int          `yaml:"retryInterval"`
	GeoIP         bool         `yaml:"geoip"`
	HealthCheck   bool         `yaml:"healthcheck"`
	GeoIPOptions  GeoIPOptions `yaml:"geoipOptions,omitempty"`
}

type Exporter struct {
	Port int `yaml:"port" default:"9210"`
}

// Config holds the overall configuration.
type Config struct {
	Exporter Exporter `yaml:"exporter,omitempty"`
	Server   Server   `yaml:"server,omitempty"`
	Options  Options  `yaml:"options,omitempty"`
}

var defaultConfig = Config{
	Exporter: Exporter{
		Port: 9210,
	},
	Server: Server{
		Hostname: "localhost",
		Port:     "8096",
		Token:    "",
		UserID:   "",
	},
	Options: Options{
		RetryInterval: 10,
		GeoIP:         false,
		HealthCheck:   true,
		GeoIPOptions: GeoIPOptions{
			AccountId:  "",
			LicenseKey: "",
		},
	},
}

// NewConfig reads the YAML configuration from the specified path and returns a Config instance.
// If no path is provided, it defaults to "./config.yml". It also applies default values for missing fields.
func NewConfig(path string) (*Config, error) {
	if path == "" {
		path = "./config.yml"
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open configuration file %s: %w", path, err)
	}
	defer file.Close()

	cfg := defaultConfig

	if err := yaml.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("cannot decode configuration file %s: %w", path, err)
	}

	// Ensure the server hostname starts with "http://" or "https://".
	if cfg.Server.Hostname != "" && !hasHTTPPrefix(cfg.Server.Hostname) {
		cfg.Server.Hostname = "http://" + cfg.Server.Hostname
	}

	// Set default RetryInterval if not provided.
	if cfg.Options.RetryInterval <= 0 {
		cfg.Options.RetryInterval = defaultConfig.Options.RetryInterval
	}

	if cfg.Exporter.Port <= 0 {
		cfg.Exporter.Port = defaultConfig.Exporter.Port
	}

	return &cfg, nil
}

// hasHTTPPrefix checks if the provided hostname starts with "http://" or "https://".
func hasHTTPPrefix(hostname string) bool {
	return strings.HasPrefix(hostname, "http://") || strings.HasPrefix(hostname, "https://")
}
