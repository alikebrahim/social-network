package routes

import (
	"database/sql"
	"net/http"
)

// PROFILES HANDLERS
// GET /profiles/{id}
func ProfileGetHandler(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	
}

// PUT /profiles/privacy
func ProfileSetPrivacyHandler(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

}

// GET /profiles/{id}/activity
func ProfileGetActivitiHandler(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

}
