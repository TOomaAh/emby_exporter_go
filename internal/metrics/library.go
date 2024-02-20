package metrics

import (
	"TOomaAh/emby_exporter_go/pkg/emby"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	libraryValue = []string{"name"}
)

type LibraryCollector struct {
	server  *emby.Server
	library *prometheus.Desc
}

func NewLibraryCollector(server *emby.Server) *LibraryCollector {
	return &LibraryCollector{
		server:  server,
		library: prometheus.NewDesc("emby_media_item", "All Media Item", libraryValue, nil),
	}
}

func (c *LibraryCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.library
}

func (c *LibraryCollector) Collect(ch chan<- prometheus.Metric) {
	libraries, err := c.server.GetLibrary()

	if err != nil {
		return
	}

	for _, l := range libraries.LibraryItem {
		ch <- prometheus.MustNewConstMetric(
			c.library,
			prometheus.GaugeValue, 1,
			l.Name,
		)
	}
}
