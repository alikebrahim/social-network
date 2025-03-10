package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// POSTS HANDLERS
// POST /posts
func (s *APIServer) HandlePostCreate(w http.ResponseWriter, r *http.Request) error {
	var post Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		return err
	}

	sessionToken, err := getSessionToken(r)
	if err != nil {
		log.Println("getSessionToken error :", err)
		return err
	}

	post.UserID, err = s.store.GetUserIdBySession(sessionToken)
	if err != nil {
		log.Println("store.GetUserIDBySession error :", err)
		return err
	}

	if post.Content == "" && post.Image == "" {
		err = errors.New("Post must have content or image")
		log.Println("Post must have content or image")
		return err
	}

	if post.Privacy == "" {
		post.Privacy = "public"
	}else if post.Privacy != "public" && post.Privacy != "private" && post.Privacy != "friends" {
		err = errors.New("privacy must be public, private or friends")
		log.Println("Privacy must be public, private or friends")
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
