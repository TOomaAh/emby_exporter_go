package series

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Medusa struct {
	*Indexer
}

func NewMedusa(url string, token string) *Medusa {
	return &Medusa{
		&Indexer{
			Url:   url,
			Token: token,
		},
	}
}

func formatSeriesOrEpisode(episode int) string {
	if episode > 10 {
		return fmt.Sprintf("%d", episode)
	}
	return fmt.Sprintf("0%d", episode)
}

func (m Medusa) GetTodayEpisodes() *[]Episode {
	resp, err := m.makeRequest("GET", "/api/v2/schedule?category[]=later&category[]=today&paused=true", "")
	if err != nil {

		return nil
	}

	var medusaSchedule MedusaSchedule
	err = json.Unmarshal(resp, &medusaSchedule)
	if err != nil {
		return nil
	}

	var episodes []Episode
	for _, s := range medusaSchedule.Today {
		episodes = append(episodes, Episode{
			Name:    s.ShowName,
			AirDate: fmt.Sprintf("%s %s", s.Airdate, s.Airs),
			Season:  fmt.Sprintf("S%dE%s", s.Season, formatSeriesOrEpisode(s.Episode)),
		})
	}

	return &episodes
}

func (m Medusa) GetHistory() *[]Episode {
	resp, err := m.makeRequest("GET", "/api/v2/history?page=1&limit=25&sort[]={\"field\":\"date\",\"type\":\"desc\"}&filter={}", "")
	if err != nil {

		return nil
	}

	var medusaHistory []MedusaHistory
	err = json.Unmarshal(resp, &medusaHistory)
	if err != nil {
		return nil
	}

	var episodes []Episode
	for _, s := range medusaHistory {
		var year, month, day, hours, minutes string
		actionDate := strconv.FormatInt(s.ActionDate, 10)
		year = actionDate[0:4]
		month = actionDate[4:6]
		day = actionDate[6:8]
		hours = actionDate[8:10]
		minutes = actionDate[10:12]

		episodes = append(episodes, Episode{
			Name:    s.EpisodeTitle,
			AirDate: fmt.Sprintf("%s-%s-%s %s:%s", year, month, day, hours, minutes),
			Season:  fmt.Sprintf("S%sE%s", formatSeriesOrEpisode(s.Season), formatSeriesOrEpisode(s.Episode)),
			Status:  s.StatusName,
		})
	}

	return &episodes
}

func (m Medusa) GetAllTVShow() *[]Series {
	resp, err := m.makeRequest("GET", "/api/v2/series?sort=", "")

	if err != nil {
		return nil
	}

	var medusaSeries []MedusaSeries
	err = json.Unmarshal(resp, &medusaSeries)
	if err != nil {
		return nil
	}
	var series []Series
	for _, s := range medusaSeries {
		series = append(series, Series{
			name: s.Title,
		})
	}

	return &series

}

func (m Medusa) makeRequest(method string, path string, body string) ([]byte, error) {
	req, _ := http.NewRequest(method, fmt.Sprintf("%s%s", m.Url, path), strings.NewReader(body))
	req.Header.Set("x-api-key", m.Token)
	req.Header.Set("Application-Type", "application/json")

	if len(body) > 0 {
		bodybytes := []byte(body)
		buf := bytes.NewBuffer(bodybytes)
		req.Body = ioutil.NopCloser(buf)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respbody, _ := ioutil.ReadAll(resp.Body)
	return respbody, nil
}
