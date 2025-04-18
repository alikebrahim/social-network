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

	// Get log file path from environment or use default
	logFilePath := os.Getenv("LOG_FILE")
	if logFilePath == "" {
		logFilePath = filepath.Join(logsDir, "app.log")
	}
	toConsole := os.Getenv("LOG_CONSOLE") == "true" || os.Getenv("LOG_CONSOLE") == ""
	useColors := os.Getenv("LOG_COLORS") != "false" // Use colors by default unless explicitly turned off

	logger.Init(logger.Config{
		Level:     logger.LevelFromString(logLevel),
		LogFile:   logFilePath,
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

	// Get server host and port from environment or use defaults
	serverHost := os.Getenv("SERVER_HOST")
	if serverHost == "" {
		serverHost = "0.0.0.0"
	}
	
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "3000"
	}
	
	// Create the listen address
	listenAddr := fmt.Sprintf("%s:%s", serverHost, serverPort)
	
	// Create and run the API server
	server := api.NewAPIServer(listenAddr, store)
	log.Info("Starting server", "address", listenAddr)
	server.Run()
}