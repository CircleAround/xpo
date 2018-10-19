package store

import (
	"context"
	"local/gaekit"
	"local/xpo/entities"
	"time"

	"google.golang.org/appengine/datastore"
)

type ReportRepository struct {
	gaekit.AppEngineService
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
