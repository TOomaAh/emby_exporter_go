package series

import "time"

type SonarrSchedule struct {
	SeriesID           int       `json:"seriesId"`
	EpisodeFileID      int       `json:"episodeFileId"`
	SeasonNumber       int       `json:"seasonNumber"`
	EpisodeNumber      int       `json:"episodeNumber"`
	Title              string    `json:"title"`
	AirDate            string    `json:"airDate"`
	AirDateUtc         time.Time `json:"airDateUtc"`
	Overview           string    `json:"overview"`
	HasFile            bool      `json:"hasFile"`
	Monitored          bool      `json:"monitored"`
	SceneEpisodeNumber int       `json:"sceneEpisodeNumber"`
	SceneSeasonNumber  int       `json:"sceneSeasonNumber"`
	TvDbEpisodeID      int       `json:"tvDbEpisodeId"`
	Series             struct {
		TvdbID           int       `json:"tvdbId"`
		TvRageID         int       `json:"tvRageId"`
		ImdbID           string    `json:"imdbId"`
		Title            string    `json:"title"`
		CleanTitle       string    `json:"cleanTitle"`
		Status           string    `json:"status"`
		Overview         string    `json:"overview"`
		AirTime          string    `json:"airTime"`
		Monitored        bool      `json:"monitored"`
		QualityProfileID int       `json:"qualityProfileId"`
		SeasonFolder     bool      `json:"seasonFolder"`
		LastInfoSync     time.Time `json:"lastInfoSync"`
		Runtime          int       `json:"runtime"`
		Images           []struct {
			CoverType string `json:"coverType"`
			URL       string `json:"url"`
		} `json:"images"`
		SeriesType        string    `json:"seriesType"`
		Network           string    `json:"network"`
		UseSceneNumbering bool      `json:"useSceneNumbering"`
		TitleSlug         string    `json:"titleSlug"`
		Path              string    `json:"path"`
		Year              int       `json:"year"`
		FirstAired        time.Time `json:"firstAired"`
		QualityProfile    struct {
			Value struct {
				Name    string `json:"name"`
				Allowed []struct {
					ID     int    `json:"id"`
					Name   string `json:"name"`
					Weight int    `json:"weight"`
				} `json:"allowed"`
				Cutoff struct {
					ID     int    `json:"id"`
					Name   string `json:"name"`
					Weight int    `json:"weight"`
				} `json:"cutoff"`
				ID int `json:"id"`
			} `json:"value"`
			IsLoaded bool `json:"isLoaded"`
		} `json:"qualityProfile"`
		Seasons []struct {
			SeasonNumber int  `json:"seasonNumber"`
			Monitored    bool `json:"monitored"`
		} `json:"seasons"`
		ID int `json:"id"`
	} `json:"series"`
	Downloading bool `json:"downloading"`
	ID          int  `json:"id"`
}

type SonarrHistory struct {
	Page          int    `json:"page"`
	PageSize      int    `json:"pageSize"`
	SortKey       string `json:"sortKey"`
	SortDirection string `json:"sortDirection"`
	TotalRecords  int    `json:"totalRecords"`
	Records       []struct {
		EpisodeID   int    `json:"episodeId"`
		SeriesID    int    `json:"seriesId"`
		SourceTitle string `json:"sourceTitle"`
		Language    struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"language"`
		Quality struct {
			Quality struct {
				ID         int    `json:"id"`
				Name       string `json:"name"`
				Source     string `json:"source"`
				Resolution int    `json:"resolution"`
			} `json:"quality"`
			Revision struct {
				Version  int  `json:"version"`
				Real     int  `json:"real"`
				IsRepack bool `json:"isRepack"`
			} `json:"revision"`
		} `json:"quality"`
		QualityCutoffNotMet  bool      `json:"qualityCutoffNotMet"`
		LanguageCutoffNotMet bool      `json:"languageCutoffNotMet"`
		Date                 time.Time `json:"date"`
		DownloadID           string    `json:"downloadId"`
		EventType            string    `json:"eventType"`
		ID                   int       `json:"id"`
		Data                 struct {
			Indexer            string    `json:"indexer"`
			NzbInfoURL         string    `json:"nzbInfoUrl"`
			ReleaseGroup       string    `json:"releaseGroup"`
			Age                string    `json:"age"`
			AgeHours           string    `json:"ageHours"`
			AgeMinutes         string    `json:"ageMinutes"`
			PublishedDate      time.Time `json:"publishedDate"`
			DownloadClient     string    `json:"downloadClient"`
			DownloadClientName string    `json:"downloadClientName"`
			Size               string    `json:"size"`
			DownloadURL        string    `json:"downloadUrl"`
			GUID               string    `json:"guid"`
			TvdbID             string    `json:"tvdbId"`
			TvRageID           string    `json:"tvRageId"`
			Protocol           string    `json:"protocol"`
			PreferredWordScore string    `json:"preferredWordScore"`
			TorrentInfoHash    string    `json:"torrentInfoHash"`
		} `json:"data,omitempty"`
	} `json:"records"`
}

type SonarrEpisode struct {
	SeriesID                 int    `json:"seriesId"`
	EpisodeFileID            int    `json:"episodeFileId"`
	SeasonNumber             int    `json:"seasonNumber"`
	EpisodeNumber            int    `json:"episodeNumber"`
	Title                    string `json:"title"`
	AirDate                  string `json:"airDate"`
	AirDateUTC               string `json:"airDateUtc"`
	HasFile                  bool   `json:"hasFile"`
	Monitored                bool   `json:"monitored"`
	AbsoluteEpisodeNumber    int    `json:"absoluteEpisodeNumber"`
	UnverifiedSceneNumbering bool   `json:"unverifiedSceneNumbering"`
	ID                       int    `json:"id"`
}
