package router

import (
	"net/http"
	"os"

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
	fsys := os.DirFS(dir)

	r.Get(path+"/*", func(w http.ResponseWriter, r *http.Request) {
		filePath := r.URL.Path[len(path):]

		http.ServeFileFS(w, r, fsys, filePath)
	})
}
