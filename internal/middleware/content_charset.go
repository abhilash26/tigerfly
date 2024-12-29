package middleware

import (
	"net/http"
	"strings"
)

func ContentCharset(charsets ...string) func(next http.Handler) http.Handler {
	validCharsets := make(map[string]bool, len(charsets))
	for _, charset := range charsets {
		validCharsets[strings.ToLower(charset)] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !isValidCharset(r.Header.Get("Content-Type"), validCharsets) {
				w.WriteHeader(http.StatusUnsupportedMediaType)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func isValidCharset(contentType string, validCharsets map[string]bool) bool {
	_, params := split(contentType, ";")
	_, charset := split(params, "charset=")

	charset = strings.TrimSpace(charset)
	_, exists := validCharsets[charset]
	return exists
}

func split(str, sep string) (string, string) {
	parts := strings.SplitN(str, sep, 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[0]), parts[1]
	}
	return str, ""
}
