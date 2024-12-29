package middleware

import (
	"net/http"
	"strings"
)

func AllowContentEncoding(allowedEncodings ...string) func(next http.Handler) http.Handler {
	allowed := make(map[string]bool)
	for _, enc := range allowedEncodings {
		allowed[strings.ToLower(enc)] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			contentEncoding := strings.ToLower(r.Header.Get("Content-Encoding"))

			if contentEncoding != "" && !allowed[contentEncoding] {
				// Respond with 415 Unsupported Media Type
				http.Error(w, "Unsupported Media Type: "+contentEncoding, http.StatusUnsupportedMediaType)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
