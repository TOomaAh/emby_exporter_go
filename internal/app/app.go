package app

import (
	"TOomaAh/emby_exporter_go/conf"
	"TOomaAh/emby_exporter_go/internal/metrics"
	"TOomaAh/emby_exporter_go/pkg/emby"
	"TOomaAh/emby_exporter_go/pkg/geoip"
	"TOomaAh/emby_exporter_go/pkg/logger"
	"TOomaAh/emby_exporter_go/series"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func logRequest(handler http.Handler, logger *logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func Run(config *conf.Config) {
	l := logger.New("info")

	_ = geoip.GetGeoIPDatabase()

	// Waiting signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	db := geoip.GetGeoIPDatabase()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	select {
	case <-interrupt:
		defer db.Reader.Close()
		l.Info("Stopping server...")
		os.Exit(0)
	}

	embyServer := emby.NewServer(config.Server.Hostname, config.Server.Token, config.Server.UserID, config.Server.Port, config.Options.GeoIP)

	errorPing := embyServer.Ping()
	if errorPing != nil {
		l.Error("Server is not reachable")
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
	l.Info("Beginning to serve on port %d", config.Exporter.Port|9210)
	l.Info("You can see the metrics on http://localhost:%d/metrics", config.Exporter.Port|9210)
	http.ListenAndServe(fmt.Sprintf(":%d", config.Exporter.Port|9210), logRequest(http.DefaultServeMux, l))
}
