package metrics

import (
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"TOomaAh/emby_exporter_go/pkg/emby"
	"TOomaAh/emby_exporter_go/pkg/logger"

	"github.com/prometheus/client_golang/prometheus"
)

// newFakeEmby returns an emby.Server talking to a fake Emby served by mux.
func newFakeEmby(t *testing.T, mux *http.ServeMux) *emby.Server {
	t.Helper()
	fake := httptest.NewServer(mux)
	t.Cleanup(fake.Close)

	host, port, err := net.SplitHostPort(strings.TrimPrefix(fake.URL, "http://"))
	if err != nil {
		t.Fatal(err)
	}
	return emby.NewServer(&emby.ServerInfo{Hostname: "http://" + host, Port: port, Token: "token"}, logger.New("error"))
}

type sample struct {
	labels map[string]string
	value  float64
}

// gather collects c and returns its samples by metric name.
func gather(t *testing.T, c prometheus.Collector) map[string][]sample {
	t.Helper()
	reg := prometheus.NewRegistry()
	reg.MustRegister(c)

	families, err := reg.Gather()
	if err != nil {
		t.Fatal(err)
	}

	out := map[string][]sample{}
	for _, mf := range families {
		for _, m := range mf.GetMetric() {
			s := sample{labels: map[string]string{}, value: m.GetGauge().GetValue()}
			for _, lp := range m.GetLabel() {
				s.labels[lp.GetName()] = lp.GetValue()
			}
			out[mf.GetName()] = append(out[mf.GetName()], s)
		}
	}
	return out
}
