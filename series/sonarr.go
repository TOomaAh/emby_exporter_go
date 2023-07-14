package series

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Sonarr struct {
	*Indexer
}

const (
	sonarrHistory      = "/api/v3/history?page=1&pageSize=10&sortDirection=descending&sortKey=date"
	sonarrEpisodeTitle = "/api/v3/episode?episodeIds=%d"
)

var (
	defaultSonarrEpisode = &SonarrEpisode{
		SeriesID:                 0,
		EpisodeFileID:            0,
		SeasonNumber:             0,
		EpisodeNumber:            0,
		Title:                    "N/A",
		AirDate:                  "N/A",
		AirDateUTC:               "N/A",
		HasFile:                  false,
		Monitored:                false,
		AbsoluteEpisodeNumber:    0,
		UnverifiedSceneNumbering: false,
		ID:                       0,
	}
)

func NewSonarr(url string, token string) *Sonarr {
	return &Sonarr{
		&Indexer{
			Url:   url,
			Token: token,
		},
	}
}

func (s *Sonarr) GetTodayEpisodes() []*Episode {
	currentTime := time.Now()
	today := currentTime.Format("2006-01-02")
	resp, err := s.makeRequest("GET", "/api/calendar?start="+today, "")

	if err != nil {
		log.Println("Sonarr - GetTodayEpisodes" + err.Error())
		return nil
	}

	var sonarrSchedule []SonarrSchedule
	err = json.Unmarshal(resp, &sonarrSchedule)
	if err != nil {
		log.Println("Sonarr - GetTodayEpisodes" + err.Error())
		return nil
	}
	var episodes []*Episode
	for _, s := range sonarrSchedule {
		episodes = append(episodes, &Episode{
			Name:    s.Title,
			AirDate: s.AirDate + " " + s.Series.AirTime,
			Season:  "S" + formatSeriesOrEpisode(s.SeasonNumber) + "E" + formatSeriesOrEpisode(s.EpisodeNumber),
		})
	}

	return episodes
}

func (s *Sonarr) getEpisodeTitle(episodeID int) *SonarrEpisode {

	resp, err := s.makeRequest("GET", fmt.Sprintf(sonarrEpisodeTitle, episodeID), "")

	if err != nil {
		log.Println("Sonarr - getEpisodeTitle" + err.Error())
		return defaultSonarrEpisode
	}

	var sonarrEpisode []SonarrEpisode
	err = json.Unmarshal(resp, &sonarrEpisode)
	if err != nil {
		log.Println("Sonarr - getEpisodeTitle" + err.Error())
		return defaultSonarrEpisode
	}

	for _, s := range sonarrEpisode {
		return &SonarrEpisode{
			SeriesID:                 s.SeriesID,
			EpisodeFileID:            s.EpisodeFileID,
			SeasonNumber:             s.SeasonNumber,
			EpisodeNumber:            s.EpisodeNumber,
			Title:                    s.Title,
			AirDate:                  s.AirDate,
			AirDateUTC:               s.AirDateUTC,
			HasFile:                  s.HasFile,
			Monitored:                s.Monitored,
			AbsoluteEpisodeNumber:    s.AbsoluteEpisodeNumber,
			UnverifiedSceneNumbering: s.UnverifiedSceneNumbering,
			ID:                       s.ID,
		}
	}
	return defaultSonarrEpisode
}

func (s *Sonarr) GetHistory() []*Episode {
	resp, err := s.makeRequest("GET", sonarrHistory, "")
	if err != nil {
		return nil
	}

	var sonarrHistory SonarrHistory
	err = json.Unmarshal(resp, &sonarrHistory)
	if err != nil {
		return nil
	}

	var episodes []*Episode
	for _, sh := range sonarrHistory.Records {
		//Verifier si l'ID n'est pas déjà présent dans le tableau
		time := sh.Date.UTC().Format("2006-01-02")
		episodeInfo := s.getEpisodeTitle(sh.EpisodeID)
		episodes = append(episodes, &Episode{
			Name:    episodeInfo.Title,
			AirDate: time,
			Season:  "S" + formatSeriesOrEpisode(episodeInfo.SeasonNumber) + "E" + formatSeriesOrEpisode(episodeInfo.EpisodeNumber),
			Status:  sh.EventType,
		})
	}

	return episodes
}

func (s *Sonarr) makeRequest(method string, path string, body string) ([]byte, error) {
	req, _ := http.NewRequest(method, fmt.Sprintf("%s%s&apikey=%s", s.Url, path, s.Token), strings.NewReader(body))
	req.Header.Set("Application-Type", "application/json")

	if len(body) > 0 {
		bodybytes := []byte(body)
		buf := bytes.NewBuffer(bodybytes)
		req.Body = ioutil.NopCloser(buf)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Sonarr - makeRequest" + err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	respbody, _ := ioutil.ReadAll(resp.Body)
	return respbody, nil
}
