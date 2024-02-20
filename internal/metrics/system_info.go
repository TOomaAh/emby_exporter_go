package metrics

import (
	"TOomaAh/emby_exporter_go/pkg/emby"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	serverInfoValue = []string{"version",
		"wanAdress",
		"localAdress",
		"hasUpdateAvailable",
		"hasPendingRestart",
	}
)

type SystemInfoCollector struct {
	server     *emby.Server
	serverInfo *prometheus.Desc
}

func NewSystemInfoCollector(server *emby.Server) *SystemInfoCollector {
	return &SystemInfoCollector{
		server:     server,
		serverInfo: prometheus.NewDesc("emby_system_info", "All Emby Info", serverInfoValue, nil),
	}
}

func (c *SystemInfoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.serverInfo
}

func (c *SystemInfoCollector) Collect(ch chan<- prometheus.Metric) {
	systemInfo, err := c.server.GetServerInfo()

	if err != nil {
		ch <- prometheus.MustNewConstMetric(
			c.serverInfo,
			prometheus.GaugeValue, 1,
			"0.0.0",
			"empty",
			"empty",
			"false",
			"false",
		)
		return
	}

	if systemInfo.WanAddress == "" {
		systemInfo.WanAddress = "N/A"
	}

	if systemInfo.LocalAddress == "" {
		systemInfo.LocalAddress = "N/A"
	}

	ch <- prometheus.MustNewConstMetric(
		c.serverInfo,
		prometheus.GaugeValue, 1,
		systemInfo.Version,
		systemInfo.WanAddress,
		systemInfo.LocalAddress,
		strconv.FormatBool(systemInfo.HasUpdateAvailable),
		strconv.FormatBool(systemInfo.HasPendingRestart),
	)
}
