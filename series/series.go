package series

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
	GetAllTVShow() *[]Series
	makeRequest(method string, path string, body string) ([]byte, error)
}
