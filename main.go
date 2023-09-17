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
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var Options struct {
	ConfFile string `short:"c" long:"config" description:"Path of your configuration file" required:"false"`
}

func setTimeZone() {
	if tz := os.Getenv("TZ"); tz != "" {
		loc, err := time.LoadLocation(tz)
		if err != nil {
			log.Printf("Timezone %s is not valid, using utc as default", tz)
			time.Local = time.UTC
			return
		}
		log.Printf("Using timezone %s", tz)
		time.Local = loc
	} else {
		time.Local = time.UTC
	}
}

func init() {
	setTimeZone()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	db := geoip.GetGeoIPDatabase()
	go func() {
		<-c
		defer db.Reader.Close()
		log.Println("Stopping server...")
		os.Exit(0)
	}()
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
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
	http.ListenAndServe(fmt.Sprintf(":%d", config.Exporter.Port|9210), logRequest(http.DefaultServeMux))

}
