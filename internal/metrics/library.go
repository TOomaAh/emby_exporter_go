package metrics

import (
	"TOomaAh/emby_exporter_go/pkg/emby"
	"TOomaAh/emby_exporter_go/pkg/logger"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	libraryValue = []string{"id", "name"}
)

type LibraryCollector struct {
	server  *emby.Server
	library *prometheus.Desc
	logger  logger.Interface
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
		c.logger.Error("Error while getting library: %s", err)
		return
	}

	for _, l := range libraries.LibraryItem {
		librarySize, err := c.server.GetLibrarySize(l.ItemID, l.LibraryOptions.ContentType)

		if err != nil {
			c.logger.Error("Error while getting library size: %s", err)
			return
		}

		ch <- prometheus.MustNewConstMetric(
			c.library,
			prometheus.GaugeValue, float64(librarySize),
			l.ItemID,
			l.Name,
		)
	}
}
