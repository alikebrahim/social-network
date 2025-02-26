package routes

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type follow struct {
	FollowerID int64  `json:"follower_id"`
	FollowedID int64  `json:"followed_id"`
	Status     string `json:"status"`
	Followers  []User `json:"followers"`
}

// FOLLOWING HANDLERS
// POST /follow/{id}
func  FollowingHandler(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	//var user User
	var follow follow

	vars := mux.Vars(r)
	followedID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		log.Print(err)
		return
	}
	follow.FollowedID = followedID

	err = json.NewDecoder(r.Body).Decode(&follow)
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
	var user User
	sessiomToken, err := GetSessionToken(r)
	if err != nil {
		http.Error(w, "Invalid session token", http.StatusUnauthorized)
		log.Print(err)
		return
	}
	user.ID, err = getUserIDFromSession(sessiomToken)
	if err != nil {
		http.Error(w, "Invalid session token", http.StatusUnauthorized)
		log.Print(err)
		return
	}
	rows, err := DB.Query(`
	SELECT f.follower_id, u.first_name, u.last_name, u.nickname, u.avatar
		FROM followers f
		JOIN users u ON f.follower_id = u.id  -- Fixed JOIN condition
		WHERE f.followed_id = ? AND f.status = 'pending'
`, user.ID)
	if err != nil {
		http.Error(w, "Error getting follower requests", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	var followRequests []User
	for rows.Next() {
		var follower User
		err := rows.Scan(&follower.ID, &follower.FirstName, &follower.LastName, &follower.Nickname, &follower.Avatar)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		followRequests = append(followRequests, follower)
	}

	if len(followRequests) == 0 {
		log.Println("No follow requests found")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]User{})
		return
	}
	log.Println("Pending follow requests retrieved:", followRequests)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(followRequests)
}

// POST /follow/{id}/accept
func FollowAcceptRequestHandler(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	vars := mux.Vars(r)
	followedId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		log.Print(err)
		return
	}
	sessionToken, err := GetSessionToken(r)
	if err != nil {
		http.Error(w, "Invalid session token", http.StatusUnauthorized)
		log.Print(err)
		return
	}

	userID, err := getUserIDFromSession(sessionToken)
	if err != nil {
		http.Error(w, "Invalid session token", http.StatusUnauthorized)
		log.Print(err)
		return
	}

	var count int
	err = DB.QueryRow("SELECT COUNT(*) FROM followers WHERE follower_id = ? AND followed_id = ? AND status = 'pending'", followedId, userID).Scan(&count)
	if err != nil {
		http.Error(w, "Error checking if follow request exists", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	if count == 0 {
		http.Error(w, "No follow request found", http.StatusNotFound)
		log.Print(err)
		return
	}

	_, err = DB.Exec(`
		UPDATE followers 
		SET status = 'accepted' 
		WHERE follower_id = ? AND followed_id = ? AND status = 'pending'
	`, followedId, userID)
	if err != nil {
		http.Error(w, "Error accepting follow request", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	log.Println("Follow request accepted")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Follow request accepted"})
}

// DELETE /follow/{id}
func FollowRejectRequestHandler(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	log.Println("Rejecting follow request")
	vars := mux.Vars(r)
	followerId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		log.Print(err)
		return
	}
	sessionToken, err := GetSessionToken(r)
	if err != nil {
		http.Error(w, "Invalid session token", http.StatusUnauthorized)
		log.Print(err)
		return
	}

	followedId, err := getUserIDFromSession(sessionToken)
	if err != nil {
		http.Error(w, "Invalid session token", http.StatusUnauthorized)
		log.Print(err)
		return
	}


	var count int
	err = DB.QueryRow("SELECT COUNT(*) FROM followers WHERE follower_id = ? AND followed_id = ? AND status = 'pending'", followerId, followedId).Scan(&count)
	if err != nil {
		http.Error(w, "Error checking if follow request exists", http.StatusInternalServerError)
		log.Print("gg")
		log.Print(err)
		return
	}
	log.Println("Count:", count)
	if count == 0 {
		http.Error(w, "No follow request found", http.StatusNotFound)
		log.Print(err)
		return
	}

	_, err = DB.Exec("DELETE FROM followers WHERE follower_id = ? AND followed_id = ? AND status = 'pending'", followerId, followedId)
	if err != nil {
		http.Error(w, "Error rejecting follow request", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	log.Println("Follow request rejected")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Follow request rejected"})
}
