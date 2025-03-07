package main

import (
	"encoding/json"
	"net/http"
)

// POSTS HANDLERS
// POST /posts
func (s *APIServer) HandlePostCreate(w http.ResponseWriter, r *http.Request) error {
	var post Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		return err
	}

	// Get user ID from session
	userID, err := s.getUserIDFromSession(r)
	if err != nil {
		return ErrUnauthorized
	}

	// Set the user ID
	post.UserID = userID
	
	id, err := s.store.CreatePost(&post)
	if err != nil {
		return err
	}

	// Set the ID in the response
	post.ID = id
	
	WriteJson(w, http.StatusCreated, post)
	return nil
}

// GET /posts/{id}
func (s *APIServer) HandlePostsGet(w http.ResponseWriter, r *http.Request) error {
	return nil

}

// PUT /posts/{id}
func (s *APIServer) HandlePostEdit(w http.ResponseWriter, r *http.Request) error {
	return nil

}

// DELETE /posts/{id}
func (s *APIServer) HandlePostDelete(w http.ResponseWriter, r *http.Request) error {
	return nil

}

// POST /posts/{id}/comments
func (s *APIServer) HandlePostComment(w http.ResponseWriter, r *http.Request) error {
	return nil

}

// POST /posts/{id}/likes
func (s *APIServer) HandlePostLike(w http.ResponseWriter, r *http.Request) error {
	return nil

}
