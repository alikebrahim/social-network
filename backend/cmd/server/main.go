package main

import (
	"fmt"
	"os"
	"path/filepath"

	"socialNetwork/pkg/api"
	"socialNetwork/pkg/logger"
	"socialNetwork/pkg/storage/sqlite"
)

func main() {
	// Initialize the logger
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "INFO"
	}

	// Ensure logs directory exists
	logsDir := filepath.Join(".", "logs")
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logsDir, 0755); err != nil {
			panic("Failed to create logs directory: " + err.Error())
		}
	}

	logFile := filepath.Join(logsDir, "app.log")
	toConsole := os.Getenv("LOG_CONSOLE") == "true" || os.Getenv("LOG_CONSOLE") == ""
	useColors := os.Getenv("LOG_COLORS") != "false" // Use colors by default unless explicitly turned off

	logger.Init(logger.Config{
		Level:     logger.LevelFromString(logLevel),
		LogFile:   logFile,
		ToConsole: toConsole,
		UseColors: useColors,
	})

	// Close logger on exit
	defer func() {
		if err := logger.GetLogger("main").Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing logger: %v\n", err)
		}
	}()

	// Get logger for main package
	log := logger.GetLogger("main")
	log.Info("Application starting with colored output")

	// Initialize the SQLite storage
	store, err := sqlite.NewSQLiteStore()
	if err != nil {
		log.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	if err := store.Init(); err != nil {
		log.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	log.Info("Database initialized successfully")

	// Create and run the API server
	server := api.NewAPIServer(":3000", store)
	log.Info("Starting server", "address", ":3000")
	server.Run()
}