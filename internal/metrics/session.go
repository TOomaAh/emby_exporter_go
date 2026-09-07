package metrics

import (
	"TOomaAh/emby_exporter_go/pkg/emby"
	"TOomaAh/emby_exporter_go/pkg/geoip"
	"TOomaAh/emby_exporter_go/pkg/logger"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	sessionsBitrates = []string{
		"sessionId",
		"username",
		"client",
	}
	sessionsValue = []string{
		"username",
		"client",
		"isPaused",
		"remoteEndPoint",
		"latitude",
		"longitude",
		"city",
		"region",
		"countryCode",
		"nowPlayingItemName",
		"tvshow",
		"season",
		"nowPlayingItemType",
		"percentPlayback",
		"playMethod",
		"transcodeReason",
		"mediaDuration",
		"currentPlayTime",
		"bitrate",
	}
)

type SessionCollector struct {
	server           *emby.Server
	sessions         *prometheus.Desc
	count            *prometheus.Desc
	sessionsBitrates *prometheus.Desc
	logger           logger.Interface
	geoip            geoip.GeoIP
}

func NewSessionCollector(server *emby.Server, geoIP geoip.GeoIP, logger logger.Interface) prometheus.Collector {
	return &SessionCollector{
		server:           server,
		sessions:         prometheus.NewDesc("emby_sessions", "All session", sessionsValue, nil),
		sessionsBitrates: prometheus.NewDesc("emby_sessions_bitrate", "Session Bitrate in kilobits", sessionsBitrates, nil),
		count:            prometheus.NewDesc("emby_sessions_count", "Session Count", []string{}, nil),
		logger:           logger,
		geoip:            geoIP,
	}
}

func (c *SessionCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.sessions
}

func (c *SessionCollector) Collect(ch chan<- prometheus.Metric) {
	sessions, err := c.server.GetSessions()

	if err != nil {
		c.logger.Error("Error while getting sessions: %s", err)
		return
	}

	count := 0

	for i, session := range sessions {
		if !session.HasPlayMethod() {
			continue
		}

		count++
		tvshow, season := "", ""
		if session.IsEpisode() {
			tvshow = session.NowPlayingItem.SeriesName
			season = session.NowPlayingItem.SeasonName
		}

		var record *geoip.Record

		record, err := c.geoip.Lookup(session.RemoteEndPoint)

		if err != nil {
			c.logger.Error("Error getting GeoIP record for %s: %s", session.RemoteEndPoint, err)
			record = geoip.LocalRecord()
		}

		ch <- prometheus.MustNewConstMetric(
			c.sessions, prometheus.GaugeValue,
			float64(i), session.UserName,
			session.Client,
			strconv.FormatBool(session.PlayState.IsPaused),
			session.RemoteEndPoint,
			strconv.FormatFloat(record.Latitude, 'f', 6, 64),
			strconv.FormatFloat(record.Longitude, 'f', 6, 64),
			record.City,
			record.Region,
			record.CountryCode,
			session.NowPlayingItem.Name,
			tvshow,
			season,
			session.NowPlayingItem.Type,
			strconv.FormatInt(session.GetPercentPlayed(), 10),
			session.GetPlayMethod(),
			session.GetTranscodeReason(),
			session.GetDuration(session.GetRuntimeTick()),
			session.GetDuration(session.PlayState.PositionTicks),
			session.GetBitrateFormat(),
		)

		ch <- prometheus.MustNewConstMetric(
			c.sessionsBitrates,
			prometheus.GaugeValue,
			float64(session.GetBitrateValue()/1024),
			strconv.Itoa(i),
			session.UserName,
			session.Client,
		)
	}

	ch <- prometheus.MustNewConstMetric(
		c.count,
		prometheus.GaugeValue,
		float64(count))
}

func (c *SessionCollector) getRecord(ip string) *geoip.Record {
	record, err := c.geoip.Lookup(ip)

	if err != nil {
		return geoip.LocalRecord()
	}
	return record
}
