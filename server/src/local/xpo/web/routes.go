package web

import (
	"net/http"

	"github.com/go-chi/chi"
)

func Routes() {

	r := chi.NewRouter()
	r.Use(CrossOriginable)

	r.Get("/", handleRoot)
	r.Get("/loggedin", handleLoggedIn)

	r.Route("/users/me", func(r chi.Router) {
		r.Use(AuthorizedCheckable)

		r.Get("/", Catch(GetMe))
		r.Post("/", Catch(PostMe))
		r.Put("/", Catch(UpdateMe))
	})

	r.Route("/reports", func(r chi.Router) {
		r.With(AuthorizedCheckable).Get(
			"/{authorId:[0-9]+}/_/{year:[0-9]+}/{month:[0-9]+}/{day:[0-9]+}",
			Catch(SearchReportsYmd),
		)

		r.Route("/{authorId:[0-9]+}/{id:[0-9]+}", func(r chi.Router) {
			r.Use(AuthorizedCheckable)

			r.Get("/", Catch(GetReport))
			r.Put("/", Catch(UpdateReport))
		})

		r.Route("/", func(r chi.Router) {
			r.Get("/", Catch(GetReports))
			r.With(AuthorizedCheckable).Post("/", Catch(PostReport))
		})
	})

	http.Handle("/", r)
}
