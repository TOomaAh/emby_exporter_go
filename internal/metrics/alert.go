package metrics

import (
	"TOomaAh/emby_exporter_go/pkg/emby"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	alertValue = []string{
		"id",
		"overview",
		"shortOverview",
		"type",
		"date",
		"severity",
	}
)

type AlertCollector struct {
	server *emby.Server
	alert  *prometheus.Desc
}

func NewAlertCollector(server *emby.Server) *AlertCollector {
	return &AlertCollector{
		server: server,
		alert:  prometheus.NewDesc("emby_alert", "Alert log", alertValue, nil),
	}
}

func (c *AlertCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.alert
}

func (c *AlertCollector) Collect(ch chan<- prometheus.Metric) {
	alerts, err := c.server.GetAlerts()

	if err != nil {
		return
	}

	for _, a := range alerts.Items {
		ch <- prometheus.MustNewConstMetric(
			c.alert,
			prometheus.GaugeValue, 1,
			strconv.Itoa(a.ID),
			a.Overview,
			a.ShortOverview,
			a.Type,
			a.Date.In(time.Local).Format("02/01/2006 15:04:05"),
			a.Severity,
		)
	}
}
