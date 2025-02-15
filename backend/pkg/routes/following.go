package routes

import (
	"fmt"
	"net/http"
)

// FOLLOWING HANDLERS
// POST /follow/{id}
func FollowingHandler(w http.ResponseWriter, r *http.Request) {

}

// GET /follow/requests
func FollowingRequestsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello there")

}

// POST /follow/{id}/accept
func FollowAcceptRequestHandler(w http.ResponseWriter, r *http.Request) {

}

// DELETE /follow/{id}
func FollowRejectRequestHandler(w http.ResponseWriter, r *http.Request) {

}
