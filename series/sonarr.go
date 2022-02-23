package series

type Sonarr struct {
	*Indexer
}

func NewSonarr(url string, token string) *Sonarr {
	return &Sonarr{
		&Indexer{
			Url:   url,
			Token: token,
		},
	}
}

func (m Sonarr) GetTodayTVShow() string {
	return "oc"
}

func (m Sonarr) GetHistory() string {
	return "oc"
}
