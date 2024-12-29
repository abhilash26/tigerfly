package router

import (
	"net/http"

	"github.com/abhilash26/tigerfly/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func New() *chi.Mux {
	r := chi.NewRouter()

	// Essentials Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CleanSlashes)
	return r
}

func Static(r *chi.Mux, path string, dir string) {
	fs := http.FileServer(http.Dir(dir))
	r.Get(path+"/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix(path, fs).ServeHTTP(w, r)
	})
}
