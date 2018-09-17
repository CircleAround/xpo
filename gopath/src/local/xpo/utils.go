package xpo

import (
	"fmt"
	"net/http"
	"os"
	"reflect"

	"github.com/mjibson/goon"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"

	"github.com/iancoleman/strcase"
	validator "gopkg.in/go-playground/validator.v9"

	"local/apikit"
	"local/gaekit"
)

func FullURL(r *http.Request, path string) string {
	hostName := r.Host
	var scheme string
	if r.TLS == nil {
		scheme = "http"
	} else {
		scheme = "https"
	}
	return fmt.Sprintf("%v://%v%v", scheme, hostName, path)
}

func LoginFullURL(r *http.Request) string {
	c := appengine.NewContext(r)
	url, _ := user.LoginURL(c, "/loggedin")
	return safeFullUrl(r, url)
}

func LogoutFullURL(r *http.Request) string {
	c := appengine.NewContext(r)
	url, _ := user.LogoutURL(c, "/")
	return safeFullUrl(r, url)
}

func safeFullUrl(r *http.Request, url string) string {
	if appengine.IsDevAppServer() {
		return FullURL(r, url)
	}
	return url
}

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
		case *gaekit.ValueNotUniqueError:
		case *gaekit.DuplicatedObjectError:
		case *validator.InvalidValidationError:
			apikit.ResponseFailure(w, r, err, http.StatusUnprocessableEntity)
			return

		case validator.ValidationErrors:
			ve := apikit.NewValidationError()
			for _, err := range err.(validator.ValidationErrors) {
				ve.PushOne(strcase.ToSnake(err.Field()), err.Tag())
			}
			apikit.ResponseFailure(w, r, ve, http.StatusUnprocessableEntity)
			return

		default:
			log.Warningf(c, "err: %v, %v\n", err.Error(), reflect.TypeOf(err))
			apikit.ResponseFailure(w, r, err, http.StatusInternalServerError)
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
	code := http.StatusUnauthorized
	apikit.ResponseFailure(w, r, LoginFullURL(r), code)
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
