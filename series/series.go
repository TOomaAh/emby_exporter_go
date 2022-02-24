package series

import "TOomaAh/emby_exporter_go/conf"

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
	Url   string
	Token string
}

type SeriesInterface interface {
	GetTodayEpisodes() *[]Episode
	GetHistory() *[]Episode
	makeRequest(method string, path string, body string) ([]byte, error)
}

func NewSeriesFromConf(conf *conf.Config) SeriesInterface {
	if conf.Series.Medusa.Url != "" {
		return NewMedusa(conf.Series.Medusa.Url, conf.Series.Medusa.Token)
	} else if conf.Series.Sonarr.Url != "" {
		return NewSonarr(conf.Series.Sonarr.Url, conf.Series.Sonarr.Token)
	}
	return nil
}
