package web

import (
	"html/template"
	"local/apikit"
	"local/xpo/app"
	"net/http"
	"os"

	"google.golang.org/appengine"
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

func parseJSONBody(r *http.Request, p interface{}) error {
	err := apikit.ParseJSONBody(r, &p)
	if err != nil {
		return apikit.NewInvalidParameterErrorWithMessage("json", err.Error())
	}
	return nil
}
