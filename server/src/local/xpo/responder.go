package xpo

import (
	"local/apikit"
	"local/gaekit"
	"net/http"

	"github.com/mjibson/goon"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

type XUserResponse struct {
	XUser
	LoginURL  string `json:"loginUrl"`
	LogoutURL string `json:"logoutUrl"`
}

type Responder struct {
	w  http.ResponseWriter
	r  *http.Request
	jr apikit.JSONRenderer
}

func NewResponder(w http.ResponseWriter, r *http.Request) *Responder {
	return &Responder{w: w, r: r, jr: apikit.NewJSONRenderer(w)}
}

func (r *Responder) RenderMeOrError(xu *XUser, err error) error {
	res := XUserResponse{
		XUser:     *xu,
		LoginURL:  gaekit.LoginFullURL(r.r, "/loggedin"),
		LogoutURL: gaekit.LogoutFullURL(r.r),
	}
	return r.jr.RenderOrError(res, err)
}

func (r *Responder) RenderObjectOrError(obj interface{}, err error) error {
	return r.jr.RenderOrError(obj, err)
}

func redirectUnlessLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	c := appengine.NewContext(r)
	u := user.Current(c)
	// ログインしてなければリダイレクト
	if u != nil {
		return true
	}

	url, _ := user.LoginURL(c, "/loggedin")
	http.Redirect(w, r, url, http.StatusFound)
	return false
}

func xUserOrRedirect(w http.ResponseWriter, r *http.Request) *XUser {
	c := appengine.NewContext(r)
	u := user.Current(c)
	g := goon.NewGoon(r)

	xu := &XUser{ID: u.ID}
	if err := g.Get(xu); err != nil {
		log.Warningf(c, "Oops! has not user!")
		url, _ := user.LoginURL(c, "/loggedin")
		http.Redirect(w, r, url, http.StatusFound)
		return nil
	}
	return xu
}

func xUserOrResponse(w http.ResponseWriter, r *http.Request) *XUser {
	c := appengine.NewContext(r)
	u := user.Current(c)
	g := goon.NewGoon(r)

	xu := &XUser{ID: u.ID}
	if err := g.Get(xu); err != nil {
		log.Warningf(c, "Oops! has not user!")
		responseUnauthorized(w, r)
		return nil
	}
	return xu
}

func responseUnauthorized(w http.ResponseWriter, r *http.Request) {
	code := http.StatusUnauthorized
	apikit.RespondFailure(w, gaekit.LoginFullURL(r, "/loggedin"), code)
}

func responseIfUnauthorized(w http.ResponseWriter, r *http.Request) bool {
	c := appengine.NewContext(r)
	u := user.Current(c)

	log.Infof(c, "user: %v", u)

	if u != nil {
		return true
	}

	// ログインしてなければリダイレクト
	responseUnauthorized(w, r)
	return false
}
