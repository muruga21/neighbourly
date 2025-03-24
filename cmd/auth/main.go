package main

import (
	"net/http"

	"github.com/muruga21/neighbourly/internal/service"
)

func main() {
	http.HandleFunc("/signup", service.HandleSignup)
	http.ListenAndServe(":5000", nil)
}
