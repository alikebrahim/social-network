package main

import (
	"log"
	"net/http"
	"time"

	"backend/config"
	"backend/pkg/api"
	"backend/pkg/db"
	"backend/pkg/middleware"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize the database
	dbConn, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer dbConn.Close()

	// Apply migrations
	if err := db.ApplyMigrations(dbConn); err != nil {
		log.Fatalf("Error applying database migrations: %v", err)
	}

	// Initialize the router
	router := mux.NewRouter()

	// Add middleware
	router.Use(middleware.RequestLogger)
	router.Use(middleware.CORS) // Use middleware.CORS instead of middleware.CORSMiddleware

	// Register API routes
	api.RegisterRoutes(router, dbConn)

	// Start the server
	server := &http.Server{
		Handler:      router,
		Addr:         ":" + cfg.ServerPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Server is running on port %s...", cfg.ServerPort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}
