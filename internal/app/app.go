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

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func Run(config *conf.Config, geoIp geoip.GeoIP, logger logger.Interface) {

	embyServer := emby.NewServer(&emby.ServerInfo{
		Hostname: config.Server.Hostname,
		Port:     config.Server.Port,
		UserID:   config.Server.UserID,
		Token:    config.Server.Token,
	}, logger)

	errorPing := embyServer.Ping()

	if errorPing != nil {
		logger.Error("Server is not reachable")
		os.Exit(-1)
	}

	activityCollector := metrics.NewActivityCollector(embyServer, logger)
	alertCollector := metrics.NewAlertCollector(embyServer, logger)
	libraryCollector := metrics.NewLibraryCollector(embyServer, logger)
	sessionCollector := metrics.NewSessionCollector(embyServer, geoIp, logger)
	systemInfoCollector := metrics.NewSystemInfoCollector(embyServer, logger)

	newRegistry := prometheus.NewRegistry()

	newRegistry.MustRegister(alertCollector, libraryCollector, sessionCollector, systemInfoCollector, activityCollector)
	handler := promhttp.HandlerFor(newRegistry, promhttp.HandlerOpts{})
	http.Handle("/metrics", handler)
	logger.Info("Beginning to serve on port %d", config.Exporter.Port|9210)
	logger.Info("You can see the metrics on http://localhost:%d/metrics", config.Exporter.Port|9210)
	http.ListenAndServe(fmt.Sprintf(":%d", config.Exporter.Port|9210), logRequest(http.DefaultServeMux))
}
