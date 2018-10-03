package web

import (
	"local/xpo/app"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

func GetMe(c context.Context, w http.ResponseWriter, r *http.Request, xu *app.XUser) error {
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

func UpdateMe(c context.Context, w http.ResponseWriter, r *http.Request, xu *app.XUser) error {
	p := app.XUserProfileParams{}
	if err := parseJSONBody(r, &p); err != nil {
		return err
	}

	log.Infof(c, "params: %v\n", p)
	return NewResponder(w, r).RenderMeOrError(Services.XUser().Update(c, xu, p))
}
