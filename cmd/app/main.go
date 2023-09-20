package main

import (
	"TOomaAh/emby_exporter_go/conf"
	"TOomaAh/emby_exporter_go/internal/app"
	"TOomaAh/emby_exporter_go/pkg/logger"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
)

var Options struct {
	ConfFile string `short:"c" long:"config" description:"Path of your configuration file" required:"false"`
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

	app.Run(config, l)

}
