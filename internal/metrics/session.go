package metrics

import (
	"TOomaAh/emby_exporter_go/pkg/emby"
	"TOomaAh/emby_exporter_go/pkg/geoip"
	"TOomaAh/emby_exporter_go/pkg/logger"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var (
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
	server   *emby.Server
	sessions *prometheus.Desc
	count    *prometheus.Desc
	logger   logger.Interface
}

func NewSessionCollector(server *emby.Server) *SessionCollector {
	return &SessionCollector{
		server:   server,
		sessions: prometheus.NewDesc("emby_sessions", "All session", sessionsValue, nil),
		count:    prometheus.NewDesc("emby_sessions_count", "Session Count", []string{}, nil),
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

	db := geoip.GetGeoIPDatabase()

	var (
		tvshow      string = ""
		season      string = ""
		city        string = ""
		region      string = ""
		countryCode string = ""
		latitude    float64
		longitude   float64
	)

	count := 0
	for i, session := range *sessions {
		if !session.HasPlayMethod() {
			continue
		}

		count++
		if session.IsEpisode() {
			tvshow = session.NowPlayingItem.SeriesName
			season = session.NowPlayingItem.SeasonName
		}

		latitude, longitude = db.GetLocation(session.RemoteEndPoint)
		city = db.GetCity(session.RemoteEndPoint)
		region = db.GetRegion(session.RemoteEndPoint)
		countryCode = db.GetCountryCode(session.RemoteEndPoint)

		ch <- prometheus.MustNewConstMetric(
			c.sessions, prometheus.GaugeValue,
			float64(i), session.UserName,
			session.Client,
			strconv.FormatBool(session.PlayState.IsPaused),
			session.RemoteEndPoint,
			strconv.FormatFloat(latitude, 'f', 6, 64),
			strconv.FormatFloat(longitude, 'f', 6, 64),
			city,
			region,
			countryCode,
			session.NowPlayingItem.Name,
			tvshow,
			season,
			session.NowPlayingItem.Type,
			strconv.FormatInt(session.GetPercentPlayed(), 10),
			session.GetPlayMethod(),
			session.GetTranscodeReason(),
			session.GetDuration(session.GetRuntimeTick()),
			session.GetDuration(session.PlayState.PositionTicks),
			session.GetBitrate(),
		)
	}

	ch <- prometheus.MustNewConstMetric(
		c.count,
		prometheus.GaugeValue,
		float64(count))

}
