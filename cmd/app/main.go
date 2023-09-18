package main

import (
	"TOomaAh/emby_exporter_go/conf"
	"TOomaAh/emby_exporter_go/internal/app"
	"TOomaAh/emby_exporter_go/pkg/logger"
	"log"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
)

var Options struct {
	ConfFile string `short:"c" long:"config" description:"Path of your configuration file" required:"false"`
}

func setTimeZone() {
	l := logger.Get()
	if tz := os.Getenv("TZ"); tz != "" {
		loc, err := time.LoadLocation(tz)
		if err != nil {
			l.Info("Timezone %s is not valid, using utc as default", tz)
			time.Local = time.UTC
			return
		}
		l.Info("Using timezone %s", tz)
		time.Local = loc
	} else {
		time.Local = time.UTC
	}
}

func init() {
	setTimeZone()
}

func main() {
	_, err := flags.ParseArgs(&Options, os.Args)

	if err != nil {
		log.Fatalln(err)
	}

	config, err := conf.NewConfig(Options.ConfFile)

	if err != nil {
		log.Fatalln(err)
	}

	app.Run(config)

}
