package metrics

import (
	"TOomaAh/emby_exporter_go/pkg/emby"
	"TOomaAh/emby_exporter_go/pkg/logger"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	activityValue = []string{
		"id",
		"name",
		"type",
		"severity",
		"date",
	}
)

type ActivityCollector struct {
	server   *emby.Server
	activity *prometheus.Desc
	logger   logger.Interface
}

func NewActivityCollector(server *emby.Server, logger logger.Interface) prometheus.Collector {
	return &ActivityCollector{
		server:   server,
		activity: prometheus.NewDesc("emby_activity", "Activity log", activityValue, nil),
		logger:   logger,
	}
}

func (c *ActivityCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.activity
}

func (c *ActivityCollector) Collect(ch chan<- prometheus.Metric) {
	activities, err := c.server.GetActivity()

	if err != nil {
		c.logger.Error("Error while getting activity: %s", err)
		return
	}

	for _, a := range activities.Items {
		ch <- prometheus.MustNewConstMetric(
			c.activity,
			prometheus.GaugeValue, 1,
			strconv.Itoa(a.ID),
			a.Name,
			a.Type,
			a.Severity,
			a.Date.In(time.Local).Format("02/01/2006 15:04:05"),
		)
	}
}
