package xpo

import (
	"local/validatekit"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// Report struct
type Report struct {
	ID        int64          `datastore:"-" goon:"id"`
	AuthorKey *datastore.Key `datastore:"-" goon:"parent"`
	Author    string         `json:"author"`
	Content   string         `json:"content"`
	Year      int16          `json:"year"`
	Month     int8           `json:"month"`
	Day       int8           `json:"day"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type ReportService struct {
	AppEngineService
}

type ReportCreationParams struct {
	Content string `json:"content" validate:"required"`
}

func NewReportService(c context.Context) *ReportService {
	s := new(ReportService)
	s.InitAppEngineService(c)
	return s
}

func (s *ReportService) RetriveAll() (reports []Report, err error) {
	limit := 10
	q := datastore.NewQuery("Report").Order("-CreatedAt").Limit(limit)
	reports = make([]Report, 0, limit)
	_, err = s.Goon.GetAll(q, &reports)
	return
}

func (s *ReportService) Create(xu *XUser, params ReportCreationParams) (report *Report, err error) {
	v := validatekit.NewValidate()
	err = v.Struct(params)
	if err != nil {
		return nil, err
	}

	report = &Report{Content: params.Content}

	now := time.Now()
	report.Author = xu.Name
	report.CreatedAt = now
	report.UpdatedAt = now
	report.AuthorKey = s.KeyOf(xu)

	err = s.Put(report)
	return
}
