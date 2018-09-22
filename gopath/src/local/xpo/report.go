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
	ID        int64          `json:"id" datastore:"-" goon:"id"`
	AuthorKey *datastore.Key `json:"-" datastore:"-" goon:"parent" validate:"required"`
	AuthorID  string         `json:"author_id" validate:"required"`
	Author    string         `json:"author" validate:"required"`
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

func (s *ReportService) Find(uid string, id int64) (report *Report, err error) {
	xu := XUser{ID: uid}
	return s.FindByXUserAndID(xu, id)
}

func (s *ReportService) FindByXUserAndID(xu XUser, id int64) (report *Report, err error) {
	ak := s.KeyOf(xu)
	report = &Report{AuthorKey: ak, ID: id}
	err = s.Get(report)
	return
}

func (s *ReportService) Create(xu XUser, params ReportCreationParams) (report *Report, err error) {
	err = validatekit.NewValidate().Struct(params)
	if err != nil {
		return
	}

	report = &Report{}
	report.Content = params.Content
	report.Author = xu.Name
	report.AuthorID = xu.ID
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
		return
	}

	report, err = s.FindByXUserAndID(xu, params.ID)
	if err != nil {
		return
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
