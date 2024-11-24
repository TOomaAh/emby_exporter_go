package main

import (
	"TOomaAh/emby_exporter_go/internal/app"
	"TOomaAh/emby_exporter_go/internal/conf"
	"TOomaAh/emby_exporter_go/pkg/geoip"
	"TOomaAh/emby_exporter_go/pkg/logger"
	"os"
	"os/signal"
	"time"

	"github.com/jessevdk/go-flags"
)

var Options struct {
	ConfFile          string `short:"c" long:"config" description:"Path of your configuration file" required:"false"`
	GeoIPDatabaseFile string `short:"g" long:"geoip" description:"Path of your GeoIP database file" required:"false"`
}

func setTimeZone() {

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

func init() {
	setTimeZone()
}

func main() {

	l := logger.New("info")

	l.Info("Using %s", time.Local.String())
	_, err := flags.ParseArgs(&Options, os.Args)

	if err != nil {
		l.Fatal(err)
	}

	config, err := conf.NewConfig(Options.ConfFile)

	if err != nil {
		l.Fatal(err)
	}

	geoipDatabase := Options.GeoIPDatabaseFile

	if geoipDatabase != "" {
		os.Setenv("GEOIP_DB", geoipDatabase)
	}

	db, err := geoip.InitGeoIPDatabase(config.Options.GeoIP)
	if err != nil {
		l.Error("GeoIP database is not initialized")
		os.Exit(-1)
	}
	db.SetLogger(l)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		<-interrupt
		defer db.Close()
		l.Info("Stopping server...")
		os.Exit(0)
	}()

	app.Run(config, l)

}
