package gaekit

import (
	"fmt"
	"html/template"
	"local/xpo/entities"
	"net/http"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

var tmpl = `
<!DOCTYPE html>
<html>
		<head>
			<meta charset="UTF-8"></meta>
      <title>Admin Page</title>
    </head>
		<body>
			<header>
				{{ if .User }}
					{{ .User }} as Admin (<a href="{{ .LogoutURL }}">Logout</a>)
				{{ else }}
					<a href="{{ .LoginURL }}">Login</a><br>
				{{ end }}
			</header>
			{{ if .User }}
				<h1>Admin Page</h1>
				<h2>Menu</h2>
				{{ range .Handler.Menues }}
					<form action="{{$.Handler.Path}}/{{.Path}}">
						<input type="submit" value="{{.Label}}">
					</form>
				{{ end }}
			{{ end }}
    </body>
</html>
`

type AdminPage struct {
	Handler   *AdminHandler
	User      *user.User
	LoginURL  string
	LogoutURL string
}

type AdminMenu struct {
	Label   string
	Path    string
	Handler http.Handler
}

type AdminHandler struct {
	Path   string
	Menues []*AdminMenu
}

func NewAdminHandler() *AdminHandler {
	return NewAdminHandlerWithPath("/admin")
}

func NewAdminHandlerWithPath(path string) *AdminHandler {
	return &AdminHandler{Path: path}
}

func (a *AdminHandler) AddMenuHandler(p string, h http.Handler, l string) {
	if l == "" {
		l = p
	}

	a.Menues = append(a.Menues, &AdminMenu{
		Label:   l,
		Path:    p,
		Handler: h,
	})
}

func (a *AdminHandler) AddMenu(p string, h http.HandlerFunc, l string) {
	a.AddMenuHandler(p, h, l)
}

func (a *AdminHandler) FindMenu(p string) *AdminMenu {
	for _, menu := range a.Menues {
		if menu.Path == p {
			return menu
		}
	}
	return nil
}

func (a *AdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := Context(r)
	path := r.URL.Path
	log.Infof(c, "Admin Access: %v", path)
	if path == a.Path || path == (a.Path+"/") {
		a.showTop(w, r)
		return
	}

	menuPath := strings.Replace(path, (a.Path + "/"), "", 1)
	menu := a.FindMenu(menuPath)
	if menu == nil {
		panic(fmt.Sprintf("menu not found; menu path: %v, path: ", menuPath, path))
	}

	menu.Handler.ServeHTTP(w, r)
}

func (a *AdminHandler) showTop(w http.ResponseWriter, r *http.Request) {
	c := Context(r)
	u := user.Current(c)
	if u != nil && !u.Admin {
		u = nil
	}

	loginURL, _ := user.LoginURL(c, a.Path)
	logoutURL, _ := user.LogoutURL(c, a.Path)

	funcMap := template.FuncMap{
		"authorName": func(r entities.Report) template.HTML {
			//nop sample
			return ""
		},
	}

	templates, err := template.New("").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templates.Execute(w, &AdminPage{
		Handler:   a,
		User:      u,
		LogoutURL: logoutURL,
		LoginURL:  loginURL,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Context(r *http.Request) context.Context {
	return appengine.NewContext(r)
}
