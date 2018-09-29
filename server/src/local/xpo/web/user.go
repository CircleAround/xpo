package web

import (
	"local/apikit"
	"local/xpo/app"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

func GetMe(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	u := user.Current(c)

	xu, err := Services.XUser().GetByUser(c, *u)
	if err == nil {
		return NewResponder(w, r).RenderMeOrError(xu, err)
	}

	if err == datastore.ErrNoSuchEntity {
		return apikit.RespondJSON(w, "BE_SIGN_UP")
	}

	return err
}

func PostMe(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	u := user.Current(c)

	p := app.XUserProfileParams{}
	if err := parseJSONBody(r, &p); err != nil {
		return err
	}

	log.Infof(c, "params: %v\n", p)

	return NewResponder(w, r).RenderMeOrError(Services.XUser().Create(c, *u, p))
}

func UpdateMe(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	u := user.Current(c)

	p := app.XUserProfileParams{}
	if err := parseJSONBody(r, &p); err != nil {
		return err
	}

	log.Infof(c, "params: %v\n", p)
	return NewResponder(w, r).RenderMeOrError(Services.XUser().Update(c, *u, p))
}
