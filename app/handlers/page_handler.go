package handlers

import (
	"fmt"
	"net/http"

	"github.com/abhilash26/tigerfly/views/page"
)

var counter int = 0

func Index(w http.ResponseWriter, r *http.Request) {
	indexPage := page.Index(counter)
	indexPage.Render(r.Context(), w)
}

func CounterAdd(w http.ResponseWriter, r *http.Request) {
	counter++
	fmt.Fprintf(w, "%d", counter)
}
