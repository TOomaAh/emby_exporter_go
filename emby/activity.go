package emby

import (
	"log"
	"time"
)

type Activity struct {
	Items []struct {
		ID       int       `json:"Id"`
		Name     string    `json:"Name"`
		Type     string    `json:"Type"`
		Date     time.Time `json:"Date"`
		UserID   string    `json:"UserId"`
		Severity string    `json:"Severity"`
	} `json:"Items"`
	TotalRecordCount int `json:"TotalRecordCount"`
}

type ActivityMetric struct {
	ID       int
	Name     string
	Type     string
	Severity string
	Date     time.Time
}

func (s *Server) GetActivity() *Activity {
	var activity Activity
	err := s.request("GET", "/System/ActivityLog/Entries?StartIndex=0&Limit=7", "", &activity)

	if err != nil {
		log.Println("Cannot get activity, maybe your server is unreachable")
		return &activity
	}

	return &activity
}
