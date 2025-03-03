package main

import (
	"fmt"
	"net/http"
)

// FOLLOWING HANDLERS
// POST /follow/{id}
func (s *APIServer) HandleFollowing(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	fmt.Println("you tried to reach id: ", id)
	return nil

}

// GET /follow/requests
func (s *APIServer) HandleFollowingRequests(w http.ResponseWriter, r *http.Request) error {
	return nil
	// fmt.Fprintln(w, "Hello there")

}

// POST /follow/{id}/accept
func (s *APIServer) HandleFollowingAccept(w http.ResponseWriter, r *http.Request) error {
	return nil

}

// DELETE /follow/{id}
func (s *APIServer) HandleFollowingReject(w http.ResponseWriter, r *http.Request) error {
	return nil

}
