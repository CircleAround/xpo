package store

import (
	"context"
	"fmt"
	"local/gaekit"
	"local/xpo/entities"
	"time"

	funk "github.com/thoas/go-funk"
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

		es := []interface{}{}
		for _, value := range report.Languages {
			l := &entities.Language{Name: value, ReportCount: 0}
			err := f.Get(c, l)
			if err != nil && err != datastore.ErrNoSuchEntity {
				return err
			}
			l.ReportCount++
			es = append(es, l)
		}
		return f.PutMulti(c, append(es, report, m, xu))
	})
}

func (f *ReportRepository) Update(c context.Context, report *entities.Report) error {
	return f.RunInXGTransaction(c, func(c context.Context) error {
		br := &entities.Report{ID: report.ID, AuthorKey: report.AuthorKey}
		err := f.Get(c, br)
		if err != nil {
			return err
		}

		es := []interface{}{}
		for _, value := range report.Languages {
			if funk.Contains(br.Languages, value) {
				continue
			}

			l := &entities.Language{Name: value, ReportCount: 0}
			err := f.Get(c, l)
			if err != nil && err != datastore.ErrNoSuchEntity {
				return err
			}
			l.ReportCount++
			es = append(es, l)
		}

		for _, value := range br.Languages {
			if funk.Contains(report.Languages, value) {
				continue
			}

			l := &entities.Language{Name: value}
			err := f.Get(c, l)
			if err != nil {
				return err
			}
			l.ReportCount--
			es = append(es, l)
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
