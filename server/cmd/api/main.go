package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"neighbourly/server/internal/database"
	"neighbourly/server/internal/middleware"
	"neighbourly/server/internal/server"
)

func main() {
	// Setup log file to write to both Terminal and server.log
	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		multiWriter := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(multiWriter)
	} else {
		log.Println("Failed to open server.log. Proceeding with terminal logging only.")
	}

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading, relying on system environment variables")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable is not set")
	}

	// Connect to Database
	if err := database.Connect(mongoURI); err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer database.Disconnect()

	// Setup Router
	mux := server.NewRouter()

	// Wrap the router with our Logger middleware
	loggedRouter := middleware.Logger(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	fmt.Printf("Starting server on port :%s\n", port)

	if err := http.ListenAndServe("0.0.0.0:"+port, loggedRouter); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
