package main

import (
	"TOomaAh/emby_exporter_go/conf"
	"TOomaAh/emby_exporter_go/emby"
	"TOomaAh/emby_exporter_go/metrics"
	"TOomaAh/emby_exporter_go/series"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var Options struct {
	ConfFile string `short:"c" long:"config" description:"Path of your configuration file" required:"false"`
}

func main() {
	_, err := flags.ParseArgs(&Options, os.Args)

	if err != nil {
		log.Fatalln(err)
	}
	var configFile string
	if Options.ConfFile == "" {
		configFile = "./config.yml"
	} else {
		configFile = Options.ConfFile
	}
	config, err := conf.NewConfig(configFile)

	if err != nil {
		log.Fatalln(err)
	}

	logger := log.Default()

	embyServer := emby.NewServer(config.Server.Hostname, config.Server.Token, config.Server.UserID, config.Server.Port)

	errorPing := embyServer.Ping()
	if errorPing != nil {
		logger.Fatalln("Server is not reachable")
	}

	client := emby.NewEmbyClient(embyServer)

	seriesInt := series.NewSeriesFromConf(config)

	embyCollector := metrics.NewEmbyCollector(client)

	newRegistry := prometheus.NewRegistry()

	if seriesInt != nil {
		serieCollector := series.NewSeriesCollector(seriesInt)
		newRegistry.MustRegister(serieCollector)
	}

	newRegistry.MustRegister(embyCollector)
	handler := promhttp.HandlerFor(newRegistry, promhttp.HandlerOpts{})
	http.Handle("/metrics", handler)
	logger.Printf("Beginning to serve on port %d", config.Exporter.Port|9210)
	logger.Printf("You can see the metrics on http://localhost:%d/metrics", config.Exporter.Port|9210)
	http.ListenAndServe(fmt.Sprintf(":%d", config.Exporter.Port|9210), nil)

}
