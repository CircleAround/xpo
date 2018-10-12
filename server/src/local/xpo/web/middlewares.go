package web

import (
	"net/http"
	"os"
	"strings"

	"google.golang.org/appengine/user"
)

func CrossOriginable(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowClient(w)
		if r.Method == "OPTIONS" {
			return
		}

		if !checkHeaders(r) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := Context(r)
		u := user.Current(c)
		if u == nil {
			responseUnauthorized(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func checkHeaders(r *http.Request) bool {
	h := r.Header.Get("X-Requested-With")
	if strings.ToLower(h) != "xmlhttprequest" {
		return false
	}

	h = r.Header.Get("Origin")
	if h != os.Getenv("ALLOW_ORIGIN") {
		return false
	}	
	return true
}

func allowOrigin(w http.ResponseWriter, origin string) {
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("X-Content-Type-Options", "nosniff")
}

func allowClient(w http.ResponseWriter) {
	allowOrigin(w, os.Getenv("ALLOW_ORIGIN"))
}
