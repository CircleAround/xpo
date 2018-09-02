package xpo

import (
	"html/template"
	"net/http"
	"os"

	"github.com/mjibson/goon"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"

	"local/apikit"
)

func init() {
	// message := fmt.Sprintf("ALLOW_ORIGIN=%s", os.Getenv("ALLOW_ORIGIN"))

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/loggedin", handleLoggedIn)

	http.HandleFunc("/reports", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getReports(w, r)
			return
		}

		if r.Method == "POST" {
			if !responseIfUnauthorized(w, r) {
				return
			}

			postReport(w, r)
			return
		}

		if r.Method == "OPTIONS" {
			allowClient(w)
		}
	})
}

func getReports(w http.ResponseWriter, r *http.Request) {
	s := NewReportService(r)
	reports, err := s.RetriveAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responseJSON(w, reports)
}

func postReport(w http.ResponseWriter, r *http.Request) {
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

	content := jsonBody["content"].(string)
	s := NewReportService(r)
	report := Report{Content: content}

	err = s.Create(xu, &report)
	if err != nil {
		switch err.(type) {
		default:
			log.Warningf(c, "err: %v\n", err.Error())
			responseFailure(w, r, apikit.NewFailure(err.Error()), http.StatusInternalServerError)
			return
		case *apikit.ValidationError:
			responseFailure(w, r, apikit.NewFailure(err), http.StatusUnprocessableEntity)
			return
		}
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
