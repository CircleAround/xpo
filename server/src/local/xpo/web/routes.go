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
		r.Use(GAuth)

		r.Get("/", Catch(GetMe))
		r.Post("/", Catch(PostMe))
		r.Put("/", Catch(UpdateMe))
	})

	r.Route("/reports", func(r chi.Router) {
		r.With(GAuth).Get(
			"/{authorId:[0-9]+}/_/{year:[0-9]+}/{month:[0-9]+}/{day:[0-9]+}",
			Catch(SearchReportsYmd),
		)

		r.Route("/{authorId:[0-9]+}/{id:[0-9]+}", func(r chi.Router) {
			r.Get("/", Catch(GetReport))

			r.With(GAuth).Put("/", Catch(Auth(UpdateReport)))
		})

		r.Route("/", func(r chi.Router) {
			r.Get("/", Catch(GetReports))

			r.With(GAuth).Post("/", Catch(Auth(PostReport)))
		})
	})

	http.Handle("/", r)
}
