package xpo

import (
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/mjibson/goon"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"

	"local/apikit"
)

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
	// message := fmt.Sprintf("ALLOW_ORIGIN=%s", os.Getenv("ALLOW_ORIGIN"))

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/loggedin", handleLoggedIn)

	http.HandleFunc("/xreports", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getXReports(w, r)
			return
		}

		if r.Method == "POST" {
			if !responseIfUnauthorized(w, r) {
				return
			}

			postXReport(w, r)
			return
		}

		if r.Method == "OPTIONS" {
			allowClient(w)
		}
	})
}

func getXReports(w http.ResponseWriter, r *http.Request) {
	g := goon.NewGoon(r)

	q := datastore.NewQuery("Report").Order("-CreatedAt").Limit(10)
	reports := make([]Report, 0, 10)
	if _, err := g.GetAll(q, &reports); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseJSON(w, reports)
}

func postXReport(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	xu := xUserOrResponse(w, r)
	if xu == nil {
		log.Warningf(c, "xu==nil. response 401")
		return
	}

	jsonBody, err := apikit.ParseJSONBody(r)
	if err != nil {
		log.Warningf(c, "err: %v\n", err.Error())
		responseFailure(w, r, apikit.NewFailure(err.Error()), http.StatusBadRequest)
		return
	}

	log.Infof(c, "JSON: %v\n", jsonBody)

	g := goon.NewGoon(r)

	report := Report{
		Author:    xu.Name,
		Content:   jsonBody["content"].(string),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		AuthorKey: g.Key(xu),
	}
	_, err = g.Put(&report)
	if err != nil {
		log.Warningf(c, "err: %v\n", err.Error())
		responseFailure(w, r, apikit.NewFailure(err.Error()), http.StatusInternalServerError)
		return
	}

	responseJSON(w, report)
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
	logoutURL, _ := user.LogoutURL(c, "/")

	data := struct {
		LogoutURL string
	}{
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

func handleLoggedIn(w http.ResponseWriter, r *http.Request) {
	if !redirectUnlessLoggedIn(w, r) {
		return
	}

	c := appengine.NewContext(r)
	u := user.Current(c)
	g := goon.NewGoon(r)

	log.Infof(c, "logged in.")

	xu := &apikit.XUser{ID: u.ID}
	err := datastore.RunInTransaction(c, func(ctx context.Context) error {
		if err := g.Get(xu); err != nil {
			if err != datastore.ErrNoSuchEntity {
				return err
			}

			log.Infof(c, "XUser not found. create new one. : "+u.ID)
			xu = &apikit.XUser{ID: u.ID, Name: "user" + u.ID, Email: u.Email}
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

	http.Redirect(w, r, os.Getenv("ALLOW_ORIGIN"), http.StatusFound)
}
