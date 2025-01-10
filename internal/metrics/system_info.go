package metrics

import (
	"TOomaAh/emby_exporter_go/pkg/emby"
	"TOomaAh/emby_exporter_go/pkg/logger"
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
	logger     logger.Interface
}

func NewSystemInfoCollector(server *emby.Server, logger logger.Interface) *SystemInfoCollector {
	return &SystemInfoCollector{
		server:     server,
		serverInfo: prometheus.NewDesc("emby_system_info", "All Emby Info", serverInfoValue, nil),
		logger:     logger,
	}
}

func (c *SystemInfoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.serverInfo
}

func (c *SystemInfoCollector) Collect(ch chan<- prometheus.Metric) {
	systemInfo, err := c.server.GetServerInfo()

	if err != nil {
		c.logger.Error("Error while getting system info: %s", err)
		return
	}

	ch <- prometheus.MustNewConstMetric(
		c.serverInfo,
		prometheus.GaugeValue, 1,
		systemInfo.Version,
		normalizeAddress(systemInfo.WanAddress),
		normalizeAddress(systemInfo.LocalAddress),
		strconv.FormatBool(systemInfo.HasUpdateAvailable),
		strconv.FormatBool(systemInfo.HasPendingRestart),
	)
}

func normalizeAddress(address string) string {
	if address == "" {
		return "N/A"
	}
	return address
}
