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
	ID             int64          `json:"id" datastore:"-" goon:"id"`
	AuthorKey      *datastore.Key `json:"-" datastore:"-" goon:"parent" validate:"required"`
	AuthorID       string         `json:"authorId" validate:"required"`
	Author         string         `json:"author" validate:"required"`
	AuthorNickname string         `json:"authorNickname" validate:"required"`
	Content        string         `json:"content" validate:"required,max=20000" datastore:"Content,noindex"`
	ContentType    string         `json:"contentType" validate:"required"`
	ReportedAt     time.Time      `json:"reportedAt"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
}

type ReportService struct {
	gaekit.AppEngineService
	timeProvider the_time.Provider
}

type ReportCreationParams struct {
	Content     string    `json:"content" validate:"required"`
	ContentType string    `json:"contentType" validate:"required"`
	ReportedAt  time.Time `json:"reportedAt"`
}

type ReportUpdatingParams struct {
	ReportCreationParams
	ID int64 `json:"id" validate:"required"`
}

type ReportSerchParams struct {
	AuthorID       string
	ReportedAtFrom time.Time
	ReportedAtTo   time.Time
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
	limit := 30
	q := datastore.NewQuery("Report").Order("-ReportedAt").Limit(limit)
	reports = make([]Report, 0, limit)
	_, err = s.Goon(c).GetAll(q, &reports)
	return
}

func (s *ReportService) SearchBy(c context.Context, authorID string, year int, month int, day int) (reports []Report, err error) {
	limit := 30

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return
	}
	from := time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc)
	to := from.AddDate(0, 0, 1)

	return s.search(c, ReportSerchParams{
		AuthorID:       authorID,
		ReportedAtFrom: from,
		ReportedAtTo:   to,
	}, limit)
}

func (s *ReportService) SearchByAuthor(c context.Context, authorID string) (reports []Report, err error) {
	limit := 30
	return s.search(c, ReportSerchParams{
		AuthorID: authorID,
	}, limit)
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
	v := validatekit.NewValidate()
	err = v.Struct(params)
	if err != nil {
		return
	}

	report = &Report{}
	report.Content = params.Content
	report.ContentType = params.ContentType
	report.Author = xu.Name
	report.AuthorID = xu.ID
	report.AuthorKey = s.KeyOf(c, xu)
	report.AuthorNickname = xu.Nickname

	now := s.now()
	if params.ReportedAt.IsZero() {
		report.ReportedAt = now
	} else {
		report.ReportedAt = params.ReportedAt
	}
	report.CreatedAt = now
	report.UpdatedAt = now

	err = v.Struct(report)
	if err != nil {
		return
	}

	if xu.ReportCount == 0 {
		rs, err := s.search(c, ReportSerchParams{
			AuthorID: xu.ID,
		}, 0)

		if err != nil {
			return nil, err
		}

		xu.ReportCount = int64(len(rs))
	}

	xu.ReportCount++

	err = s.RunInXGTransaction(c, func(c context.Context) error {
		// for idempotent
		oxu := XUser{ ID: xu.ID }
		err = s.Get(c, &oxu)
		if err != nil {
			return err
		}
		if oxu.ReportCount == xu.ReportCount {
			// already put
			return nil
		}

		err = s.Put(c, report)
		if err != nil {
			return err
		}

		return s.Put(c, &xu)
	})
	return
}

func (s *ReportService) Update(c context.Context, xu XUser, params ReportUpdatingParams) (report *Report, err error) {
	v := validatekit.NewValidate()
	err = v.Struct(params)
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
	report.AuthorNickname = xu.Nickname

	if !params.ReportedAt.IsZero() {
		log.Infof(c, "update ReportedAt: %v", params.ReportedAt)
		report.ReportedAt = params.ReportedAt
	}
	report.UpdatedAt = s.now()

	err = v.Struct(report)
	if err != nil {
		return
	}
	err = s.Put(c, report)
	return
}

func (s *ReportService) now() time.Time {
	return s.timeProvider.Now()
}

func (s *ReportService) search(c context.Context, p ReportSerchParams, limit int) (reports []Report, err error) {
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

	reports = make([]Report, 0, limit)
	_, err = s.Goon(c).GetAll(q, &reports)
	return
}
