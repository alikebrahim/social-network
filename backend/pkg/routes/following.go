package routes

import (
	"database/sql"
	"fmt"
	"net/http"
)

// FOLLOWING HANDLERS
// POST /follow/{id}
func FollowingHandler(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
}

// GET /follow/requests
func FollowingRequestsHandler(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	fmt.Fprintln(w, "Hello there")

}

// POST /follow/{id}/accept
func FollowAcceptRequestHandler(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

}

// DELETE /follow/{id}
func FollowRejectRequestHandler(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

}
