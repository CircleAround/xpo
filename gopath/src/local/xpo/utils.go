package xpo

import (
	"encoding/json"
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

func redirectUnlessLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	c := appengine.NewContext(r)
	u := user.Current(c)
	// ログインしてなければリダイレクト
	if u != nil {
		return true
	}

	url, _ := user.LoginURL(c, "/loggedin")
	allowClient(w)
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

func responseJSON(w http.ResponseWriter, obj interface{}) {
	res, err := json.Marshal(obj)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	allowClient(w)
	w.Write(res)
}

func responseOk(w http.ResponseWriter) {
	responseJSON(w, apikit.NewSuccess())
}

func responseFailure(w http.ResponseWriter, r *http.Request, failure apikit.Failure, code int) {
	allowClient(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	res, err := json.Marshal(failure)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(res)
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
	failure := apikit.NewFailure(url)
	code := http.StatusUnauthorized

	responseFailure(w, r, failure, code)
}

func responseIfUnauthorized(w http.ResponseWriter, r *http.Request) bool {
	c := appengine.NewContext(r)
	u := user.Current(c)
	// ログインしてなければリダイレクト
	if u != nil {
		return true
	}

	responseUnauthorized(w, r)
	return false
}
