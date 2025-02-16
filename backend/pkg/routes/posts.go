package routes

import (
	"database/sql"
	"net/http"
)

// POSTS HANDLERS
// POST /posts
func PostCreateHandler(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

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
