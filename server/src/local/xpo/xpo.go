package xpo

import (
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"

	"local/apikit"
)

//go:generate go-assets-builder --output=assets/reserved_username_list.go -p=assets ../../../assets/reserved_username_list

type XUserResponse struct {
	XUser
	LoginURL  string `json:"loginUrl"`
	LogoutURL string `json:"logoutUrl"`
}

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/", handleRoot)
	r.HandleFunc("/loggedin", handleLoggedIn)

	r.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		allowClient(w)
		if r.Method == "OPTIONS" {
			return
		}

		if !responseIfUnauthorized(w, r) {
			return
		}

		if r.Method == "GET" {
			safeFilter(w, r, getMe(w, r))
			return
		}
		if r.Method == "POST" {
			safeFilter(w, r, postMe(w, r))
			return
		}
		if r.Method == "PUT" {
			safeFilter(w, r, updateMe(w, r))
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	r.HandleFunc("/reports/{authorId:[0-9]+}/_/{year:[0-9]+}/{month:[0-9]+}/{day:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		allowClient(w)
		if r.Method == "OPTIONS" {
			return
		}
		if r.Method == "GET" {
			if !responseIfUnauthorized(w, r) {
				return
			}

			safeFilter(w, r, searchReportsYmd(w, r))
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	r.HandleFunc("/reports/{authorId:[0-9]+}/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		allowClient(w)
		if r.Method == "OPTIONS" {
			return
		}

		if r.Method == "GET" {
			if !responseIfUnauthorized(w, r) {
				return
			}

			safeFilter(w, r, getReport(w, r))
			return
		}

		if r.Method == "PUT" {
			if !responseIfUnauthorized(w, r) {
				return
			}

			safeFilter(w, r, updateReport(w, r))
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	r.HandleFunc("/reports", func(w http.ResponseWriter, r *http.Request) {
		allowClient(w)
		if r.Method == "OPTIONS" {
			return
		}

		if r.Method == "GET" {
			safeFilter(w, r, getReports(w, r))
			return
		}

		if r.Method == "POST" {
			if !responseIfUnauthorized(w, r) {
				return
			}

			safeFilter(w, r, postReport(w, r))
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	http.Handle("/", r)
}

func getMe(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	u := user.Current(c)

	xu, err := NewXUserService().GetByUser(c, *u)
	if err == nil {
		return NewResponder(w, r).RenderMeOrError(xu, err)
	}

	if err == datastore.ErrNoSuchEntity {
		return apikit.RespondJSON(w, "BE_SIGN_UP")
	}

	return err
}

func postMe(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	u := user.Current(c)

	p := XUserProfileParams{}
	if err := parseJSONBody(r, &p); err != nil {
		return err
	}

	log.Infof(c, "params: %v\n", p)

	return NewResponder(w, r).RenderMeOrError(NewXUserService().Create(c, *u, p))
}

func updateMe(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	u := user.Current(c)

	p := XUserProfileParams{}
	if err := parseJSONBody(r, &p); err != nil {
		return err
	}

	log.Infof(c, "params: %v\n", p)
	return NewResponder(w, r).RenderMeOrError(NewXUserService().Update(c, *u, p))
}

func getReport(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	p := apikit.URLParams(r)

	uid := p.Get("authorId")
	id, err := p.AsInt64("id")
	if err != nil {
		return err
	}

	return NewResponder(w, r).RenderObjectOrError(NewReportService().Find(c, uid, id))
}

func getReports(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	return NewResponder(w, r).RenderObjectOrError(NewReportService().RetriveAll(c))
}

func searchReportsYmd(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)

	p := apikit.URLParams(r)
	uid := p.Get("authorId")
	y, err := p.AsInt("year")
	if err != nil {
		return err
	}
	m, err := p.AsInt("month")
	if err != nil {
		return err
	}
	d, err := p.AsInt("day")
	if err != nil {
		return err
	}

	return NewResponder(w, r).RenderObjectOrError(NewReportService().SearchBy(c,uid, y, m, d))
}

func postReport(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)

	xu := xUserOrResponse(w, r)
	if xu == nil {
		log.Warningf(c, "xu==nil. response 401")
		return nil
	}

	p := &ReportCreationParams{}
	if err := parseJSONBody(r, &p); err != nil {
		return err
	}

	log.Infof(c, "params: %v\n", p)
	return NewResponder(w, r).RenderObjectOrError(NewReportService().Create(c, *xu, *p))
}

func updateReport(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)

	xu := xUserOrResponse(w, r)
	if xu == nil {
		log.Warningf(c, "xu==nil. response 401")
		return nil
	}

	p := &ReportUpdatingParams{}
	{
		if err := parseJSONBody(r, &p); err != nil {
			return err
		}

		id, err := apikit.URLParams(r).AsInt64("id")
		if err != nil {
			return err
		}
		p.ID = id
	}

	log.Infof(c, "params: %v\n", p)

	return NewResponder(w, r).RenderObjectOrError(NewReportService().Update(c, *xu, *p))
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
		Version   string
	}{
		LogoutURL: logoutURL,
		Version:   os.Getenv("CURRENT_VERSION_ID"),
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
	allowClient(w)

	if !redirectUnlessLoggedIn(w, r) {
		return
	}

	http.Redirect(w, r, os.Getenv("ALLOW_ORIGIN"), http.StatusFound)
}

func parseJSONBody(r *http.Request, p interface{}) error {
	err := apikit.ParseJSONBody(r, &p)
	if err != nil {
		return apikit.NewInvalidParameterErrorWithMessage("json", err.Error())
	}
	return nil
}
