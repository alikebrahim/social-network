package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type UserProfile struct {
	UserID      string `json:"user_id"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DateOfBirth string `json:"date_of_birth"`
	Avatar      string `json:"avatar,omitempty"`
	Nickname    string `json:"nickname,omitempty"`
	AboutMe     string `json:"about_me,omitempty"`
	IsPublic    bool   `json:"is_public"`
}

// GetProfileHandler retrieves a user's profile
func GetProfileHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "User ID is required", http.StatusBadRequest)
			return
		}

		// Retrieve user profile from the database
		var profile UserProfile
		query := `SELECT id, email, first_name, last_name, date_of_birth, avatar, nickname, about_me, is_public
				  FROM users WHERE id = ?`
		err := db.QueryRow(query, userID).Scan(&profile.UserID, &profile.Email, &profile.FirstName, &profile.LastName,
			&profile.DateOfBirth, &profile.Avatar, &profile.Nickname, &profile.AboutMe, &profile.IsPublic)
		if err == sql.ErrNoRows {
			http.Error(w, "Profile not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Failed to retrieve profile", http.StatusInternalServerError)
			return
		}

		// Respond with the profile
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(profile)
	}
}

// UpdateProfileHandler updates a user's profile
func UpdateProfileHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("user_id").(string) // Assuming user ID is stored in the context

		var profile UserProfile
		if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Update profile in the database
		query := `UPDATE users SET email = ?, first_name = ?, last_name = ?, date_of_birth = ?, 
				  avatar = ?, nickname = ?, about_me = ?, is_public = ? WHERE id = ?`
		_, err := db.Exec(query, profile.Email, profile.FirstName, profile.LastName, profile.DateOfBirth,
			profile.Avatar, profile.Nickname, profile.AboutMe, profile.IsPublic, userID)
		if err != nil {
			http.Error(w, "Failed to update profile", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated successfully"})
	}
}
