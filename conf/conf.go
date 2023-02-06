package conf

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Exporter struct {
		Port int `yaml:"port"`
	} `yaml:"exporter,omitempty"`
	Server struct {
		Hostname string `yaml:"url"`
		Port     int    `yaml:"port"`
		Token    string `yaml:"token"`
		UserID   string `yaml:"userID"`
	} `yaml:"server"`
	Series struct {
		Sonarr struct {
			Url   string `yaml:"url"`
			Token string `yaml:"token"`
		} `yaml:"sonarr,omitempty"`
		Medusa struct {
			Url   string `yaml:"url"`
			Token string `yaml:"token"`
		} `yaml:"medusa,omitempty"`
	} `yaml:"series,omitempty"`
	Options struct {
		GeoIP bool `yaml:"geoip"`
	} `yaml:"options,omitempty"`
}

func NewConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Panicf("%s", err)
		os.Exit(-1)
	}

	defer file.Close()

	var config Config

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)

	if err != nil {
		log.Panicf("Cannot decode config file: %s", err)
		os.Exit(-1)
	}

	return &config, nil

}
