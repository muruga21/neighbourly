package handlers

import (
	"fmt"
	"net/http"
)

// HealthHandler returns a simple string to signify the service is alive.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "im alive")
}
