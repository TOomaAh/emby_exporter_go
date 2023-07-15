package series

import (
	"TOomaAh/emby_exporter_go/conf"
	"net/http"
	"strconv"
)

type Series struct {
	name string
}

type Episode struct {
	Name    string
	AirDate string
	Season  string
	Status  string
}

type Indexer struct {
	Url    string
	Token  string
	client *http.Client
}

type SeriesInterface interface {
	GetTodayEpisodes() []*Episode
	GetHistory() []*Episode
	setToken(req *http.Request)
}

func NewSeriesFromConf(conf *conf.Config) SeriesInterface {
	if conf.Series.Medusa.Url != "" {
		return NewMedusa(conf.Series.Medusa.Url, conf.Series.Medusa.Token)
	} else if conf.Series.Sonarr.Url != "" {
		return NewSonarr(conf.Series.Sonarr.Url, conf.Series.Sonarr.Token)
	}
	return nil
}

func formatSeasonEpisodes(season int, episode int) string {
	return "S" + formatSeriesOrEpisode(season) + "E" + formatSeriesOrEpisode(episode)
}

func formatSeriesOrEpisode(episode int) string {
	if episode > 10 {
		return strconv.Itoa(episode)
	}
	return "0" + strconv.Itoa(episode)
}
