package conf

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// Server holds the Emby server configuration.
type Server struct {
	Hostname string `yaml:"url" default:"localhost"`
	Port     string `yaml:"port" default:"8096"`
	Token    string `yaml:"token"`
	UserID   string `yaml:"userID"`
}

// Config holds the overall configuration.
type Config struct {
	Exporter struct {
		Port int `yaml:"port" default:"9210"`
	} `yaml:"exporter,omitempty"`
	Server  Server `yaml:"server,omitempty"`
	Options struct {
		RetryInterval int  `yaml:"retryInterval" default:"10"`
		GeoIP         bool `yaml:"geoip" default:"false"`
		HealthCheck   bool `yaml:"healthcheck" default:"false"`
	} `yaml:"options,omitempty"`
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

	var cfg Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("cannot decode config file: %w", err)
	}

	// Ensure the server hostname starts with "http://" or "https://".
	if cfg.Server.Hostname != "" && !hasHTTPPrefix(cfg.Server.Hostname) {
		cfg.Server.Hostname = "http://" + cfg.Server.Hostname
	}

	// Set default RetryInterval if not provided.
	if cfg.Options.RetryInterval == 0 {
		cfg.Options.RetryInterval = 10
	}

	return &cfg, nil
}

// hasHTTPPrefix checks if the provided hostname starts with "http://" or "https://".
func hasHTTPPrefix(hostname string) bool {
	return strings.HasPrefix(hostname, "http://") || strings.HasPrefix(hostname, "https://")
}
