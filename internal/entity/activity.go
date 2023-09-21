package entity

import (
	"time"
)

type Activity struct {
	Items            []ActivityItem `json:"Items"`
	TotalRecordCount int            `json:"TotalRecordCount"`
}

type ActivityItem struct {
	Name     string    `json:"Name"`
	Type     string    `json:"Type"`
	UserID   string    `json:"UserId"`
	Severity string    `json:"Severity"`
	Date     time.Time `json:"Date"`
	ID       int       `json:"Id"`
}
type ActivityMetric struct {
	ID       int
	Name     string
	Type     string
	Severity string
	Date     time.Time
}
