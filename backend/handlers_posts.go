package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
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
	} else if post.Privacy != "public" && post.Privacy != "private" {
		err = errors.New("privacy must be public or private")
		log.Println("Privacy must be public or private")
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

// toDo fix the for loop in s.store.GetPostByID
func (s *APIServer) HandlePostsGet(w http.ResponseWriter, r *http.Request) error {

	id := r.PathValue("id")
	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println("strconv error :", err)
		return err
	}
	sessionToken, err := getSessionToken(r)
	if err != nil {
		log.Println("getSessionToken error :", err)
		return err
	}
	followerID, err := s.store.GetUserIdBySession(sessionToken)
	if err != nil {
		log.Println("store.GetUserIDBySession error :", err)
		return err
	}
	posts, err := s.store.GetPostByID(userID, followerID)
	if err != nil {
		log.Println("store.GetPost error :", err)
		return err
	}
	err = WriteJson(w, http.StatusOK, posts)
	if err != nil {
		log.Println("WriteJson error :", err)
		return err
	}
	return nil

}

// PUT /posts/{id}
func (s *APIServer) HandlePostEdit(w http.ResponseWriter, r *http.Request) error {
	var post Post
	var err error

	// Get post ID from URL
	post.ID, err = strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		log.Println("strconv error :", err)
		return err
	}

	// Extract session token from cookies
	session, err := getSessionToken(r)
	if err != nil {
		log.Println("getSessionToken error :", err)
		return err
	}

	
	post.UserID, err = s.store.GetUserIdBySession(session)
	if err != nil {
		log.Println("store.GetUserIdBySession error :", err)
		return err
	}

	isOwner, err := s.store.IsPostOwner(post.ID, post.UserID)
	if err != nil {
		log.Println("store.IsPostOwner error :", err)
		return err
	}
	if !isOwner {
		err = errors.New("You are not the owner of this post")
		log.Println(err)
		return err
	}

	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		log.Println("json.Decode error :", err)
		return err
	}

	log.Println("Updated post data:", post)

	err = s.store.EditPost(post)
	if err != nil {
		log.Println("EditPost error :", err)
		return err
	}

	log.Println("Post edited successfully")
	return WriteJson(w, http.StatusOK, post)
}

// DELETE /posts/{id}
func (s *APIServer) HandlePostDelete(w http.ResponseWriter, r *http.Request) error {
	var post Post
	var err error
	post.ID, err = strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		log.Println("strconv error :", err)
		return err
	}
	session, err := getSessionToken(r)
	if err != nil {
		log.Println("getSessionToken error :", err)
		return err
	}
	post.UserID, err = s.store.GetUserIdBySession(session)
	if err != nil {
		log.Println("store.GetUserIDBySession error :", err)
		return err
	}
	isOwner, err := s.store.IsPostOwner(post.ID, post.UserID)
	if err != nil {
		log.Println("store.IsPostOwner error :", err)
		return err
	}
	if !isOwner {
		err = errors.New("You are not the owner of this post")
		log.Println("You are not the owner of this post")
		return err
	} else {
		s.store.DeletePost(post.ID)
		return WriteJson(w, http.StatusOK, post)
	}
	return nil

}

// POST /posts/{id}/comments
func (s *APIServer) HandlePostComment(w http.ResponseWriter, r *http.Request) error {
	var comment Comment
	var err error
	comment.PostID, err = strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		log.Println("strconv error :", err)
		return err
	}
	session, err := getSessionToken(r)
	if err != nil {
		log.Println("getSessionToken error :", err)
		return err
	}
	comment.ComenterID, err = s.store.GetUserIdBySession(session)
	if err != nil {
		log.Println("store.GetUserIDBySession error :", err)
		return err
	}
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		log.Println("json.Decode error :", err)
		return err
	}
	err = s.store.CreateComment(comment)
	if err != nil {
		log.Println("store.CreateComment error :", err)
		return err
	}
	WriteJson(w, http.StatusCreated, comment)
	return nil
}

// POST /posts/{id}/likes
func (s *APIServer) HandlePostLike(w http.ResponseWriter, r *http.Request) error {
	var like like
	err := json.NewDecoder(r.Body).Decode(&like)
	if err != nil {
		log.Println("json.Decode error :", err)
		return err
	}
	id := r.PathValue("id")
	like.PostID, err = strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println("strconv error :", err)
		return err
	}
	sessionToken, err := getSessionToken(r)
	if err != nil {
		log.Println("getSessionToken error :", err)
		return err
	}
	like.UserID, err = s.store.GetUserIdBySession(sessionToken)
	if err != nil {
		log.Println("store.GetUserIDBySession error :", err)
		return err
	}
	
	err = s.store.CreateLike(like)
	if err != nil {
		log.Println("store.CreateLike error :", err)
		return err
	}
	WriteJson(w, http.StatusCreated, like)
	return nil

}

func (s *APIServer) HandlePostUnlike(w http.ResponseWriter, r *http.Request) error {
	var like like
	//err := json.NewDecoder(r.Body).Decode(&like)
	var err error
	if err != nil {
		log.Println("json.Decode error :", err)
		return err
	}
	id := r.PathValue("id")
	like.PostID, err = strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println("strconv error :", err)
		return err
	}
	sessionToken, err := getSessionToken(r)
	if err != nil {
		log.Println("getSessionToken error :", err)
		return err
	}
	like.UserID, err = s.store.GetUserIdBySession(sessionToken)
	if err != nil {
		log.Println("store.GetUserIDBySession error :", err)
		return err
	}

	err = s.store.RemoveLikes(like)
	if err != nil {
		log.Println("store.CreateLike error :", err)
		return err
	}
	
	return nil
}
