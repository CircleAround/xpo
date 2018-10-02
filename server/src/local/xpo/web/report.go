package web

import (
	"local/apikit/exchi"
	"local/xpo/app"
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/appengine/log"
)

func GetReport(w http.ResponseWriter, r *http.Request) error {
	c := Context(r)
	p := exchi.URLParams(r)

	uid := p.Get("authorId")
	id, err := p.AsInt64("id")
	if err != nil {
		return err
	}

	return NewResponder(w, r).RenderObjectOrError(Services.Report().Find(c, uid, id))
}

func GetReports(w http.ResponseWriter, r *http.Request) error {
	c := Context(r)
	return NewResponder(w, r).RenderObjectOrError(Services.Report().RetriveAll(c))
}

func SearchReportsYmd(w http.ResponseWriter, r *http.Request) error {
	c := Context(r)

	p := exchi.URLParams(r)
	uid := p.Get("authorId")
	y, err := p.AsInt("year")
	if err != nil {
		return err
	}
	m, err := p.AsInt("month")
	if err != nil {
		return err
	}
	d, err := p.AsInt("day")
	if err != nil {
		return err
	}

	return NewResponder(w, r).RenderObjectOrError(Services.Report().SearchBy(c, uid, y, m, d))
}

func PostReport(c context.Context, w http.ResponseWriter, r *http.Request, xu *app.XUser) error {
	p := &app.ReportCreationParams{}
	if err := parseJSONBody(r, &p); err != nil {
		return err
	}

	log.Infof(c, "params: %v\n", p)
	return NewResponder(w, r).RenderObjectOrError(Services.Report().Create(c, *xu, *p))
}

func UpdateReport(c context.Context, w http.ResponseWriter, r *http.Request, xu *app.XUser) error {
	p := &app.ReportUpdatingParams{}
	{
		if err := parseJSONBody(r, &p); err != nil {
			return err
		}

		id, err := exchi.URLParams(r).AsInt64("id")
		if err != nil {
			return err
		}
		p.ID = id
	}

	log.Infof(c, "params: %v\n", p)

	return NewResponder(w, r).RenderObjectOrError(Services.Report().Update(c, *xu, *p))
}
