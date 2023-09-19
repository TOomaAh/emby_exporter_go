package series

import "time"

type MedusaSeries struct {
	ID struct {
		Tvdb  int    `json:"tvdb"`
		Imdb  string `json:"imdb"`
		Slug  string `json:"slug"`
		Trakt int    `json:"trakt"`
	} `json:"id"`
	Externals struct {
		Imdb   string `json:"imdb"`
		Tvdb   int    `json:"tvdb"`
		Tvrage int    `json:"tvrage"`
		Tvmaze int    `json:"tvmaze"`
		Tmdb   int    `json:"tmdb"`
		Trakt  int    `json:"trakt"`
	} `json:"externals"`
	Title           string `json:"title"`
	Name            string `json:"name"`
	Indexer         string `json:"indexer"`
	Network         string `json:"network"`
	Type            string `json:"type"`
	Status          string `json:"status"`
	Airs            string `json:"airs"`
	AirsFormatValid bool   `json:"airsFormatValid"`
	Language        string `json:"language"`
	ShowType        string `json:"showType"`
	ImdbInfo        struct {
		ImdbInfoID   int    `json:"imdbInfoId"`
		Indexer      int    `json:"indexer"`
		IndexerID    int    `json:"indexerId"`
		ImdbID       string `json:"imdbId"`
		Title        string `json:"title"`
		Year         int    `json:"year"`
		Akas         string `json:"akas"`
		Runtimes     int    `json:"runtimes"`
		Genres       string `json:"genres"`
		Countries    string `json:"countries"`
		CountryCodes string `json:"countryCodes"`
		Certificates string `json:"certificates"`
		Rating       string `json:"rating"`
		Votes        int    `json:"votes"`
		LastUpdate   int    `json:"lastUpdate"`
		Plot         string `json:"plot"`
	} `json:"imdbInfo"`
	Year struct {
		Start int `json:"start"`
	} `json:"year"`
	PrevAirDate time.Time   `json:"prevAirDate"`
	NextAirDate interface{} `json:"nextAirDate"`
	LastUpdate  string      `json:"lastUpdate"`
	Runtime     int         `json:"runtime"`
	Genres      []string    `json:"genres"`
	Rating      struct {
		Imdb struct {
			Rating string `json:"rating"`
			Votes  int    `json:"votes"`
		} `json:"imdb"`
	} `json:"rating"`
	Classification string `json:"classification"`
	Cache          struct {
		Poster string `json:"poster"`
		Banner string `json:"banner"`
	} `json:"cache"`
	Countries    []string `json:"countries"`
	CountryCodes []string `json:"countryCodes"`
	Plot         string   `json:"plot"`
	Config       struct {
		Location      string `json:"location"`
		RootDir       string `json:"rootDir"`
		LocationValid bool   `json:"locationValid"`
		Qualities     struct {
			Allowed   []int         `json:"allowed"`
			Preferred []interface{} `json:"preferred"`
		} `json:"qualities"`
		Paused               bool          `json:"paused"`
		AirByDate            bool          `json:"airByDate"`
		SubtitlesEnabled     bool          `json:"subtitlesEnabled"`
		DvdOrder             bool          `json:"dvdOrder"`
		SeasonFolders        bool          `json:"seasonFolders"`
		Anime                bool          `json:"anime"`
		Scene                bool          `json:"scene"`
		Sports               bool          `json:"sports"`
		Templates            bool          `json:"templates"`
		DefaultEpisodeStatus string        `json:"defaultEpisodeStatus"`
		Aliases              []interface{} `json:"aliases"`
		Release              struct {
			IgnoredWords         []interface{} `json:"ignoredWords"`
			RequiredWords        []interface{} `json:"requiredWords"`
			IgnoredWordsExclude  bool          `json:"ignoredWordsExclude"`
			RequiredWordsExclude bool          `json:"requiredWordsExclude"`
		} `json:"release"`
		AirdateOffset int      `json:"airdateOffset"`
		ShowLists     []string `json:"showLists"`
	} `json:"config"`
	XemNumbering           []interface{} `json:"xemNumbering"`
	SceneAbsoluteNumbering []interface{} `json:"sceneAbsoluteNumbering"`
	XemAbsoluteNumbering   []interface{} `json:"xemAbsoluteNumbering"`
	SceneNumbering         []interface{} `json:"sceneNumbering"`
}

type MedusaSchedule struct {
	Today []struct {
		Airdate      string    `json:"airdate"`
		Airs         string    `json:"airs"`
		LocalAirTime time.Time `json:"localAirTime"`
		EpName       string    `json:"epName"`
		EpPlot       string    `json:"epPlot"`
		Season       int       `json:"season"`
		Episode      int       `json:"episode"`
		EpisodeSlug  string    `json:"episodeSlug"`
		IndexerID    int       `json:"indexerId"`
		Indexer      int       `json:"indexer"`
		Network      string    `json:"network"`
		Paused       int       `json:"paused"`
		Quality      int       `json:"quality"`
		ShowSlug     string    `json:"showSlug"`
		ShowName     string    `json:"showName"`
		ShowStatus   string    `json:"showStatus"`
		Tvdbid       int       `json:"tvdbid"`
		Weekday      int       `json:"weekday"`
		Runtime      int       `json:"runtime"`
		Externals    struct {
			ImdbID   string `json:"imdb_id"`
			TvdbID   int    `json:"tvdb_id"`
			TvmazeID int    `json:"tvmaze_id"`
			TmdbID   int    `json:"tmdb_id"`
			TraktID  int    `json:"trakt_id"`
		} `json:"externals"`
	} `json:"today"`
}

type MedusaHistory struct {
	ID               int         `json:"id"`
	Series           string      `json:"series"`
	Status           int         `json:"status"`
	StatusName       string      `json:"statusName"`
	ActionDate       int64       `json:"actionDate"`
	Quality          int         `json:"quality"`
	Resource         string      `json:"resource"`
	Size             int         `json:"size"`
	ProperTags       string      `json:"properTags"`
	Season           int         `json:"season"`
	Episode          int         `json:"episode"`
	EpisodeTitle     string      `json:"episodeTitle"`
	ManuallySearched bool        `json:"manuallySearched"`
	InfoHash         interface{} `json:"infoHash"`
	ReleaseName      interface{} `json:"releaseName"`
	ReleaseGroup     string      `json:"releaseGroup"`
	FileName         string      `json:"fileName"`
	SubtitleLanguage interface{} `json:"subtitleLanguage"`
	ShowSlug         string      `json:"showSlug"`
	ShowTitle        string      `json:"showTitle"`
	ProviderType     interface{} `json:"providerType"`
	ClientStatus     interface{} `json:"clientStatus"`
	PartOfBatch      bool        `json:"partOfBatch"`
	Provider         struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		ImageName string `json:"imageName"`
	} `json:"provider,omitempty"`
}
