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

	r.Get("/", handleRoot)
	r.Get(loggedInPath, handleLoggedIn)

	r.Route("/users", func(r chi.Router) {
		r.Use(CrossOriginable)

		r.Route("/me", func(r chi.Router) {
			r.Use(GAuth)

			r.Post("/", Handler(PostMe))

			r.Get("/", Handler(Auth(GetMe)))
			r.Put("/", Handler(Auth(UpdateMe)))
		})

		r.Get("/{authorName:[a-z][0-9a-z_]+}", Handler(GetUserByName))
		r.Get("/{authorId:[1-9][0-9]*}", Handler(GetUserByID))
	})

	r.Route("/reports", func(r chi.Router) {
		r.Use(CrossOriginable)

		r.Route("/{authorId:[0-9]+}", func(r chi.Router) {
			r.Get("/", Handler(SearchReportsByAuthor))
			r.Get("/languages/{language}", Handler(SearchReportsByAuthorAndLanguage))

			r.Get(
				"/_/{year:[0-9]+}/{month:[0-9]+}/{day:[0-9]+}",
				Handler(SearchReportsYmd),
			)

			r.Route("/{id:[0-9]+}", func(r chi.Router) {
				r.Get("/", Handler(GetReport))

				r.With(GAuth).Put("/", Handler(Auth(UpdateReport)))
			})
		})

		r.Route("/", func(r chi.Router) {
			r.Get("/", Handler(GetReports))

			r.With(GAuth).Post("/", Handler(Auth(PostReport)))
		})
	})

	r.Route("/languages", func(r chi.Router) {
		r.Use(CrossOriginable)

		r.Get("/", Handler(GetLanguages))

		r.Route("/{language}", func(r chi.Router) {
			r.Get("/reports", Handler(SearchReportsByLanguage))
		})
	})

	return r
}
