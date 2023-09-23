package entity

import (
	"time"
)

type Alert struct {
	Items            []AlertItem `json:"Items"`
	TotalRecordCount int         `json:"TotalRecordCount"`
}

type AlertItem struct {
	Overview      string    `json:"Overview,omitempty"`
	ShortOverview string    `json:"ShortOverview"`
	Type          string    `json:"Type"`
	Severity      string    `json:"Severity"`
	Date          time.Time `json:"Date"`
	ID            int       `json:"Id"`
}

type AlertMetrics struct {
	Name          string
	Overview      string
	ShortOverview string
	Type          string
	Severity      string
	Date          time.Time
	ID            int
}
