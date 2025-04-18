package sqlite

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"socialNetwork/pkg/logger"
)

// SQLiteStore implements the Storage interface using SQLite
type SQLiteStore struct {
	db     *sql.DB
	logger logger.Logger
	hub    interface{} // For WebSocket notifications
}

// Init verifies the database connection
func (s *SQLiteStore) Init() error {
	s.logger.Info("Initializing database connection")
	return s.db.Ping()
}

// NewSQLiteStore creates a new SQLite storage instance
func NewSQLiteStore() (*SQLiteStore, error) {
	// Get logger for sqlite package
	log := logger.GetLogger("storage/sqlite")
	
	// Get database path from environment variable or use default
	connStr := os.Getenv("DB_PATH")
	if connStr == "" {
		connStr = "./pkg/db/sqlite/main.db"
	}
	log.Info("Opening database connection", "path", connStr)

	db, err := sql.Open("sqlite3", connStr)
	if err != nil {
		log.Error("Failed to open database", "error", err)
		return nil, err
	}

	// Configure connection pool
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)
	
	log.Debug("Database connection pool configured",
		"max_open", 10,
		"max_idle", 5,
		"max_lifetime", "1h")

	return &SQLiteStore{
		db:     db,
		logger: log,
	}, nil
}


// Helper methods
func (s *SQLiteStore) isFollowing(followerID, followedID int64) (bool, error) {
	start := time.Now()
	s.logger.Debug("Checking follow status", 
		"follower_id", followerID, 
		"followed_id", followedID)
	
	query := `SELECT COUNT(*) FROM followers 
              WHERE follower_id = ? AND followed_id = ? AND status = 'accepted'`

	var count int
	err := s.db.QueryRow(query, followerID, followedID).Scan(&count)
	if err != nil {
		s.logger.Error("Failed to check follow status", 
			"error", err, 
			"follower_id", followerID, 
			"followed_id", followedID)
		return false, err
	}

	duration := time.Since(start)
	s.logger.Debug("Follow status checked", 
		"is_following", count > 0,
		"duration_ms", duration.Milliseconds(),
		"follower_id", followerID, 
		"followed_id", followedID)
	
	return count > 0, nil
}

// NOTE: IsFollowing is now implemented in following.go

// DB returns the underlying database connection for direct queries
func (s *SQLiteStore) DB() *sql.DB {
	return s.db
}

// SetHub sets the WebSocket hub for real-time notifications
func (s *SQLiteStore) SetHub(hub interface{}) {
	s.hub = hub
}