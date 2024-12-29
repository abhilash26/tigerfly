package handlers

import (
	"fmt"
	"net/http"

	"github.com/abhilash26/tigerfly/internal/view"
)

var counter int = 0

func Index(w http.ResponseWriter, r *http.Request) {
	view.RenderTemplate(w, "page/index", counter)
}

func CounterAdd(w http.ResponseWriter, r *http.Request) {
	counter++
	fmt.Fprintf(w, "%d", counter)
}
