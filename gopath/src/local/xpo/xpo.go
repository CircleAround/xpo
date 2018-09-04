package xpo

import (
	"html/template"
	"net/http"
	"os"

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

	http.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
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

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	http.HandleFunc("/reports", func(w http.ResponseWriter, r *http.Request) {
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
}

func getMe(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	u := user.Current(c)

	s := NewXUserService(c)
	xu, err := s.GetByUser(u)
	if err == nil {
		apikit.ResponseJSON(w, xu)
		return nil
	}

	if err == datastore.ErrNoSuchEntity {
		apikit.ResponseJSON(w, "BE_SIGN_UP")
		return nil
	}

	return err
}

func postMe(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	u := user.Current(c)

	jsonBody, err := apikit.ParseJSONBody(r)
	if err != nil {
		log.Warningf(c, "err: %v\n", err.Error())
		apikit.ResponseFailure(w, r, err, http.StatusBadRequest)
		return nil
	}

	log.Infof(c, "JSON: %v\n", jsonBody)

	name := jsonBody["name"].(string)
	nickname := jsonBody["nickname"].(string)

	s := NewXUserService(c)
	xu, err := s.Create(u, name, nickname)

	if err != nil {
		switch err.(type) {
		case *AlreadyKeptNameError:
			apikit.ResponseFailure(w, r, err, http.StatusUnprocessableEntity)
			return nil
		}
		return err
	}

	apikit.ResponseJSON(w, xu)
	return nil
}

func getReports(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	s := NewReportService(c)
	reports, err := s.RetriveAll()
	if err != nil {
		return err
	}
	apikit.ResponseJSON(w, reports)
	return nil
}

func postReport(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)

	xu := xUserOrResponse(w, r)
	if xu == nil {
		log.Warningf(c, "xu==nil. response 401")
		return nil
	}

	jsonBody, err := apikit.ParseJSONBody(r)
	if err != nil {
		log.Warningf(c, "err: %v\n", err.Error())
		apikit.ResponseFailure(w, r, err, http.StatusBadRequest)
		return nil
	}

	log.Infof(c, "JSON: %v\n", jsonBody)

	content := jsonBody["content"].(string)
	s := NewReportService(c)
	report := Report{Content: content}

	err = s.Create(xu, &report)
	if err != nil {
		return err
	}

	apikit.ResponseJSON(w, report)
	return nil
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
	allowClient(w)

	if !redirectUnlessLoggedIn(w, r) {
		return
	}

	http.Redirect(w, r, os.Getenv("ALLOW_ORIGIN"), http.StatusFound)
}
