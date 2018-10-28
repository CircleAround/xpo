package app

import (
	"local/the_time"
	"local/validatekit"
	"local/xpo/entities"
	"local/xpo/store"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

type ReportService struct {
	timeProvider the_time.Provider
	rrep         *store.ReportRepository
}

type ReportCreationParams struct {
	Content     string    `json:"content" validate:"required"`
	ContentType string    `json:"contentType" validate:"required"`
	ReportedAt  time.Time `json:"reportedAt"`
	Languages   []string  `json:"languages" validate:"oneof=1c abnf accesslog actionscript ada apache applescript cpp arduino armasm xml asciidoc aspectj autohotkey autoit avrasm awk axapta bash basic bnf brainfuck cal capnproto ceylon clean clojure clojure-repl cmake coffeescript coq cos crmsh crystal cs csp css d markdown dart delphi diff django dns dockerfile dos dsconfig dts dust ebnf elixir elm ruby erb erlang-repl erlang excel fix flix fortran fsharp gams gauss gcode gherkin glsl go golo gradle groovy haml handlebars haskell haxe hsp htmlbars http hy inform7 ini irpf90 java javascript jboss-cli json julia julia-repl kotlin lasso ldif leaf less lisp livecodeserver livescript llvm lsl lua makefile mathematica matlab maxima mel mercury mipsasm mizar perl mojolicious monkey moonscript n1ql nginx nimrod nix nsis objectivec ocaml openscad oxygene parser3 pf php pony powershell processing profile prolog protobuf puppet purebasic python q qml r rib roboconf routeros rsl ruleslanguage rust scala scheme scilab scss shell smali smalltalk sml sqf sql stan stata step21 stylus subunit swift taggerscript yaml tap tcl tex thrift tp twig typescript vala vbnet vbscript vbscript-html verilog vhdl vim x86asm xl xquery zephir"`
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
	s.rrep = store.NewReportRepository()
	return s
}

func (s *ReportService) RetriveAll(c context.Context) (reports []entities.Report, err error) {
	limit := 30
	return s.rrep.Search(c, store.ReportSearchParams{}, limit)
}

func (s *ReportService) SearchBy(c context.Context, authorID string, year int, month int, day int) (reports []entities.Report, err error) {
	limit := 30

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return
	}
	from := time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc)
	to := from.AddDate(0, 0, 1)

	return s.rrep.Search(c, store.ReportSearchParams{
		AuthorID:       authorID,
		ReportedAtFrom: from,
		ReportedAtTo:   to,
	}, limit)
}

func (s *ReportService) SearchByAuthor(c context.Context, authorID string) (reports []entities.Report, err error) {
	limit := 30
	return s.rrep.Search(c, store.ReportSearchParams{
		AuthorID: authorID,
	}, limit)
}

func (s *ReportService) Find(c context.Context, uid string, id int64) (report *entities.Report, err error) {
	xu := entities.XUser{ID: uid}
	return s.FindByXUserAndID(c, xu, id)
}

func (s *ReportService) FindByXUserAndID(c context.Context, xu entities.XUser, id int64) (report *entities.Report, err error) {
	ak := s.rrep.KeyOf(c, xu)
	report = &entities.Report{AuthorKey: ak, ID: id}
	err = s.rrep.Get(c, report)
	return
}

func (s *ReportService) Create(c context.Context, xu entities.XUser, params ReportCreationParams) (report *entities.Report, err error) {
	v := validatekit.NewValidate()
	err = v.Struct(params)
	if err != nil {
		return
	}

	report = &entities.Report{}
	report.Content = params.Content
	report.ContentType = params.ContentType
	report.Author = xu.Name
	report.AuthorID = xu.ID
	report.AuthorKey = s.rrep.KeyOf(c, xu)
	report.AuthorNickname = xu.Nickname
	report.Languages = params.Languages

	now := s.now()
	var ra time.Time
	if params.ReportedAt.IsZero() {
		ra = now
	} else {
		ra = params.ReportedAt
	}
	report.ReportedAt = ra
	report.CreatedAt = now
	report.UpdatedAt = now

	err = v.Struct(report)
	if err != nil {
		return
	}

	err = s.rrep.Create(c, &xu, report)
	return
}

func (s *ReportService) Update(c context.Context, xu entities.XUser, params ReportUpdatingParams) (report *entities.Report, err error) {
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
	err = s.rrep.Put(c, report)
	return
}

func (s *ReportService) now() time.Time {
	return s.timeProvider.Now()
}
