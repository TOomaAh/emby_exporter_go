package app

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"TOomaAh/emby_exporter_go/internal/conf"
	"TOomaAh/emby_exporter_go/internal/metrics"
	"TOomaAh/emby_exporter_go/pkg/emby"
	"TOomaAh/emby_exporter_go/pkg/geoip"
	"TOomaAh/emby_exporter_go/pkg/logger"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// PingManager encapsulates the ping state and a mutex to ensure only one ping loop is running.
type PingManager struct {
	mu        sync.Mutex
	isPinging int32
}

// NewPingManager creates a new instance of PingManager.
func NewPingManager() *PingManager {
	return &PingManager{}
}

// IsPinging returns true if a ping loop is already running.
func (pm *PingManager) IsPinging() bool {
	return atomic.LoadInt32(&pm.isPinging) == 1
}

// SetPinging sets the ping loop state (true = running, false = not running).
func (pm *PingManager) SetPinging(active bool) {
	if active {
		atomic.StoreInt32(&pm.isPinging, 1)
	} else {
		atomic.StoreInt32(&pm.isPinging, 0)
	}
}

// logRequest is a middleware that logs each HTTP request.
func logRequest(handler http.Handler) http.Handler {
	log := logger.New("info")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

// metricHandlerMiddleware blocks access to the /metrics endpoint when the server is unreachable.
// If the server is unreachable, it starts a ping loop (if not already running) and returns a 503 error.
func metricHandlerMiddleware(next http.Handler, server *emby.Server, cfg *conf.Config, pm *PingManager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the server is unreachable at the beginning of the request,
		// start the ping loop (if not already running) and return immediately.
		if !server.IsRecheable {
			if !pm.IsPinging() {
				go startPingLoop(server, logger.New("info"), cfg, pm)
			}
			// Here, we return a 200 status.
			w.WriteHeader(http.StatusOK)
			return
		}

		// If the server is reachable, serve the metrics.
		next.ServeHTTP(w, r)

		// After serving, check again: if the server became unreachable during the request,
		// launch the ping loop (if not already running).
		if !server.IsRecheable && !pm.IsPinging() {
			go startPingLoop(server, logger.New("info"), cfg, pm)
		}
	})
}

// HealthCheck simply returns HTTP 200, which can be used to verify the container's health.
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// startPingLoop pings the server at regular intervals until it responds.
// This function blocks until the server becomes reachable and updates server.IsRecheable.
func startPingLoop(server *emby.Server, log logger.Interface, cfg *conf.Config, pm *PingManager) {
	// Ensure that no other ping loop is already running.
	if pm.IsPinging() {
		return
	}
	pm.mu.Lock()
	pm.SetPinging(true)
	defer func() {
		pm.SetPinging(false)
		pm.mu.Unlock()
	}()

	for {
		if err := server.Ping(); err != nil {
			log.Error("Server is not reachable, retrying in %d seconds...", cfg.Options.RetryInterval)
			server.IsRecheable = false
		} else {
			server.IsRecheable = true
			log.Info("Server is reachable")
			break
		}
		time.Sleep(time.Duration(cfg.Options.RetryInterval) * time.Second)
	}
}

// Run initializes the Prometheus collectors and starts the HTTP server to export metrics.
func Run(cfg *conf.Config, geoIp geoip.GeoIP, log logger.Interface) {
	// Create the Emby server instance.
	embyServer := emby.NewServer(&emby.ServerInfo{
		Hostname: cfg.Server.Hostname,
		Port:     cfg.Server.Port,
		UserID:   cfg.Server.UserID,
		Token:    cfg.Server.Token,
	}, log)

	// Create a new PingManager instance.
	pingManager := NewPingManager()

	// Perform an initial ping test at startup.
	if err := embyServer.Ping(); err != nil {
		log.Error("Initial ping failed: %v", err)
		embyServer.IsRecheable = false
		// Start the ping loop in the background.
		go startPingLoop(embyServer, log, cfg, pingManager)
	} else {
		embyServer.IsRecheable = true
		log.Info("Initial ping succeeded, server is reachable")
	}

	// Initialize collectors.
	activityCollector := metrics.NewActivityCollector(embyServer, log)
	alertCollector := metrics.NewAlertCollector(embyServer, log)
	libraryCollector := metrics.NewLibraryCollector(embyServer, log)
	sessionCollector := metrics.NewSessionCollector(embyServer, geoIp, log)
	systemInfoCollector := metrics.NewSystemInfoCollector(embyServer, log)

	// Create a Prometheus registry and register the collectors.
	registry := prometheus.NewRegistry()
	registry.MustRegister(alertCollector, libraryCollector, sessionCollector, systemInfoCollector, activityCollector)
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})

	// Wrap the /metrics endpoint with middleware that controls access based on the server status.
	http.Handle("/metrics", metricHandlerMiddleware(handler, embyServer, cfg, pingManager))

	// Start the HTTP server.
	// Note: using the bitwise OR operator (|) here seems to provide a default value.
	port := cfg.Exporter.Port | 9210
	log.Info("Beginning to serve on port %d", port)
	log.Info("You can see the metrics on http://localhost:%d/metrics", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), logRequest(http.DefaultServeMux))
	if err != nil {
		log.Error("HTTP server error: %v", err)
	}
}
