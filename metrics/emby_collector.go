package metrics

import (
	"TOomaAh/emby_exporter_go/emby"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type EmbyCollector struct {
	embyClient *emby.EmbyClient
	serverInfo *prometheus.Desc
	library    *prometheus.Desc
	sessions   *prometheus.Desc
	count      *prometheus.Desc
}

func NewEmbyCollector(e *emby.EmbyClient) *EmbyCollector {
	return &EmbyCollector{
		embyClient: e,
		serverInfo: prometheus.NewDesc("emby_system_info", "All Emby Info", []string{"version", "wanAdress", "localAdress", "hasUpdateAvailable", "hasPendingRestart"}, nil),
		library:    prometheus.NewDesc("emby_media_item", "All Media Item", []string{"name", "size"}, nil),
		sessions:   prometheus.NewDesc("emby_sessions", "All session", []string{"username", "client", "isPaused", "remoteEndPoint", "latitude", "longitude", "city", "region", "countryCode", "nowPlayingItemName", "nowPlayingItemType"}, nil),
		count:      prometheus.NewDesc("emby_sessions_count", "Session Count", []string{"count"}, nil),
	}
}

func (c *EmbyCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.serverInfo
	ch <- c.library
	ch <- c.sessions
	ch <- c.count
}

func (c *EmbyCollector) Collect(ch chan<- prometheus.Metric) {
	embyMetrics := c.embyClient.GetMetrics()
	ch <- prometheus.MustNewConstMetric(c.serverInfo, prometheus.GaugeValue, 1, embyMetrics.Info.Version, embyMetrics.Info.WanAddress, embyMetrics.Info.LocalAddress, strconv.FormatBool(embyMetrics.Info.HasUpdateAvailable), strconv.FormatBool(embyMetrics.Info.HasPendingRestart))
	for i, session := range embyMetrics.Sessions {
		ch <- prometheus.MustNewConstMetric(c.sessions, prometheus.GaugeValue, float64(i), session.Username, session.Client, strconv.FormatBool(session.IsPaused), session.RemoteEndPoint, strconv.FormatFloat(session.Latitude, 'f', 6, 64), strconv.FormatFloat(session.Longitude, 'f', 6, 64), session.City, session.Region, session.CountryCode, session.NowPlayingItemName, session.NowPlayingItemType)
	}
	for i, library := range embyMetrics.LibraryMetrics {
		ch <- prometheus.MustNewConstMetric(c.library, prometheus.GaugeValue, float64(i), library.Name, strconv.FormatInt(int64(library.Size), 10))
	}
	ch <- prometheus.MustNewConstMetric(c.count, prometheus.GaugeValue, 1, strconv.FormatInt(int64(len(embyMetrics.Sessions)), 10))
}
