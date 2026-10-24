package metrics

import (
	"TOomaAh/emby_exporter_go/pkg/emby"

	"github.com/prometheus/client_golang/prometheus"
)

type UpCollector struct {
	server *emby.Server
	up     *prometheus.Desc
}

func NewUpCollector(server *emby.Server) prometheus.Collector {
	return &UpCollector{
		server: server,
		up:     prometheus.NewDesc("emby_up", "Whether the Emby server is reachable (1) or not (0)", nil, nil),
	}
}

func (c *UpCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
}

func (c *UpCollector) Collect(ch chan<- prometheus.Metric) {
	value := 0.0
	if c.server.IsReachable() {
		value = 1
	}
	ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, value)
}
