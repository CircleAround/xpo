package web

import (
	"local/apikit/exchi"
	"local/xpo/app"
	"local/xpo/entities"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

func GetMe(c context.Context, w http.ResponseWriter, r *http.Request, xu *entities.XUser) error {
	return NewResponder(w, r).RenderMeOrError(xu, nil)
}

func PostMe(w http.ResponseWriter, r *http.Request) error {
	c := Context(r)
	u := user.Current(c)

	p := app.XUserProfileParams{}
	if err := parseJSONBody(r, &p); err != nil {
		return err
	}

	log.Infof(c, "params: %v\n", p)

	return NewResponder(w, r).RenderMeOrError(Services.XUser().Create(c, *u, p))
}

func GetByName(w http.ResponseWriter, r *http.Request) error {
	c := Context(r)
	p := exchi.URLParams(r)
	n := p.Get("authorName")
	return NewResponder(w, r).RenderMeOrError(Services.XUser().GetByName(c, n))
}

func UpdateMe(c context.Context, w http.ResponseWriter, r *http.Request, xu *entities.XUser) error {
	p := app.XUserProfileParams{}
	if err := parseJSONBody(r, &p); err != nil {
		return err
	}

	log.Infof(c, "params: %v\n", p)
	return NewResponder(w, r).RenderMeOrError(Services.XUser().Update(c, xu, p))
}
