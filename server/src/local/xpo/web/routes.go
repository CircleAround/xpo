package web

import (
	"net/http"

	"github.com/go-chi/chi"
)

func Routes() {
	http.Handle("/", Router())
}

func Router() *chi.Mux {
	r := chi.NewRouter()
	r.Use(CrossOriginable)

	r.Get("/", handleRoot)
	r.Get(loggedInPath, handleLoggedIn)

	r.Route("/users/me", func(r chi.Router) {
		r.Use(GAuth)

		r.Post("/", Handler(PostMe))

		r.Get("/", Handler(Auth(GetMe)))
		r.Put("/", Handler(Auth(UpdateMe)))
	})

	r.Route("/reports", func(r chi.Router) {
		r.With(GAuth).Get(
			"/{authorId:[0-9]+}/_/{year:[0-9]+}/{month:[0-9]+}/{day:[0-9]+}",
			Handler(SearchReportsYmd),
		)

		r.Route("/{authorId:[0-9]+}/{id:[0-9]+}", func(r chi.Router) {
			r.Get("/", Handler(GetReport))

			r.With(GAuth).Put("/", Handler(Auth(UpdateReport)))
		})

		r.Route("/", func(r chi.Router) {
			r.Get("/", Handler(GetReports))

			r.With(GAuth).Post("/", Handler(Auth(PostReport)))
		})
	})
	return r
}
