package store

import (
	"context"
	"fmt"
	"local/gaekit"
	"local/xpo/entities"
	"time"

	"google.golang.org/appengine/datastore"
)

type ReportRepository struct {
	gaekit.DatastoreAccessObject
}

type ReportSearchParams struct {
	AuthorID       string
	ReportedAtFrom time.Time
	ReportedAtTo   time.Time
}

func NewReportRepository() *ReportRepository {
	return new(ReportRepository)
}

func (s *ReportRepository) Search(c context.Context, p ReportSearchParams, limit int) (reports []entities.Report, err error) {
	q := datastore.NewQuery("Report").Order("-ReportedAt")
	if limit != 0 {
		q = q.Limit(limit)
	}

	if p.AuthorID != "" {
		q = q.Filter("AuthorID=", p.AuthorID)
	}

	if !p.ReportedAtFrom.IsZero() {
		q = q.Filter("ReportedAt>=", p.ReportedAtFrom)
	}

	if !p.ReportedAtTo.IsZero() {
		q = q.Filter("ReportedAt<", p.ReportedAtTo)
	}

	reports = make([]entities.Report, 0, limit)
	_, err = s.Goon(c).GetAll(q, &reports)
	return
}

func (f *ReportRepository) Create(c context.Context, xu *entities.XUser, report *entities.Report) error {
	if xu.ReportCount == 0 {
		rs, err := f.Search(c, ReportSearchParams{
			AuthorID: xu.ID,
		}, 0)

		if err != nil {
			return err
		}

		xu.ReportCount = int64(len(rs))
	}

	return f.RunInXGTransaction(c, func(c context.Context) error {
		ra := report.ReportedAt
		m, err := f.MontlyReportOverview(c, xu, ra.Year(), int(ra.Month()))
		if err != nil && err != datastore.ErrNoSuchEntity {
			return err
		}
		m.DailyReportCounts[ra.Day()]++
		m.ReportCount++

		xu.ReportCount++

		es, err := f.buildLanguageEntities(c, xu, []string{}, report.Languages)
		if err != nil {
			return err
		}
		return f.PutMulti(c, append(es, report, m, xu))
	})
}

func (f *ReportRepository) Update(c context.Context, xu *entities.XUser, report *entities.Report) error {
	return f.RunInXGTransaction(c, func(c context.Context) error {
		br := &entities.Report{ID: report.ID, AuthorKey: report.AuthorKey}
		err := f.Get(c, br)
		if err != nil {
			return err
		}

		es, err := f.buildLanguageEntities(c, xu, br.Languages, report.Languages)
		if err != nil {
			return err
		}

		return f.PutMulti(c, append(es, report))
	})
}

func (s *ReportRepository) MontlyReportOverview(c context.Context, xu *entities.XUser, year, month int) (*entities.MonthlyReportOverview, error) {
	m := s.NewMonthlyReportOverview(s.KeyOf(c, xu), year, month)
	err := s.Get(c, m)
	return m, err
}

func (s *ReportRepository) NewMonthlyReportOverview(ak *datastore.Key, y, m int) *entities.MonthlyReportOverview {
	return &entities.MonthlyReportOverview{
		AuthorKey:         ak,
		Year:              y,
		Month:             m,
		ID:                fmt.Sprintf("%d/%0d", y, m),
		DailyReportCounts: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, // DailyReportCounts[0] is unusedvalue
	}
}

func (f *ReportRepository) buildLanguageEntities(c context.Context, xu *entities.XUser, b []string, a []string) ([]interface{}, error) {
	es, err := newLanguageLabelToCounterUpdatingDao().BuildEntities(c, b, a)
	if err != nil {
		return nil, err
	}

	ues, err := newLanguageLabelToXUserCounterUpdatingDao(xu).BuildEntities(c, b, a)
	if err != nil {
		return nil, err
	}

	return append(es, ues...), nil
}

func newLanguageLabelToCounterUpdatingDao() *gaekit.LabelToCounterUpdatingDao {
	var lng *entities.Language
	return &gaekit.LabelToCounterUpdatingDao{
		NewFunc: func(c context.Context, label string) interface{} {
			lng = &entities.Language{Name: label, ReportCount: 0}
			return lng
		},
		IncrementFunc: func() {
			lng.ReportCount++
		},
		DecrementFunc: func() {
			lng.ReportCount--
		},
	}
}

func newLanguageLabelToXUserCounterUpdatingDao(xu *entities.XUser) *gaekit.LabelToCounterUpdatingDao {
	var lng *entities.XUserLanguage
	dao := &gaekit.LabelToCounterUpdatingDao{}
	dao.NewFunc = func(c context.Context, label string) interface{} {
		lng = &entities.XUserLanguage{Name: label, ReportCount: 0, XUserKey: dao.KeyOf(c, xu)}
		return lng
	}
	dao.IncrementFunc = func() {
		lng.ReportCount++
	}
	dao.DecrementFunc = func() {
		lng.ReportCount--
	}
	return dao
}
