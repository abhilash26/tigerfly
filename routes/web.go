package routes

import (
	"net/http"

	"github.com/abhilash26/tigerfly/app/handlers"
	"github.com/abhilash26/tigerfly/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterWeb(r *chi.Mux) {
	r.Route("/", func(r chi.Router) {
		r.Use(middleware.AllowContentType(
			"application/x-www-form-urlencoded",
			"multipart/form-data", "text/html", "text/plain",
		))
		r.Use(middleware.SetContentType("text/html"))
		r.Use(middleware.Compress(6))

		r.NotFound(handlers.NotFound)

		r.Get("/checkhealth", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		// Define your routes below
		r.Get("/", handlers.Index)
		r.Get("/counter-add", handlers.CounterAdd)
	})
}
