package xpo

import (
	"html/template"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"

	"local/apikit"
	"local/apikit/exchi"
)

//go:generate go-assets-builder --output=assets/reserved_username_list.go -p=assets ../../../assets/reserved_username_list

func XOriginable(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowClient(w)
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AuthorizedCheckable(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !responseIfUnauthorized(w, r) {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Catch(handler func(http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		safeFilter(w, r, handler(w, r))
	}
}

func init() {
	r := chi.NewRouter()
	r.Use(XOriginable)

	r.Get("/", handleRoot)
	r.Get("/loggedin", handleLoggedIn)

	r.Route("/users/me", func(r chi.Router) {
		r.Use(AuthorizedCheckable)

		r.Get("/", Catch(getMe))
		r.Post("/", Catch(postMe))
		r.Put("/", Catch(updateMe))
	})

	r.Route("/reports", func(r chi.Router) {
		r.With(AuthorizedCheckable).Get(
			"/{authorId:[0-9]+}/_/{year:[0-9]+}/{month:[0-9]+}/{day:[0-9]+}",
			Catch(searchReportsYmd),
		)

		r.Route("/{authorId:[0-9]+}/{id:[0-9]+}", func(r chi.Router) {
			r.Use(AuthorizedCheckable)

			r.Get("/", Catch(getReport))
			r.Put("/", Catch(updateReport))
		})

		r.Route("/", func(r chi.Router) {
			r.Get("/", Catch(getReports))
			r.With(AuthorizedCheckable).Post("/", Catch(postReport))
		})
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
	p := exchi.URLParams(r)

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

	p := exchi.URLParams(r)
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

	return NewResponder(w, r).RenderObjectOrError(NewReportService().SearchBy(c, uid, y, m, d))
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

		id, err := exchi.URLParams(r).AsInt64("id")
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
