package handlers

import (
	"net/http"

	"github.com/abhilash26/tigerfly/views/page"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	errorPage := page.Error(404, "The page you are looking for does not exist", true)
	errorPage.Render(r.Context(), w)
}

func InternalServerError(w http.ResponseWriter, r *http.Request) {
	errorPage := page.Error(500, "Internal Server Error", true)
	errorPage.Render(r.Context(), w)
}
