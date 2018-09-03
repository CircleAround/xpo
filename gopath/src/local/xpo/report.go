package xpo

import (
	"local/apikit"
	"net/http"
	"time"

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

func NewReportService(r *http.Request) *ReportService {
	s := new(ReportService)
	s.InitAppEngineService(r)
	return s
}

func (s *ReportService) RetriveAll() (reports []Report, err error) {
	limit := 10
	q := datastore.NewQuery("Report").Order("-CreatedAt").Limit(limit)
	reports = make([]Report, 0, limit)
	_, err = s.Goon.GetAll(q, &reports)
	return
}

func (s *ReportService) Create(xu *XUser, report *Report) (err error) {
	verr := apikit.NewValidationError()
	if report.Content == "" {
		verr.PushOne("content", apikit.Required)
	}

	if verr.HasItem() {
		return verr
	}

	now := time.Now()
	report.Author = xu.Name
	report.CreatedAt = now
	report.UpdatedAt = now
	report.AuthorKey = s.KeyOf(xu)

	err = s.Put(report)
	return
}
