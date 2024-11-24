package conf

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Server struct {
	Hostname string `yaml:"url"`
	Port     string `yaml:"port"`
	Token    string `yaml:"token"`
	UserID   string `yaml:"userID"`
}

type Config struct {
	Exporter struct {
		Port int `yaml:"port"`
	} `yaml:"exporter,omitempty"`
	Server  Server `yaml:"server,omitempty"`
	Options struct {
		GeoIP bool `yaml:"geoip" default:"false"`
	} `yaml:"options,omitempty"`
}

func NewConfig(path string) (*Config, error) {
	if path == "" {
		path = "./config.yml"
	}
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("cannot open configuration file %s\n", path)
		os.Exit(-1)
	}

	defer file.Close()

	var config Config

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)

	if config.Server.Hostname == "" {
		config.Server.Hostname = "http://localhost"
	} else {
		// Add http:// if not present
		if len(config.Server.Hostname) < 7 || config.Server.Hostname[:7] != "http://" {
			config.Server.Hostname = "http://" + config.Server.Hostname
		}
	}

	if config.Server.Port == "" {
		config.Server.Port = "8096"
	}

	if err != nil {
		fmt.Printf("Cannot decode config file: %s", err)
		os.Exit(-1)
	}

	return &config, nil

}
