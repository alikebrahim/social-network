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

	id, err := s.store.CreatePost(post)
	if err != nil {
		return err
	}

	WriteJson(w, http.StatusCreated, id)

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
