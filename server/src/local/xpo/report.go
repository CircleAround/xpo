package xpo

import (
	"local/gaekit"
	"local/the_time"
	"local/validatekit"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// Report struct
type Report struct {
	ID          int64          `json:"id" datastore:"-" goon:"id"`
	AuthorKey   *datastore.Key `json:"-" datastore:"-" goon:"parent" validate:"required"`
	AuthorID    string         `json:"authorId" validate:"required"`
	Author      string         `json:"author" validate:"required"`
	Content     string         `json:"content" validate:"required"`
	ContentType string         `json:"content_type" validate:"required"`
	ReportedAt  time.Time      `json:"repoorted_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type ReportService struct {
	gaekit.AppEngineService
	timeProvider the_time.Provider
}

type ReportCreationParams struct {
	Content     string    `json:"content" validate:"required"`
	ContentType string    `json:"content_type" validate:"required"`
	ReportedAt  time.Time `json:"reported_at"`
}

type ReportUpdatingParams struct {
	ReportCreationParams
	ID int64 `json:"id" validate:"required"`
}

func NewReportService(c context.Context) *ReportService {
	return NewReportServiceWithTheTime(c, the_time.Real())
}

func NewReportServiceWithTheTime(c context.Context, tp the_time.Provider) *ReportService {
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

func (s *ReportService) SearchBy(authorID string, year int, month int, day int) (reports []Report, err error) {
	limit := 10

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return
	}
	from := time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc)
	to := from.AddDate(0, 0, 1)

	q := datastore.NewQuery("Report").Order("-ReportedAt").Limit(limit).
		Filter("ReportedAt>=", from).
		Filter("ReportedAt<", to).
		Filter("AuthorID=", authorID)
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
	report.ContentType = params.ContentType
	report.Author = xu.Name
	report.AuthorID = xu.ID
	report.AuthorKey = s.KeyOf(xu)

	now := s.now()
	if params.ReportedAt.IsZero() {
		report.ReportedAt = now
	} else {
		report.ReportedAt = params.ReportedAt
	}
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
	report.ContentType = params.ContentType
	report.Author = xu.Name
	if !params.ReportedAt.IsZero() {
		report.ReportedAt = params.ReportedAt
	}
	report.UpdatedAt = s.now()

	err = s.Put(report)
	return
}

func (s *ReportService) now() time.Time {
	return s.timeProvider.Now()
}
