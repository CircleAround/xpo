package entities

import (
	"google.golang.org/appengine/datastore"
)

type MonthlyReportOverview struct {
	ID                string         `json:"id" datastore:"-" goon:"id" validate:"required"` // "YYYY/MM"
	AuthorKey         *datastore.Key `json:"-" datastore:"-" goon:"parent" validate:"required"`
	Year              int            `json:"year" validate:"required"`
	Month             int            `json:"month" validate:"required"`
	ReportCount       int64          `json:"reportCount" validate:"required"`
	DailyReportCounts []int          `json:"dailyReportCounts" validate:"required"`
}
