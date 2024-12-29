package middleware

import (
	"net/http"
	"path"
	"strings"
)

// Clean double slashes from the path (e.g., "/users//1" -> "/users/1")
// Strip the trailing slash if exists (e.g., "/users/1/" -> "/users/1")

func CleanSlashes(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimSuffix(path.Clean(r.URL.Path), "/")

		// Continue processing the request
		next.ServeHTTP(w, r)
	})
}
