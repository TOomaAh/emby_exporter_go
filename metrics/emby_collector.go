package metrics

import (
	"TOomaAh/emby_exporter_go/emby"
	"fmt"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type EmbyCollector struct {
	embyClient *emby.EmbyClient
	serverInfo *prometheus.Desc
	library    *prometheus.Desc
	sessions   *prometheus.Desc
	count      *prometheus.Desc
	activity   *prometheus.Desc
	alert      *prometheus.Desc
}

func NewEmbyCollector(e *emby.EmbyClient) *EmbyCollector {
	return &EmbyCollector{
		embyClient: e,
		serverInfo: prometheus.NewDesc("emby_system_info", "All Emby Info", []string{"version", "wanAdress", "localAdress", "hasUpdateAvailable", "hasPendingRestart"}, nil),
		library:    prometheus.NewDesc("emby_media_item", "All Media Item", []string{"name"}, nil),
		sessions:   prometheus.NewDesc("emby_sessions", "All session", []string{"username", "client", "isPaused", "remoteEndPoint", "latitude", "longitude", "city", "region", "countryCode", "nowPlayingItemName", "tvshow", "season", "nowPlayingItemType", "percentPlayback", "playMethod"}, nil),
		count:      prometheus.NewDesc("emby_sessions_count", "Session Count", []string{}, nil),
		activity:   prometheus.NewDesc("emby_activity", "Activity log", []string{"id", "name", "type", "severity", "date"}, nil),
		alert:      prometheus.NewDesc("emby_alert", "Alert log", []string{"id", "name", "overview", "shortOverview", "type", "date", "severity"}, nil),
	}
}

func (c *EmbyCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.serverInfo
	ch <- c.library
	ch <- c.sessions
	ch <- c.count
	ch <- c.activity
}

func (c *EmbyCollector) Collect(ch chan<- prometheus.Metric) {
	embyMetrics := c.embyClient.GetMetrics()
	ch <- prometheus.MustNewConstMetric(c.serverInfo, prometheus.GaugeValue, 1, embyMetrics.Info.Version, embyMetrics.Info.WanAddress, embyMetrics.Info.LocalAddress, strconv.FormatBool(embyMetrics.Info.HasUpdateAvailable), strconv.FormatBool(embyMetrics.Info.HasPendingRestart))
	for i, session := range embyMetrics.Sessions {
		ch <- prometheus.MustNewConstMetric(c.sessions, prometheus.GaugeValue, float64(i), session.Username, session.Client, strconv.FormatBool(session.IsPaused), session.RemoteEndPoint, strconv.FormatFloat(session.Latitude, 'f', 6, 64), strconv.FormatFloat(session.Longitude, 'f', 6, 64), session.City, session.Region, session.CountryCode, session.NowPlayingItemName, session.TVShow, session.Season, session.NowPlayingItemType, strconv.FormatInt(session.PlaybackPercent, 10), session.PlayMethod)
	}

	for _, library := range embyMetrics.LibraryMetrics {
		ch <- prometheus.MustNewConstMetric(c.library, prometheus.GaugeValue, float64(library.Size), library.Name)
	}

	ch <- prometheus.MustNewConstMetric(c.count, prometheus.GaugeValue, float64(len(embyMetrics.Sessions)))

	for i, activity := range embyMetrics.Activity {
		ch <- prometheus.MustNewConstMetric(c.activity, prometheus.GaugeValue, float64(i), strconv.Itoa(activity.ID), activity.Name, activity.Type, activity.Severity, activity.Date.Format("02/01/2006 15:04:05"))
	}

	for i, alert := range embyMetrics.Alert {
		ch <- prometheus.MustNewConstMetric(c.alert, prometheus.GaugeValue, float64(i), fmt.Sprint(alert.ID), alert.Name, alert.Overview, alert.ShortOverview, alert.Type, alert.Date.Format("02/01/2006 15:04:05"), alert.Severity)
	}
}
