package metrics

import (
	"TOomaAh/emby_exporter_go/pkg/emby"
	"TOomaAh/emby_exporter_go/pkg/logger"
	"sync"

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

func NewLibraryCollector(server *emby.Server, logger logger.Interface) prometheus.Collector {
	return &LibraryCollector{
		server:  server,
		library: prometheus.NewDesc("emby_media_item", "All Media Item", libraryValue, nil),
		logger:  logger,
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

	var wg sync.WaitGroup

	for _, l := range libraries.LibraryItem {
		wg.Go(func() {
			librarySize, err := c.server.GetLibrarySize(l.ItemID, l.LibraryOptions.ContentType)

			if err != nil {
				c.logger.Error("Error while getting library size for %s: %s", l.Name, err)
				return
			}

			ch <- prometheus.MustNewConstMetric(
				c.library,
				prometheus.GaugeValue, float64(librarySize),
				l.ItemID,
				l.Name,
			)
		})
	}
	wg.Wait()
}
