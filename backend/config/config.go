package config

import (
	"log"
	"os"
)

type Config struct {
	ServerPort  string
	DatabaseURL string
	JWTSecret   string
	StaticFiles string
}

func LoadConfig() *Config {
	// Load environment variables
	serverPort := getEnv("SERVER_PORT", "8080")
	databaseURL := getEnv("DATABASE_URL", "file:app.db?cache=shared&mode=rwc")
	jwtSecret := getEnv("JWT_SECRET", "supersecretkey")
	staticFiles := getEnv("STATIC_FILES_PATH", "./static")

	// Log the loaded configuration
	log.Printf("Loaded configuration: ServerPort=%s, DatabaseURL=%s, StaticFiles=%s",
		serverPort, databaseURL, staticFiles)

	return &Config{
		ServerPort:  serverPort,
		DatabaseURL: databaseURL,
		JWTSecret:   jwtSecret,
		StaticFiles: staticFiles,
	}
}

// getEnv reads an environment variable and provides a default value if not set.
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
