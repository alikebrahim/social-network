package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// NOTE: migration handling
// goose -dir migrationsDir sqlite3 up
// goose -dir migrationsDir sqlite3 down

type Storage interface {
	// test data
	AddTestAccount() error
	// USER ACCOUNTS
	CreateUserAccount(*UserAccount) (id int64, err error)
	GetSeesionToken(int64) (string, error)
	AuthenticateUser(string, string) (string, error)
	DeleteSession(string) error
	GetUserIdBySession(string) (int64, error)

	// follow requests
	CreateFollowRequest(follow) error
	AcceptFollowRequest(int64, int64) error
	DeleteFollowRequest(int64, int64) error
	GetFollowRequests(int64) ([]UserAccount, error)
	
	

	// DeleteUserAccount(int) error
	// EditUserAccount(*UserAccount) error
	// GetUserAccountByID(int) (*UserAccount, error)
	// GetUserAccounts() ([]*UserAccount, error)

	// POSTS
	CreatePost(Post) (id int64, err error)
}

type SQLiteStore struct {
	db *sql.DB
}

func (s *SQLiteStore) Init() error {
	return s.db.Ping()
}

func NewSQLiteStore() (*SQLiteStore, error) {
	connStr := "./pkg/db/sqlite/database.sql"

	db, err := sql.Open("sqlite3", connStr)
	if err != nil {
		return nil, err
	}

	return &SQLiteStore{
		db: db,
	}, nil
}
