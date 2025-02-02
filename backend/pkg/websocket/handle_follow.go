package websocket

import "log"

func handleFollowRequest(msg Message) {
	response := Message{Type: "follow_response"}

	followerID, err := getUserIDFromSession(msg.SessionToken)
	if err != nil {
		log.Printf("Follow request failed: %v", err)
		response.Content = "Authentication failed"
		sendResponseToClients(response)
		return
	}

	followedID := msg.Follow.FollowedID
	log.Print(followedID)
	if followedID == 0 || followedID == followerID {
		response.Content = "Invalid follow request"
		sendResponseToClients(response)
		return
	} 

	var existingStatus string
	err = DB.QueryRow("SELECT status FROM followers WHERE follower_id = ? AND followed_id = ?", followerID, followedID).Scan(&existingStatus)
	if err == nil {
		if existingStatus == "pending" {
			response.Content = "Follow request already sent"
		} else if existingStatus == "accepted" {
			response.Content = "You are already following this user"
		}
		sendResponseToClients(response)
		return
	}

	var profileType string
	err = DB.QueryRow("SELECT profile_type FROM users WHERE id = ?", followedID).Scan(&profileType)
	if err != nil {
		response.Content = "User not found"
		sendResponseToClients(response)
		return
	}

	status := "pending"
	if profileType == "public" {
		status = "accepted"
	}

	_, err = DB.Exec("INSERT INTO followers (follower_id, followed_id, status) VALUES (?, ?, ?)", followerID, followedID, status)
	if err != nil {
		response.Content = "Follow request failed"
		sendResponseToClients(response)
		return
	}

	response.Content = "Follow request sent"
	if status == "accepted" {
		response.Content = "You are now following this user"
	}

	sendResponseToClients(response)
}

func handleFollowResponse(msg Message){
	response := Message{Type: "follow_response"}

	userID, err := getUserIDFromSession(msg.SessionToken)
	if err != nil {
		log.Printf("Follow response failed: %v", err)
		response.Content = "Authentication failed"
		sendResponseToClients(response)
		return
	}


	_, err = DB.Exec("UPDATE followers SET status = ? WHERE followed_id = ? AND follower_id = ? AND status = 'pending'",
		msg.Follow.Status, userID, msg.Follow.FollowerID)
	if err != nil {
		response.Content = "Follow request update failed"
		sendResponseToClients(response)
		return
	}
	response.Content = "Follow request " + msg.Follow.Status
	sendResponseToClients(response)
}


func handleUnfollow(msg Message) {
	response := Message{Type: "unfollow_response"}


	followerID, err := getUserIDFromSession(msg.SessionToken)
	if err != nil {
		log.Printf("Unfollow failed: %v", err)
		response.Content = "Authentication failed"
		sendResponseToClients(response)
		return
	}

	followedID := msg.Follow.FollowerID
	if followedID == 0 {
		response.Content = "Invalid unfollow request"
		sendResponseToClients(response)
		return
	}
	_, err = DB.Exec("DELETE FROM followers WHERE follower_id = ? AND followed_id = ?", followerID, followedID)
	if err != nil {
		response.Content = "Unfollow failed"
		sendResponseToClients(response)
		return
	}

	response.Content = "Unfollowed successfully"
	sendResponseToClients(response)

}

func handleGetFollowers(msg Message) {
	response := Message{Type: "get_followers_response"}

	userID, err := getUserIDFromSession(msg.SessionToken)
	if err != nil {
		log.Printf("Get followers failed: %v", err)
		response.Content = "Authentication failed"
		sendResponseToClients(response)
		return
	}

	rows, err := DB.Query(`
		SELECT users.id, users.nickname, users.avatar, users.first_name, users.last_name 
		FROM followers 
		JOIN users ON followers.follower_id = users.id 
		WHERE followers.followed_id = ? AND followers.status = 'accepted'`, userID)
	if err != nil {
		response.Content = "Failed to retrieve followers"
		sendResponseToClients(response)
		return
	}
	defer rows.Close()

	var followers []User
	for rows.Next() {
		var follower User
		err := rows.Scan(&follower.ID, &follower.Nickname, &follower.Avatar, &follower.FirstName, &follower.LastName)
		if err != nil {
			log.Println("Error scanning followers:", err)
			continue
		}
		followers = append(followers, follower)
	}

	response.Follow.Followers = followers
	sendResponseToClients(response)
}
