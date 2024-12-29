package middleware

import (
	"net/http"
)

func SetContentType(contentType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", contentType)

			next.ServeHTTP(w, r)
		})
	}
}
