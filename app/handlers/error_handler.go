package handlers

import (
	"net/http"

	"github.com/abhilash26/tigerfly/internal/view"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	view.Render404(w)
}

func InternalServerError(w http.ResponseWriter, r *http.Request) {
	view.Render500(w)
}
