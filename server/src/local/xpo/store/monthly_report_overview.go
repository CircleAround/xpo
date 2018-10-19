package store

import (
	"context"
	"fmt"
	"local/gaekit"
	"local/xpo/entities"

	"google.golang.org/appengine/datastore"
)

type MonthlyReportOverviewRepository struct {
	gaekit.AppEngineService
}

func NewMonthlyReportOverviewRepository() *MonthlyReportOverviewRepository {
	return new(MonthlyReportOverviewRepository)
}

func (s *MonthlyReportOverviewRepository) MontlyReportOverview(c context.Context, xu *entities.XUser, year, month int) (*entities.MonthlyReportOverview, error) {
	m := NewMonthlyReportOverview(s.KeyOf(c, xu), year, month)
	err := s.Get(c, m)
	return m, err
}

func NewMonthlyReportOverview(ak *datastore.Key, y, m int) *entities.MonthlyReportOverview {
	return &entities.MonthlyReportOverview{
		AuthorKey:         ak,
		Year:              y,
		Month:             m,
		ID:                fmt.Sprintf("%d/%0d", y, m),
		DailyReportCounts: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, // DailyReportCounts[0] is unusedvalue
	}
}
