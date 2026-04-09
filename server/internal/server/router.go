package server

import (
	"net/http"

	"neighbourly/server/internal/handlers"
)

// NewRouter creates a new HTTP router and registers all routes.
func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/", handlers.HelloHandler)
	mux.HandleFunc("/health", handlers.HealthHandler)
	mux.HandleFunc("/api/auth/signup", handlers.SignupHandler)
	mux.HandleFunc("/api/auth/login", handlers.LoginHandler)
	mux.HandleFunc("/api/providers", handlers.ListProvidersHandler)
	mux.HandleFunc("/api/profile/update", handlers.UpdateProfileHandler)
	mux.HandleFunc("/api/profile/{userid}", handlers.ProfileHandler)

	return mux
}
