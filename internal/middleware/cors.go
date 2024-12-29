package middleware

import (
	"net/http"
	"strconv"
	"strings"
)

type CORSOptions struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

func NewCORSOptions(origins, methods, headers, exposedHeaders []string, allowCredentials bool, maxAge int) *CORSOptions {
	return &CORSOptions{
		AllowedOrigins:   origins,
		AllowedMethods:   methods,
		AllowedHeaders:   headers,
		ExposedHeaders:   exposedHeaders,
		AllowCredentials: allowCredentials,
		MaxAge:           maxAge,
	}
}

func CORS(options *CORSOptions) func(http.Handler) http.Handler {
	allowedOrigins := make(map[string]struct{}, len(options.AllowedOrigins))
	for _, origin := range options.AllowedOrigins {
		allowedOrigins[strings.ToLower(origin)] = struct{}{}
	}

	allowedMethods := make(map[string]struct{}, len(options.AllowedMethods))
	for _, method := range options.AllowedMethods {
		allowedMethods[strings.ToUpper(method)] = struct{}{}
	}

	allowedHeaders := make(map[string]struct{}, len(options.AllowedHeaders))
	for _, header := range options.AllowedHeaders {
		allowedHeaders[strings.ToLower(header)] = struct{}{}
	}

	exposedHeaders := make(map[string]struct{}, len(options.ExposedHeaders))
	for _, header := range options.ExposedHeaders {
		exposedHeaders[strings.ToLower(header)] = struct{}{}
	}

	setCORSHeaders := func(w http.ResponseWriter) {
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(options.AllowedMethods, ", "))
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(options.AllowedHeaders, ", "))
		if len(options.ExposedHeaders) > 0 {
			w.Header().Set("Access-Control-Expose-Headers", strings.Join(options.ExposedHeaders, ", "))
		}
		if options.AllowCredentials {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		w.Header().Set("Access-Control-Max-Age", strconv.Itoa(options.MaxAge))
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if r.Method == http.MethodOptions {
				if _, ok := allowedOrigins["*"]; !ok {
					if _, ok := allowedOrigins[strings.ToLower(origin)]; !ok {
						http.Error(w, "Forbidden", http.StatusForbidden)
						return
					}
				}

				setCORSHeaders(w)
				w.WriteHeader(http.StatusOK)
				return
			}

			if _, ok := allowedOrigins["*"]; !ok {
				originLower := strings.ToLower(origin)
				if _, ok := allowedOrigins[originLower]; !ok {
					http.Error(w, "Forbidden: "+origin, http.StatusForbidden)
					return
				}
			}

			setCORSHeaders(w)
			next.ServeHTTP(w, r)
		})
	}
}
