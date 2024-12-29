package routes

import (
	"net/http"

	"github.com/abhilash26/tigerfly/internal/env"
	"github.com/abhilash26/tigerfly/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterAPI(r *chi.Mux) {
	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/json"))

		corsOptions := middleware.CORSOptions{
			AllowedOrigins:   env.GetSlice("CORS_ALLOWED_ORIGINS", "*"),
			AllowedMethods:   env.GetSlice("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS"),
			AllowedHeaders:   env.GetSlice("CORS_ALLOWED_HEADERS", "Accept,Authorization,Content-Type,X-CSRF-Token"),
			ExposedHeaders:   env.GetSlice("CORS_EXPOSED_HEADERS", "Link"),
			AllowCredentials: env.GetBool("CORS_ALLOW_CREDENTIALS", false),
			MaxAge:           env.GetInt("CORS_MAX_AGE", 300),
		}

		r.Use(middleware.CORS(&corsOptions))

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"message": "Hi from TigerFly!"}`))
		})
	})
}
