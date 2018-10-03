package web

import (
	"html/template"
	"local/xpo/app"
	"net/http"
	"os"

	"github.com/mjibson/goon"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	if !redirectUnlessLoggedIn(w, r) {
		return
	}

	xu := xUserOrRedirect(w, r)
	if xu == nil {
		return
	}

	c := Context(r)
	logoutURL, _ := user.LogoutURL(c, "/")

	data := struct {
		LogoutURL string
		Version   string
	}{
		LogoutURL: logoutURL,
		Version:   os.Getenv("CURRENT_VERSION_ID"),
	}

	// view構築する

	funcMap := template.FuncMap{
		"authorName": func(r app.Report) template.HTML {
			//nop sample
			return ""
		},
	}
	templates := template.Must(template.New("").Funcs(funcMap).ParseFiles("views/layout.html", "views/my.html"))

	if err := templates.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleLoggedIn(w http.ResponseWriter, r *http.Request) {
	allowClient(w)

	if !redirectUnlessLoggedIn(w, r) {
		return
	}

	http.Redirect(w, r, os.Getenv("ALLOW_ORIGIN"), http.StatusFound)
}

func redirectUnlessLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	c := Context(r)
	u := user.Current(c)
	// ログインしてなければリダイレクト
	if u != nil {
		return true
	}

	url, _ := user.LoginURL(c, loggedInPath)
	http.Redirect(w, r, url, http.StatusFound)
	return false
}

func xUserOrRedirect(w http.ResponseWriter, r *http.Request) *app.XUser {
	c := Context(r)
	u := user.Current(c)
	g := goon.NewGoon(r)

	xu := &app.XUser{ID: u.ID}
	if err := g.Get(xu); err != nil {
		log.Warningf(c, "Oops! has not user!")
		url, _ := user.LoginURL(c, loggedInPath)
		http.Redirect(w, r, url, http.StatusFound)
		return nil
	}
	return xu
}
