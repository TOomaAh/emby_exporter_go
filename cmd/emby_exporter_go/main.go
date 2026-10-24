package main

import (
	"TOomaAh/emby_exporter_go/internal/app"
	"TOomaAh/emby_exporter_go/internal/conf"
	"TOomaAh/emby_exporter_go/pkg/geoip"
	"TOomaAh/emby_exporter_go/pkg/logger"
	"fmt"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
)

// version is set at build time with -ldflags "-X main.version=...".
var version = "dev"

var Options struct {
	Version           bool   `short:"v" long:"version" description:"Print the version and exit"`
	ConfFile          string `short:"c" long:"config" env:"CONFIG_FILE" description:"Path of your configuration file" required:"false"`
	GeoIPDatabaseFile string `short:"g" long:"geoip" description:"Path of your GeoIP database file" required:"false"`
	LogLevel          string `short:"l" long:"log-level" env:"LOG_LEVEL" description:"Log level (debug, info, warn, error)" required:"false" default:"info"`
}

func init() {
	tz := os.Getenv("TZ")
	if tz == "" {
		time.Local = time.UTC
		return
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		time.Local = time.UTC
		return
	}
	time.Local = loc
}

func main() {
	if _, err := flags.ParseArgs(&Options, os.Args); err != nil {
		if fe, ok := err.(*flags.Error); ok && fe.Type == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	if Options.Version {
		fmt.Println(version)
		os.Exit(0)
	}

	l := logger.New(Options.LogLevel)

	l.Info("emby_exporter %s", version)
	l.Info("Using %s", time.Local.String())

	config, err := conf.NewConfig(Options.ConfFile)

	if err != nil {
		l.Fatal(err)
	}

	geoipDatabase := Options.GeoIPDatabaseFile

	if geoipDatabase != "" {
		os.Setenv("GEOIP_DB", geoipDatabase)
	}

	geoIp, err := geoip.InitGeoIPDatabase(config, l)

	if err != nil {
		l.Error(err)
	}

	app.Run(config, geoIp, l)

}
