package xpo

import (
	"local/gaekit"
	"local/timekit"
	"local/validatekit"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// Report struct
type Report struct {
	ID        int64          `datastore:"-" goon:"id"`
	AuthorKey *datastore.Key `datastore:"-" goon:"parent" validate:"required"`
	Author    string         `json:"author"`
	Content   string         `json:"content" validate:"required"`
	Year      int16          `json:"year"`
	Month     int8           `json:"month"`
	Day       int8           `json:"day"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type ReportService struct {
	gaekit.AppEngineService
	timeProvider timekit.TimeProvider
}

type ReportCreationParams struct {
	Content string `json:"content" validate:"required"`
}

type ReportUpdatingParams struct {
	ReportCreationParams
	ID int64 `json:"id" validate:"required"`
}

func NewReportService(c context.Context) *ReportService {
	return NewReportServiceWithTimeProvider(c, new(timekit.RealTimeProvider))
}

func NewReportServiceWithTimeProvider(c context.Context, tp timekit.TimeProvider) *ReportService {
	s := new(ReportService)
	s.InitAppEngineService(c)
	s.timeProvider = tp
	return s
}

func (s *ReportService) RetriveAll() (reports []Report, err error) {
	limit := 10
	q := datastore.NewQuery("Report").Order("-CreatedAt").Limit(limit)
	reports = make([]Report, 0, limit)
	_, err = s.Goon.GetAll(q, &reports)
	return
}

func (s *ReportService) Create(xu XUser, params ReportCreationParams) (report *Report, err error) {
	err = validatekit.NewValidate().Struct(params)
	if err != nil {
		return nil, err
	}

	report = &Report{}
	report.Content = params.Content
	report.Author = xu.Name
	report.AuthorKey = s.KeyOf(xu)

	now := s.now()
	report.CreatedAt = now
	report.UpdatedAt = now

	err = s.Put(report)
	return
}

func (s *ReportService) Update(xu XUser, params ReportUpdatingParams) (report *Report, err error) {
	err = validatekit.NewValidate().Struct(params)
	if err != nil {
		return nil, err
	}

	report = &Report{ID: params.ID, AuthorKey: s.KeyOf(xu)}
	err = s.Get(report)
	if err != nil {
		return nil, err
	}

	report.Content = params.Content
	report.Author = xu.Name

	report.UpdatedAt = s.now()

	err = s.Put(report)
	return
}

func (s *ReportService) now() time.Time {
	return s.timeProvider.Now()
}
