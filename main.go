package main

import (
	"TOomaAh/emby_exporter_go/conf"
	"TOomaAh/emby_exporter_go/emby"
	"TOomaAh/emby_exporter_go/geoip"
	"TOomaAh/emby_exporter_go/metrics"
	"TOomaAh/emby_exporter_go/series"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/jessevdk/go-flags"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var Options struct {
	ConfFile string `short:"c" long:"config" description:"Path of your configuration file" required:"false"`
}

func init() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	db := geoip.GetGeoIPDatabase()
	go func() {
		for _ = range c {
			defer db.Reader.Close()
			log.Println("Stopping server...")
			os.Exit(0)
		}
	}()
}

func main() {
	_, err := flags.ParseArgs(&Options, os.Args)

	if err != nil {
		log.Fatalln(err)
	}

	_ = geoip.GetGeoIPDatabase()

	config, err := conf.NewConfig(Options.ConfFile)

	if err != nil {
		log.Fatalln(err)
	}

	logger := log.Default()

	embyServer := emby.NewServer(config.Server.Hostname, config.Server.Token, config.Server.UserID, config.Server.Port, config.Options.GeoIP)

	errorPing := embyServer.Ping()
	if errorPing != nil {
		logger.Fatalln("Server is not reachable")
	}

	client := emby.NewEmbyClient(embyServer)

	seriesInt := series.NewSeriesFromConf(config)

	embyCollector := metrics.NewEmbyCollector(client)

	newRegistry := prometheus.NewRegistry()

	if seriesInt != nil {
		serieCollector := series.NewSeriesCollector(&seriesInt)
		newRegistry.MustRegister(serieCollector)
	}

	newRegistry.MustRegister(embyCollector)
	handler := promhttp.HandlerFor(newRegistry, promhttp.HandlerOpts{})
	http.Handle("/metrics", handler)
	logger.Printf("Beginning to serve on port %d", config.Exporter.Port|9210)
	logger.Printf("You can see the metrics on http://localhost:%d/metrics", config.Exporter.Port|9210)
	http.ListenAndServe(fmt.Sprintf(":%d", config.Exporter.Port|9210), nil)

}
