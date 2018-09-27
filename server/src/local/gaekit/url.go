package gaekit

import (
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
)

func FullURL(r *http.Request, path string) string {
	hostName := r.Host
	var scheme string
	if r.TLS == nil {
		scheme = "http"
	} else {
		scheme = "https"
	}
	return fmt.Sprintf("%v://%v%v", scheme, hostName, path)
}

func LoginFullURL(r *http.Request, dest string) string {
	c := appengine.NewContext(r)
	url, _ := user.LoginURL(c, dest)
	return safeFullUrl(r, url)
}

func LogoutFullURL(r *http.Request) string {
	c := appengine.NewContext(r)
	url, _ := user.LogoutURL(c, "/")
	return safeFullUrl(r, url)
}

func safeFullUrl(r *http.Request, url string) string {
	if appengine.IsDevAppServer() {
		return FullURL(r, url)
	}
	return url
}
