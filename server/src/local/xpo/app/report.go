package app

import (
	"local/gaekit"
	"local/the_time"
	"local/validatekit"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

// Report struct
type Report struct {
	ID          int64          `json:"id" datastore:"-" goon:"id"`
	AuthorKey   *datastore.Key `json:"-" datastore:"-" goon:"parent" validate:"required"`
	AuthorID    string         `json:"authorId" validate:"required"`
	Author      string         `json:"author" validate:"required"`
	Content     string         `json:"content" validate:"required,max=20000"`
	ContentType string         `json:"content_type" validate:"required"`
	ReportedAt  time.Time      `json:"reportedAt"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
}

type ReportService struct {
	gaekit.AppEngineService
	timeProvider the_time.Provider
}

type ReportCreationParams struct {
	Content     string    `json:"content" validate:"required"`
	ContentType string    `json:"content_type" validate:"required"`
	ReportedAt  time.Time `json:"reportedAt"`
}

type ReportUpdatingParams struct {
	ReportCreationParams
	ID int64 `json:"id" validate:"required"`
}

func NewReportService() *ReportService {
	return NewReportServiceWithTheTime(the_time.Real())
}

func NewReportServiceWithTheTime(tp the_time.Provider) *ReportService {
	s := new(ReportService)
	s.timeProvider = tp
	return s
}

func (s *ReportService) RetriveAll(c context.Context) (reports []Report, err error) {
	limit := 10
	q := datastore.NewQuery("Report").Order("-CreatedAt").Limit(limit)
	reports = make([]Report, 0, limit)
	_, err = s.Goon(c).GetAll(q, &reports)
	return
}

func (s *ReportService) SearchBy(c context.Context, authorID string, year int, month int, day int) (reports []Report, err error) {
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
	_, err = s.Goon(c).GetAll(q, &reports)
	return
}

func (s *ReportService) Find(c context.Context, uid string, id int64) (report *Report, err error) {
	xu := XUser{ID: uid}
	return s.FindByXUserAndID(c, xu, id)
}

func (s *ReportService) FindByXUserAndID(c context.Context, xu XUser, id int64) (report *Report, err error) {
	ak := s.KeyOf(c, xu)
	report = &Report{AuthorKey: ak, ID: id}
	err = s.Get(c, report)
	return
}

func (s *ReportService) Create(c context.Context, xu XUser, params ReportCreationParams) (report *Report, err error) {
	err = validatekit.NewValidate().Struct(params)
	if err != nil {
		return
	}

	report = &Report{}
	report.Content = params.Content
	report.ContentType = params.ContentType
	report.Author = xu.Name
	report.AuthorID = xu.ID
	report.AuthorKey = s.KeyOf(c, xu)

	now := s.now()
	if params.ReportedAt.IsZero() {
		report.ReportedAt = now
	} else {
		report.ReportedAt = params.ReportedAt
	}
	report.CreatedAt = now
	report.UpdatedAt = now

	err = s.Put(c, report)
	return
}

func (s *ReportService) Update(c context.Context, xu XUser, params ReportUpdatingParams) (report *Report, err error) {
	err = validatekit.NewValidate().Struct(params)
	if err != nil {
		return
	}

	report, err = s.FindByXUserAndID(c, xu, params.ID)
	if err != nil {
		return
	}

	report.Content = params.Content
	report.ContentType = params.ContentType
	report.Author = xu.Name
	if !params.ReportedAt.IsZero() {
		log.Infof(c, "update ReportedAt: %v", params.ReportedAt)
		report.ReportedAt = params.ReportedAt
	}
	report.UpdatedAt = s.now()

	err = s.Put(c, report)
	return
}

func (s *ReportService) now() time.Time {
	return s.timeProvider.Now()
}
