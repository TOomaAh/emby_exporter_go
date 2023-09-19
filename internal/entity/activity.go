package entity

import (
	"time"
)

type Activity struct {
	Items            []ActivityItem `json:"Items"`
	TotalRecordCount int            `json:"TotalRecordCount"`
}

type ActivityItem struct {
	ID       int       `json:"Id"`
	Name     string    `json:"Name"`
	Type     string    `json:"Type"`
	Date     time.Time `json:"Date"`
	UserID   string    `json:"UserId"`
	Severity string    `json:"Severity"`
}
type ActivityMetric struct {
	ID       int
	Name     string
	Type     string
	Severity string
	Date     time.Time
}
