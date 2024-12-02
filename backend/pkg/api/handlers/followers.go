package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type FollowRequest struct {
	TargetUserID string `json:"target_user_id"`
}

type FollowResponse struct {
	Status string `json:"status"`
}

// FollowHandler handles sending a follow request or directly following a user with a public profile
func FollowHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req FollowRequest

		// Parse JSON body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Get target user's profile privacy status
		var isPublic bool
		query := `SELECT is_public FROM users WHERE id = ?`
		err := db.QueryRow(query, req.TargetUserID).Scan(&isPublic)
		if err == sql.ErrNoRows {
			http.Error(w, "Target user not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Error retrieving target user", http.StatusInternalServerError)
			return
		}

		// Add follow entry to the database
		if isPublic {
			// Auto-follow for public profiles
			query = `INSERT INTO followers (follower_id, following_id, status) VALUES (?, ?, 'accepted')`
			_, err = db.Exec(query, r.Context().Value("user_id"), req.TargetUserID)
			if err != nil {
				http.Error(w, "Error following user", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(FollowResponse{Status: "followed"})
		} else {
			// Follow request for private profiles
			query = `INSERT INTO followers (follower_id, following_id, status) VALUES (?, ?, 'pending')`
			_, err = db.Exec(query, r.Context().Value("user_id"), req.TargetUserID)
			if err != nil {
				http.Error(w, "Error sending follow request", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(FollowResponse{Status: "request_sent"})
		}
	}
}

// AcceptFollowRequestHandler handles accepting a follow request
func AcceptFollowRequestHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		followerID := r.URL.Query().Get("follower_id")
		if followerID == "" {
			http.Error(w, "Follower ID is required", http.StatusBadRequest)
			return
		}

		// Update the follow request to accepted
		query := `UPDATE followers SET status = 'accepted' WHERE follower_id = ? AND following_id = ? AND status = 'pending'`
		result, err := db.Exec(query, followerID, r.Context().Value("user_id"))
		if err != nil {
			http.Error(w, "Error updating follow request", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "Follow request not found or already accepted", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(FollowResponse{Status: "request_accepted"})
	}
}

// UnfollowHandler handles unfollowing a user
func UnfollowHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req FollowRequest

		// Parse JSON body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Delete the follow entry from the database
		query := `DELETE FROM followers WHERE follower_id = ? AND following_id = ?`
		result, err := db.Exec(query, r.Context().Value("user_id"), req.TargetUserID)
		if err != nil {
			http.Error(w, "Error unfollowing user", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "Follow relationship not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(FollowResponse{Status: "unfollowed"})
	}
}
