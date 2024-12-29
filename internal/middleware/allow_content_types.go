package middleware

import (
	"net/http"
	"strings"
)

func AllowContentType(allowedTypes ...string) func(next http.Handler) http.Handler {
	allowed := make(map[string]bool, len(allowedTypes))
	for _, t := range allowedTypes {
		allowed[strings.ToLower(t)] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			contentType := strings.ToLower(r.Header.Get("Content-Type"))

			if idx := strings.IndexByte(contentType, ';'); idx != -1 {
				contentType = contentType[:idx]
			}

			if contentType != "" && !allowed[contentType] {
				http.Error(w, "Unsupported Media Type: "+contentType, http.StatusUnsupportedMediaType)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
