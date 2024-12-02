package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type CreateGroupRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type InviteToGroupRequest struct {
	GroupID string `json:"group_id"`
	UserID  string `json:"user_id"`
}

type JoinGroupRequest struct {
	GroupID string `json:"group_id"`
}

type CreateEventRequest struct {
	GroupID     string `json:"group_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DayTime     string `json:"day_time"` // Format: YYYY-MM-DD HH:MM:SS
}

type Response struct {
	Message string `json:"message"`
}

// CreateGroupHandler handles group creation
func CreateGroupHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateGroupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		groupID := uuid.New().String()
		query := `INSERT INTO groups (id, title, description, created_at) VALUES (?, ?, ?, ?)`
		_, err := db.Exec(query, groupID, req.Title, req.Description, time.Now())
		if err != nil {
			http.Error(w, "Failed to create group", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Response{Message: "Group created successfully"})
	}
}

// InviteToGroupHandler handles inviting a user to a group
func InviteToGroupHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req InviteToGroupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		query := `INSERT INTO group_invitations (group_id, user_id, status, created_at) VALUES (?, ?, 'pending', ?)`
		_, err := db.Exec(query, req.GroupID, req.UserID, time.Now())
		if err != nil {
			http.Error(w, "Failed to send group invitation", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{Message: "Invitation sent successfully"})
	}
}

// JoinGroupHandler handles a user joining a group
func JoinGroupHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req JoinGroupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		query := `INSERT INTO group_members (group_id, user_id, joined_at) VALUES (?, ?, ?)`
		_, err := db.Exec(query, req.GroupID, r.Context().Value("user_id"), time.Now())
		if err != nil {
			http.Error(w, "Failed to join group", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{Message: "Joined group successfully"})
	}
}

// CreateEventHandler handles creating an event in a group
func CreateEventHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateEventRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		eventID := uuid.New().String()
		query := `INSERT INTO group_events (id, group_id, title, description, day_time, created_at) VALUES (?, ?, ?, ?, ?, ?)`
		_, err := db.Exec(query, eventID, req.GroupID, req.Title, req.Description, req.DayTime, time.Now())
		if err != nil {
			http.Error(w, "Failed to create event", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Response{Message: "Event created successfully"})
	}
}

// GetGroupHandler handles the request to get a group
func GetGroupHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Your handler logic here
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("GetGroupHandler response"))
	}
}
