package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/abhilash26/tigerfly/app/handlers"
)

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic occurred: %v\n%s", err, debug.Stack())

				// Respond with a generic 500 Internal Server Error
				handlers.InternalServerError(w, r)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
