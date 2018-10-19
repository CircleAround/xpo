package domain

import (
	"context"
	"local/gaekit"
	"local/xpo/entities"
	"local/xpo/store"

	"google.golang.org/appengine/datastore"
)

type ReportCreatorFactory struct {
	mrep *store.MonthlyReportOverviewRepository
	rrep *store.ReportRepository
}

func NewReportCreatorFactory(rrep *store.ReportRepository, mrep *store.MonthlyReportOverviewRepository) *ReportCreatorFactory {
	return &ReportCreatorFactory{
		rrep: rrep,
		mrep: mrep,
	}
}

func (f *ReportCreatorFactory) Create(c context.Context, xu *entities.XUser, report *entities.Report) (*ReportCreator, error) {
	if xu.ReportCount == 0 {
		rs, err := f.rrep.Search(c, store.ReportSearchParams{
			AuthorID: xu.ID,
		}, 0)

		if err != nil {
			return nil, err
		}

		xu.ReportCount = int64(len(rs))
	}

	ra := report.ReportedAt
	m, err := f.mrep.MontlyReportOverview(c, xu, ra.Year(), int(ra.Month()))
	if err != nil && err != datastore.ErrNoSuchEntity {
		return nil, err
	}
	m.DailyReportCounts[ra.Day()]++
	m.ReportCount++

	xu.ReportCount++

	return &ReportCreator{
		xu:  xu,
		m:   m,
		rep: report,
	}, nil
}

type ReportCreator struct {
	gaekit.AppEngineService

	xu  *entities.XUser
	rep *entities.Report
	m   *entities.MonthlyReportOverview
}

func (rc *ReportCreator) Create(c context.Context) error {
	// for idempotent
	oxu := entities.XUser{ID: rc.xu.ID}
	err := rc.Get(c, &oxu)
	if err != nil {
		return err
	}

	if oxu.ReportCount == rc.xu.ReportCount {
		// already put
		return nil
	}

	return rc.PutAll(c, []interface{}{rc.rep, rc.m, rc.xu})
}
