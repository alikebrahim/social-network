package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

)

type follow struct {
	FollowerID string    `json:"follower_id"`
	FollowedID string    `json:"followed_id"`
	Status     string `json:"status"`
	Followers  []User `json:"followers"`
}

// FOLLOWING HANDLERS
// POST /follow/{id}
func FollowingHandler(w http.ResponseWriter, r *http.Request, DB *sql.DB, followedID string) {
	//var user User
	var follow follow

	follow.FollowedID = followedID

	err := json.NewDecoder(r.Body).Decode(&follow)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	//followedID := follow.FollowedID
	

	sessionToken, err := GetSessionToken(r)
	if err != nil {
		http.Error(w, "Invalid session token", http.StatusUnauthorized)
		log.Print(err)
		return
	}

	log.Println(sessionToken)


	follow.FollowerID, err = getUserIDFromSession(sessionToken)
	if err != nil {
		http.Error(w, "Invalid session token", http.StatusUnauthorized)
		return
	}

	if follow.FollowerID == follow.FollowedID {
		http.Error(w, "Cannot follow yourself", http.StatusBadRequest)
		return
	}

	var profileType string
	err = DB.QueryRow("SELECT profile_type FROM users WHERE id = ?", follow.FollowedID).Scan(&profileType)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		log.Print(err)
		return
	}

	var count int
	err = DB.QueryRow("SELECT COUNT(*) FROM followers WHERE follower_id = ? AND followed_id = ?", follow.FollowerID, follow.FollowedID).Scan(&count)
	if err != nil {
		http.Error(w, "Error checking if already following", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	if count > 0 {
		http.Error(w, "Already following", http.StatusBadRequest)
		log.Print(err)
		return
	}

	status := "pending"
	if profileType == "public" {
		status = "accepted"
	}

	_, err = DB.Exec("INSERT INTO followers (follower_id, followed_id, status) VALUES (?, ?, ?)", follow.FollowerID, follow.FollowedID, status)
	if err != nil {
		http.Error(w, "Error following user", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	response := map[string]string{"message": "Follow request sent"}
	if status == "accepted" {
		response["message"] = "Follow request accepted"
	}
	
	w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}

// GET /follow/requests
func FollowingRequestsHandler(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	fmt.Fprintln(w, "Hello there")

}

// POST /follow/{id}/accept
func FollowAcceptRequestHandler(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

}

// DELETE /follow/{id}
func FollowRejectRequestHandler(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

}
