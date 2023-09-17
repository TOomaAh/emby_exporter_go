package emby

import (
	"log"
	"time"
)

type Alert struct {
	Items            []AlertItem `json:"Items"`
	TotalRecordCount int         `json:"TotalRecordCount"`
}

type AlertItem struct {
	ID            int       `json:"Id"`
	Overview      string    `json:"Overview,omitempty"`
	ShortOverview string    `json:"ShortOverview"`
	Type          string    `json:"Type"`
	Date          time.Time `json:"Date"`
	Severity      string    `json:"Severity"`
}

type AlertMetrics struct {
	ID            int
	Name          string
	Overview      string
	ShortOverview string
	Type          string
	Date          time.Time
	Severity      string
}

func (s *Server) GetAlert() *Alert {
	var alert Alert
	err := s.request("GET", "/System/ActivityLog/Entries?StartIndex=0&Limit=4&hasUserId=false", "", &alert)

	if err != nil {
		log.Println("Cannot get alert, maybe your server is unreachable")
		alert.Items = make([]AlertItem, 0)
		return &alert
	}

	return &alert
}
