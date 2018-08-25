package xpo

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/mjibson/goon"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
)

// XUser struct
type XUser struct {
	ID    string `datastore:"-" goon:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Report struct
type Report struct {
	ID        int64          `datastore:"-" goon:"id"`
	AuthorKey *datastore.Key `datastore:"-" goon:"parent"`
	Author    string         `json:"author"`
	Content   string         `json:"content"`
	Year      int16          `json:"year"`
	Month     int8           `json:"month"`
	Day       int8           `json:"day"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func init() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/reports", handleReports)
	http.HandleFunc("/loggedin", handleLoggedIn)
	http.HandleFunc("/xreports", handleXReports)

	// http.HandleFunc("/my", handleMy)
}

func xUserKey(c context.Context, ID string) *datastore.Key {
	return datastore.NewKey(c, "XUser", ID, 0, nil)
}

func redirectUnlessLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	c := appengine.NewContext(r)
	u := user.Current(c)
	// ログインしてなければリダイレクト
	if u == nil {
		url, _ := user.LoginURL(c, "/loggedin")
		http.Redirect(w, r, url, http.StatusFound)
		return false
	}
	return true
}

func xUserOrRedirect(w http.ResponseWriter, r *http.Request) *XUser {
	c := appengine.NewContext(r)
	u := user.Current(c)
	g := goon.NewGoon(r)

	xu := &XUser{ID: u.ID}
	if err := g.Get(xu); err != nil {
		log.Print("Oops! has not user!")
		url, _ := user.LoginURL(c, "/loggedin")
		http.Redirect(w, r, url, http.StatusFound)
		return nil
	}
	return xu
}

func handleXReports(w http.ResponseWriter, r *http.Request) {
	g := goon.NewGoon(r)

	q := datastore.NewQuery("Report").Order("-CreatedAt").Limit(10)
	reports := make([]Report, 0, 10)
	if _, err := g.GetAll(q, &reports); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(reports)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Write(res)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	if !redirectUnlessLoggedIn(w, r) {
		return
	}

	xu := xUserOrRedirect(w, r)
	if xu == nil {
		return
	}

	c := appengine.NewContext(r)
	g := goon.NewGoon(r)

	q := datastore.NewQuery("Report").Ancestor(g.Key(xu)).Order("-CreatedAt").Limit(10)
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

	xu := xUserOrRedirect(w, r)
	if xu == nil {
		return
	}

	g := goon.NewGoon(r)

	report := Report{
		Author:    xu.Name,
		Content:   r.FormValue("content"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		AuthorKey: g.Key(xu),
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
	log.Println("logged in.")

	c := appengine.NewContext(r)
	u := user.Current(c)
	g := goon.NewGoon(r)

	xu := &XUser{ID: u.ID}
	err := datastore.RunInTransaction(c, func(ctx context.Context) error {
		if err := g.Get(xu); err != nil {
			if err != datastore.ErrNoSuchEntity {
				return err
			}

			log.Println("XUser not found. create new one. : " + u.ID)
			xu = &XUser{ID: u.ID, Name: "user" + u.ID, Email: u.Email}
			_, ierr := g.Put(xu)
			if ierr != nil {
				return ierr
			}
		}
		return nil
	}, nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
