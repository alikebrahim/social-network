package websocket

import "log"

func handleGetPosts(msg Message) {
	response := Message{Type: "get_posts_response"}

	userID, _ := getUserIDFromSession(msg.SessionToken)

	rows, err := DB.Query(`
		SELECT posts.id, posts.user_id, users.profile_type, posts.content, posts.image, posts.privacy, posts.created_at 
		FROM posts 
		JOIN users ON posts.user_id = users.id
		WHERE posts.privacy = 'public' 
		OR (posts.privacy = 'private' AND posts.user_id IN 
		    (SELECT followed_id FROM followers WHERE follower_id = ? AND status = 'accepted'))
		ORDER BY posts.created_at DESC`, userID)

	if err != nil {
		response.Content = "Failed to retrieve posts"
		sendResponseToClients(response)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var profileType string

		err := rows.Scan(&post.ID, &post.UserID, &profileType, &post.Content, &post.Image, &post.Privacy, &post.CreatedAt)
		if err != nil {
			log.Println("Error scanning post:", err)
			continue
		}

		if profileType == "private" {
			isFollower := false

			err := DB.QueryRow("SELECT 1 FROM followers WHERE follower_id = ? AND followed_id = ? AND status = 'accepted'",
				userID, post.UserID).Scan(&isFollower)

			if err != nil {
				continue 
			}
		}

		posts = append(posts, post)
	}

	response.Posts = posts
	sendResponseToClients(response)
}
