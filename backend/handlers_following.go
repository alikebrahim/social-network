package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
)

// FOLLOWING HANDLERS
// POST /follow/{id}
func (s *APIServer) HandleFollowing(w http.ResponseWriter, r *http.Request) error {
	var follow follow
	id := r.PathValue("id")
	followedID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println("strconv error :", err)
		return err
	}
	follow.FollowedID = followedID
	log.Print("gg")
	log.Print(r.Cookie("session_token"))
	session, err := getSessionToken(r)
	if err != nil {
		log.Println("getSessionToken error :", err)
		return err
	}

	follow.FollowerID, err = s.store.GetUserIdBySession(session)
	if err != nil {
		log.Println("store.GetUserIDBySession error :", err)
		return err
	}

	if follow.FollowerID == follow.FollowedID {
		err = errors.New("You can't follow yourself")
		log.Println("You can't follow yourself")
		return err
	}

	err = s.store.CreateFollowRequest(follow)
	if err != nil {
		log.Println("store.CreateFollowRequest error :", err)
		return err
	}

	return nil

}

// GET /follow/requests
func (s *APIServer) HandleFollowingRequests(w http.ResponseWriter, r *http.Request) error {
	var follow follow
	session, err := getSessionToken(r)
	if err != nil {
		log.Println("getSessionToken error :", err)
		return err
	}
	follow.FollowerID, err = s.store.GetUserIdBySession(session)
	if err != nil {
		log.Println("store.GetUserIDBySession error :", err)
		return err
	}
	followers, err := s.store.GetFollowRequests(follow.FollowerID)
	if err != nil {
		log.Println("store.GetFollowRequests error :", err)
		return err
	}
	follow.Followers = followers
	err = WriteJson(w, http.StatusOK, follow)
	if err != nil {
		log.Println("WriteJson error :", err)
		return err
	}

	return nil
}

// POST /follow/{id}/accept
func (s *APIServer) HandleFollowingAccept(w http.ResponseWriter, r *http.Request) error {
	var follow follow
	id := r.PathValue("id")
	follower_id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println("strconv error :", err)
		return err
	}
	follow.FollowerID = follower_id
	session, err := getSessionToken(r)
	if err != nil {
		log.Println("getSessionToken error :", err)
		return err
	}

	follow.FollowedID, err = s.store.GetUserIdBySession(session)
	if err != nil {
		log.Println("store.GetUserIDBySession error :", err)
		return err
	}

	if follow.FollowerID == follow.FollowedID {
		log.Println("You can't follow yourself")
		return errors.New("you can't follow yourself")
	}

	err = s.store.AcceptFollowRequest(follow.FollowedID, follow.FollowerID)
	if err != nil {
		log.Println("store.AcceptFollowRequest error :", err)
		return err
	}

	return nil

}

// DELETE /follow/{id}
func (s *APIServer) HandleFollowingReject(w http.ResponseWriter, r *http.Request) error {
	var follow follow
	id := r.PathValue("id")
	follower_id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println("strconv error :", err)
		return err
	}
	follow.FollowerID = follower_id
	session, err := getSessionToken(r)
	if err != nil {
		log.Println("getSessionToken error :", err)
		return err
	}
	follow.FollowedID, err = s.store.GetUserIdBySession(session)
	if err != nil {
		log.Println("store.GetUserIDBySession error :", err)
		return err
	}
	err = s.store.DeleteFollowRequest(follow.FollowedID, follow.FollowerID)
	if err != nil {
		log.Println("store.DeleteFollowRequest error :", err)
		return err
	}
	return nil

}
