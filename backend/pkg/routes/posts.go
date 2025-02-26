package routes

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type post struct {
	ID        int64    `json:"id"`
	UserID    int64   `json:"user_id"`
	Content   string `json:"content"`
	Image     string `json:"image"`
	Privacy  string `json:"privacy"`
	CreatedAt string `json:"created_at"`
}

// POSTS HANDLERS
// POST /posts
func PostCreateHandler(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	var post post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Print(err)
		return
	}

	sessionToken, err := GetSessionToken(r)
	if err != nil {
		http.Error(w, "Invalid session token", http.StatusUnauthorized)
		log.Print(err)
		return
	}

	post.UserID, err = getUserIDFromSession(sessionToken)
	if err != nil {
		http.Error(w, "Invalid session token", http.StatusUnauthorized)
		return
	}
	
	if post.Content == "" && post.Image == "" {
		http.Error(w, "Post cannot be empty", http.StatusBadRequest)
		log.Print("Post cannot be empty")
		return
	}

	if post.Privacy == "" {
		post.Privacy = "public"
	}else if post.Privacy != "public" && post.Privacy != "private"  {
		http.Error(w, "Invalid privacy setting", http.StatusBadRequest)
		log.Print("Invalid privacy setting")
		return
	}

	_, err = DB.Exec("INSERT INTO posts (user_id, content, image, privacy, created_at) VALUES (?, ?, ?, ?, ?)",
		post.UserID, post.Content, post.Image, post.Privacy, post.CreatedAt)
	if err != nil {
		http.Error(w, "Post creation failed", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	log.Print("Post created successfully")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "post created successfully"})
}

// GET /posts/{id}
func PostsGetHandler(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

}

// PUT /posts/{id}
func PostEditHandler(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

}

// DELETE /posts/{id}
func PostDeleteHandler(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

}

// POST /posts/{id}/comments
func PostCommentHandler(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

}

// POST /posts/{id}/likes
func PostLikeHandler(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

}
