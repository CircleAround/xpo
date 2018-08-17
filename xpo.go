package xpo

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mjibson/goon"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
)

// Report struct
type Report struct {
	ID        string         `datastore:"-" goon:"id"`
	AuthorKey *datastore.Key `datastore:"-" goon:"parent"`
	Author    string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func init() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/reports", handleReports)

	// http.HandleFunc("/my", handleMy)
	// http.HandleFunc("/loggedin", handleLoggedIn)
}

func userKey(c context.Context, email string) *datastore.Key {
	return datastore.NewKey(c, "User", email, 0, nil)
}

func redirectUnlessLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	c := appengine.NewContext(r)
	u := user.Current(c)
	// ログインしてなければリダイレクト
	if u == nil {
		url, _ := user.LoginURL(c, "/")
		http.Redirect(w, r, url, http.StatusFound)
		return false
	}
	return true
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	// ログインしてなければリダイレクト
	if !redirectUnlessLoggedIn(w, r) {
		return
	}

	c := appengine.NewContext(r)
	u := user.Current(c)
	g := goon.NewGoon(r)

	q := datastore.NewQuery("Report").Ancestor(userKey(c, u.String())).Order("-CreatedAt").Limit(10)
	reports := make([]Report, 0, 10)
	if _, err := g.GetAll(q, &reports); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logoutURL, _ := user.LogoutURL(c, "/")

	data := struct {
		Reports   []Report
		LogoutURL string
	}{
		Reports:   reports,
		LogoutURL: logoutURL,
	}

	// view構築する

	funcMap := template.FuncMap{
		"authorName": func(r Report) template.HTML {
			//nop sample
			return ""
		},
	}
	templates := template.Must(template.New("").Funcs(funcMap).ParseFiles("views/layout.html", "views/my.html"))

	if err := templates.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleReports(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)

	// ログインしてなければリダイレクト
	if u == nil {
		url, _ := user.LoginURL(c, "/")
		http.Redirect(w, r, url, http.StatusFound)
		return
	}

	g := goon.NewGoon(r)

	id, ierr := uuid.NewRandom()
	if ierr != nil {
		http.Error(w, ierr.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Hello log")
	authorKey := userKey(c, u.String())
	log.Println(authorKey)

	report := Report{
		ID:        id.String(),
		Author:    u.ID,
		Content:   r.FormValue("content"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		AuthorKey: authorKey,
	}
	_, err := g.Put(&report)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
