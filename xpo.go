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

// XUser struct
type XUser struct {
	ID    string `datastore:"-" goon:"id"`
	Name  string
	Email string

	// 他のアカウントの関連付けも同様にする。各アカウントタイプでIDをユニークにするためにあえて抽象化しない
	AppEngineUserKey *datastore.Key
}

// AppEngineUser struct
type AppEngineUser struct {
	ID      string         `datastore:"-" goon:"id"`
	UserKey *datastore.Key // parent を付けないのはログイン時、単独で本structを検索させるため
	Email   string
}

// Report struct
type Report struct {
	ID        string         `datastore:"-" goon:"id"`
	AuthorKey *datastore.Key `datastore:"-" goon:"parent"`
	Author    string
	Content   string
	Year      int16
	Month     int8
	Day       int8
	CreatedAt time.Time
	UpdatedAt time.Time
}

func init() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/reports", handleReports)
	http.HandleFunc("/loggedin", handleLoggedIn)

	// http.HandleFunc("/my", handleMy)
}

func userKey(c context.Context, email string) *datastore.Key {
	return datastore.NewKey(c, "User", email, 0, nil)
}

func xUserKey(c context.Context, ID string) *datastore.Key {
	return datastore.NewKey(c, "XUser", ID, 0, nil)
}

func appengineUserKey(c context.Context, ID string) *datastore.Key {
	return datastore.NewKey(c, "AppEngineUser", ID, 0, nil)
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
	if !redirectUnlessLoggedIn(w, r) {
		return
	}

	c := appengine.NewContext(r)
	u := user.Current(c)
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

func handleLoggedIn(w http.ResponseWriter, r *http.Request) {
	if !redirectUnlessLoggedIn(w, r) {
		return
	}

	c := appengine.NewContext(r)
	u := user.Current(c)
	g := goon.NewGoon(r)

	aeu := &AppEngineUser{ID: u.ID}
	if err := g.Get(aeu); err != nil {
		if err != datastore.ErrNoSuchEntity {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		xu := &XUser{Name: u.ID, Email: u.Email}
		xkey, ierr := g.Put(xu)
		if ierr != nil {
			http.Error(w, ierr.Error(), http.StatusInternalServerError)
			return
		}

		aeu.Email = u.Email
		aeu.UserKey = xkey
		akey, ierr := g.Put(aeu)
		if ierr != nil {
			http.Error(w, ierr.Error(), http.StatusInternalServerError)
			return
		}

		xu.AppEngineUserKey = akey
		xkey, ierr = g.Put(xu)
		if ierr != nil {
			http.Error(w, ierr.Error(), http.StatusInternalServerError)
			return
		}
	}

	q := datastore.NewQuery("AppEngineUser").Order("-CreatedAt").Limit(10)
	reports := make([]Report, 0, 10)
	if _, err := g.GetAll(q, &reports); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
