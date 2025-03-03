package routes

import (
	"database/sql"
	// "errors"
	// "log"
	"net/http"
	// "strconv"
	// "strings"
	"github.com/gorilla/mux"
)

func SetupRoutes(database *sql.DB) *mux.Router  {
	DB = database

	r := mux.NewRouter()

	r.HandleFunc("/auth/register", func(w http.ResponseWriter, r *http.Request) {
		RegisterHandler(w, r, DB)
	}).Methods("POST")
	r.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		LoginHandler(w, r, DB)
	}).Methods("POST")
	r.HandleFunc("/auth/logout", func(w http.ResponseWriter, r *http.Request) {
		LogoutHandler(w, r, DB)
	}).Methods("POST")


	r.HandleFunc("/follow/{id}", func(w http.ResponseWriter, r *http.Request) {
		FollowingHandler(w, r, DB)
	}).Methods("POST")
	r.HandleFunc("/follow/requests", func(w http.ResponseWriter, r *http.Request) {
		FollowingRequestsHandler(w, r, DB)
	}).Methods("GET")
	r.HandleFunc("/follow/accept/{id}", func(w http.ResponseWriter, r *http.Request) {
		FollowAcceptRequestHandler(w, r, DB)
	}).Methods("POST")
	r.HandleFunc("/follow/reject/{id}", func(w http.ResponseWriter, r *http.Request) {
		FollowRejectRequestHandler(w, r, DB)
	}).Methods("DELETE")

	// ✅ GROUPS HANDLERS
	r.HandleFunc("/groups", func(w http.ResponseWriter, r *http.Request) {
		GroupCreateHandler(w, r, DB)
	}).Methods("POST")
	r.HandleFunc("/groups/{id}/invite", func(w http.ResponseWriter, r *http.Request) {
		GroupInviteHandler(w, r, DB)
	}).Methods("POST")
	r.HandleFunc("/groups/search", func(w http.ResponseWriter, r *http.Request) {
		GroupSearchHandler(w, r, DB)
	}).Methods("GET")
	r.HandleFunc("/groups/{id}/join", func(w http.ResponseWriter, r *http.Request) {
		GroupRequestHandler(w, r, DB)
	}).Methods("POST")
	r.HandleFunc("/groups/{id}/requests", func(w http.ResponseWriter, r *http.Request) {
		GroupListRequestsHandler(w, r, DB)
	}).Methods("GET")
	r.HandleFunc("/groups/{id}/requests/{userId}/accept", func(w http.ResponseWriter, r *http.Request) {
		GroupAcceptRequestHandler(w, r, DB)
	}).Methods("POST")
	r.HandleFunc("/groups/{id}/events", func(w http.ResponseWriter, r *http.Request) {
		GroupCreateEventHandler(w, r, DB)
	}).Methods("POST")

	// ✅ POSTS HANDLERS
	r.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		PostCreateHandler(w, r, DB)
	}).Methods("POST")
	r.HandleFunc("/posts/{id}", func(w http.ResponseWriter, r *http.Request) {
		// is't suposed the id geting it from the cookies why the id in the url?
		PostsGetHandler(w, r, DB)
	}).Methods("GET")
	r.HandleFunc("/posts/{id}", func(w http.ResponseWriter, r *http.Request) {
		PostEditHandler(w, r, DB)
	}).Methods("PUT")
	r.HandleFunc("/posts/{id}", func(w http.ResponseWriter, r *http.Request) {
		PostDeleteHandler(w, r, DB)
	}).Methods("DELETE")
	r.HandleFunc("/posts/{id}/comments", func(w http.ResponseWriter, r *http.Request) {
		PostCommentHandler(w, r, DB)
	}).Methods("POST")
	r.HandleFunc("/posts/{id}/likes", func(w http.ResponseWriter, r *http.Request) {
		PostLikeHandler(w, r, DB)
	}).Methods("POST")

	// ✅ PROFILES HANDLERS
	r.HandleFunc("/profiles/{id}", func(w http.ResponseWriter, r *http.Request) {
		ProfileGetHandler(w, r, DB)
	}).Methods("GET")
	r.HandleFunc("/profiles/privacy", func(w http.ResponseWriter, r *http.Request) {
		ProfileSetPrivacyHandler(w, r, DB)
	}).Methods("PUT")
	r.HandleFunc("/profiles/{id}/activity", func(w http.ResponseWriter, r *http.Request) {
		ProfileGetActivitiHandler(w, r, DB)
	}).Methods("GET")
	return r
}

// func getIdFromRoute(w http.ResponseWriter, r *http.Request, key string) (int64, error) {
// 	pathParts := strings.Split(r.URL.Path, "/")

// 	for i, part := range pathParts {
// 		if part == key && i+1 < len(pathParts) {
// 			idStr := pathParts[i+1]
// 			id, err := strconv.ParseInt(idStr, 10, 64)
// 			if err != nil {
// 				http.Error(w, "Invalid URL", http.StatusBadRequest)
// 				log.Print(err)
// 				return 0, err
// 			}
// 			return id, nil
// 		}
// 	}

// 	http.Error(w, "Invalid URL", http.StatusBadRequest)
// 	log.Print("Key not found in URL")
// 	return 0, errors.New("Invalid URL")
// }
