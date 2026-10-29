package app

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"TOomaAh/emby_exporter_go/internal/metrics"
	"TOomaAh/emby_exporter_go/pkg/emby"
	"TOomaAh/emby_exporter_go/pkg/logger"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func newEmbyServer(t *testing.T, url string) *emby.Server {
	t.Helper()
	host, port, err := net.SplitHostPort(strings.TrimPrefix(url, "http://"))
	if err != nil {
		t.Fatal(err)
	}
	return emby.NewServer(&emby.ServerInfo{Hostname: "http://" + host, Port: port, Token: "token"}, logger.New("error"))
}

// newMetricsHandler wires the middleware like Run does, with a fake full handler.
func newMetricsHandler(server *emby.Server) http.Handler {
	upRegistry := prometheus.NewRegistry()
	upRegistry.MustRegister(metrics.NewUpCollector(server))
	upHandler := promhttp.HandlerFor(upRegistry, promhttp.HandlerOpts{})

	full := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all metrics"))
	})

	// A long interval: the ping loop started by the middleware must not spin during the test.
	return metricHandlerMiddleware(full, upHandler, server, time.Hour)
}

func scrape(t *testing.T, h http.Handler) (int, string) {
	t.Helper()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	body, _ := io.ReadAll(rec.Result().Body)
	return rec.Code, string(body)
}

func TestMetricsWhenEmbyIsUnreachable(t *testing.T) {
	fake := httptest.NewServer(http.NotFoundHandler())
	server := newEmbyServer(t, fake.URL)
	fake.Close()

	code, body := scrape(t, newMetricsHandler(server))

	if code != http.StatusOK {
		t.Errorf("status = %d, want 200", code)
	}
	if !strings.Contains(body, "emby_up 0") {
		t.Errorf("body does not contain emby_up 0:\n%s", body)
	}
	if strings.Contains(body, "all metrics") {
		t.Error("the full metrics must not be collected while Emby is unreachable")
	}
}

func TestMetricsWhenEmbyIsReachable(t *testing.T) {
	fake := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	t.Cleanup(fake.Close)
	server := newEmbyServer(t, fake.URL)
	if err := server.Ping(); err != nil {
		t.Fatal(err)
	}

	_, body := scrape(t, newMetricsHandler(server))

	if body != "all metrics" {
		t.Errorf("body = %q, want the full metrics", body)
	}
}

func TestHealthCheck(t *testing.T) {
	rec := httptest.NewRecorder()
	HealthCheck(rec, httptest.NewRequest(http.MethodGet, "/health", nil))
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
}
