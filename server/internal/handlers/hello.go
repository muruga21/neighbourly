package handlers

import (
	"fmt"
	"net/http"
)

// HelloHandler handles the root path request.
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Hello, World!")
}
