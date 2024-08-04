package app

import (
	"TOomaAh/emby_exporter_go/internal/conf"
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

	db, err := geoip.InitGeoIPDatabase(config.Options.GeoIP)
	if err != nil {
		logger.Error("GeoIP database is not initialized")
		os.Exit(-1)
	}
	db.SetLogger(logger)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		<-interrupt
		defer db.Close()
		logger.Info("Stopping server...")
		os.Exit(0)
	}()

	embyServer := emby.NewServer(&emby.ServerInfo{
		Hostname: config.Server.Hostname,
		Port:     config.Server.Port,
		UserID:   config.Server.UserID,
		Token:    config.Server.Token,
	}, logger)
	errorPing := embyServer.Ping()
	if errorPing != nil {
		logger.Error("Server is not reachable")
	}

	activityCollector := metrics.NewActivityCollector(embyServer)
	alertCollector := metrics.NewAlertCollector(embyServer)
	libraryCollector := metrics.NewLibraryCollector(embyServer)
	sessionCollector := metrics.NewSessionCollector(embyServer)
	systemInfoCollector := metrics.NewSystemInfoCollector(embyServer)

	newRegistry := prometheus.NewRegistry()

	newRegistry.MustRegister(alertCollector, libraryCollector, sessionCollector, systemInfoCollector, activityCollector)
	handler := promhttp.HandlerFor(newRegistry, promhttp.HandlerOpts{})
	http.Handle("/metrics", handler)
	logger.Info("Beginning to serve on port %d", config.Exporter.Port|9210)
	logger.Info("You can see the metrics on http://localhost:%d/metrics", config.Exporter.Port|9210)
	http.ListenAndServe(fmt.Sprintf(":%d", config.Exporter.Port|9210), logRequest(http.DefaultServeMux))
}
