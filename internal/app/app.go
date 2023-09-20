package app

import (
	"TOomaAh/emby_exporter_go/conf"
	"TOomaAh/emby_exporter_go/internal/metrics"
	"TOomaAh/emby_exporter_go/pkg/emby"
	"TOomaAh/emby_exporter_go/pkg/geoip"
	"TOomaAh/emby_exporter_go/pkg/logger"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func logRequest(handler http.Handler) http.Handler {
	logger := logger.New("info")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func Run(config *conf.Config, logger logger.Interface) {

	_ = geoip.GetGeoIPDatabase()

	// Waiting signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	db := geoip.GetGeoIPDatabase()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		<-interrupt
		defer db.Reader.Close()
		logger.Info("Stopping server...")
		os.Exit(0)
	}()

	embyServer := emby.NewServer(config.Server.Hostname, config.Server.Token, config.Server.UserID, config.Server.Port, config.Options.GeoIP, logger)
	errorPing := embyServer.Ping()
	if errorPing != nil {
		logger.Error("Server is not reachable")
	}

	client := emby.NewEmbyClient(embyServer, logger)
	embyCollector := metrics.NewEmbyCollector(client)
	newRegistry := prometheus.NewRegistry()

	newRegistry.MustRegister(embyCollector)
	handler := promhttp.HandlerFor(newRegistry, promhttp.HandlerOpts{})
	http.Handle("/metrics", handler)
	logger.Info("Beginning to serve on port %d", config.Exporter.Port|9210)
	logger.Info("You can see the metrics on http://localhost:%d/metrics", config.Exporter.Port|9210)
	http.ListenAndServe(fmt.Sprintf(":%d", config.Exporter.Port|9210), logRequest(http.DefaultServeMux))
}
