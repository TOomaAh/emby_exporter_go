package entity

import (
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
