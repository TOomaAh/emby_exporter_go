package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"TOomaAh/emby_exporter_go/internal/conf"
	"TOomaAh/emby_exporter_go/internal/metrics"
	"TOomaAh/emby_exporter_go/pkg/emby"
	"TOomaAh/emby_exporter_go/pkg/geoip"
	"TOomaAh/emby_exporter_go/pkg/logger"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// logRequest is a middleware that logs each HTTP request.
func logRequest(handler http.Handler, log logger.Interface) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

// metricHandlerMiddleware serves only emby_up while the server is unreachable.
// If the server is unreachable, it starts a ping loop (if not already running).
func metricHandlerMiddleware(next http.Handler, unreachable http.Handler, server *emby.Server, interval time.Duration) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !server.IsReachable() {
			go server.StartPingLoop(interval)
			unreachable.ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r)

		if !server.IsReachable() {
			go server.StartPingLoop(interval)
		}
	})
}

// HealthCheck simply returns HTTP 200, which can be used to verify the container's health.
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// Run initializes the Prometheus collectors and starts the HTTP server to export metrics.
func Run(cfg *conf.Config, geoIp geoip.GeoIP, log logger.Interface) {
	// Create the Emby server instance.
	embyServer := emby.NewServer(&emby.ServerInfo{
		Hostname: cfg.Server.Hostname,
		Port:     cfg.Server.Port,
		Token:    cfg.Server.Token,
	}, log)

	retryInterval := time.Duration(cfg.Options.RetryInterval) * time.Second

	if err := embyServer.Ping(); err != nil {
		log.Error("Initial ping to Emby server failed: %s", err)
		go embyServer.StartPingLoop(retryInterval)
	} else {
		log.Info("Successfully connected to Emby server")
	}

	// Initialize collectors.
	activityCollector := metrics.NewActivityCollector(embyServer, log)
	alertCollector := metrics.NewAlertCollector(embyServer, log)
	libraryCollector := metrics.NewLibraryCollector(embyServer, log)
	sessionCollector := metrics.NewSessionCollector(embyServer, geoIp, log)
	systemInfoCollector := metrics.NewSystemInfoCollector(embyServer, log)
	upCollector := metrics.NewUpCollector(embyServer)

	// Create a Prometheus registry and register the collectors.
	registry := prometheus.NewRegistry()
	registry.MustRegister(alertCollector, libraryCollector, sessionCollector, systemInfoCollector, activityCollector, upCollector)
	registry.MustRegister(collectors.NewGoCollector(), collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	// Registry served while Emby is unreachable: only emby_up.
	upRegistry := prometheus.NewRegistry()
	upRegistry.MustRegister(upCollector)

	serverMux := http.NewServeMux()

	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	upHandler := promhttp.HandlerFor(upRegistry, promhttp.HandlerOpts{})

	// Wrap the /metrics endpoint with middleware that controls access based on the server status.
	serverMux.Handle("/metrics", metricHandlerMiddleware(handler, upHandler, embyServer, retryInterval))
	if cfg.Options.HealthCheck {
		serverMux.HandleFunc("/health", HealthCheck)
	}

	// Start the HTTP server.
	port := cfg.Exporter.Port

	if port == 0 {
		port = 9210
	}

	log.Info("Beginning to serve on port %d", port)
	log.Info("You can see the metrics on http://localhost:%d/metrics", port)

	// listen for signals
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 30 * time.Second,
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      logRequest(serverMux, log),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("HTTP server error: %v", err)
		}
	}()

	<-ch

	log.Info("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Error shutting down server: %v", err)
	} else {
		log.Info("Server shut down gracefully")
	}

	if err := geoIp.Close(); err != nil {
		log.Error("Error closing GeoIP database: %v", err)
	}
}
