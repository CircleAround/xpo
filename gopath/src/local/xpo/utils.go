package xpo

import (
	"net/http"
	"os"

	"github.com/mjibson/goon"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"

	"local/apikit"
)

func allowOrigin(w http.ResponseWriter, origin string) {
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
}

func allowClient(w http.ResponseWriter) {
	allowOrigin(w, os.Getenv("ALLOW_ORIGIN"))
}

func safeFilter(w http.ResponseWriter, r *http.Request, err error) {
	c := appengine.NewContext(r)

	if err != nil {
		switch err.(type) {
		default:
			log.Warningf(c, "err: %v\n", err.Error())
			apikit.ResponseFailure(w, r, err, http.StatusInternalServerError)
			return

		case *ValueNotUniqueError:
			apikit.ResponseFailure(w, r, err, http.StatusUnprocessableEntity)
			return

		case *apikit.ValidationError:
			apikit.ResponseFailure(w, r, err, http.StatusUnprocessableEntity)
			return
		}
	}
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
	c := appengine.NewContext(r)
	url, _ := user.LoginURL(c, "/loggedin")
	code := http.StatusUnauthorized

	apikit.ResponseFailure(w, r, url, code)
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
