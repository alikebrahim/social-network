package routes

import (
	"net/http"
)

// func SetupRoutes(db *sql.DB) *mux.Router {
// 	r := mux.NewRouter()
//
// 	r.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
// 		var user struct {
// 			Username string `json:"username"`
// 			Password string `json:"password"`
// 			Bio      string `json:"bio"`
// 		}
// 		json.NewDecoder(r.Body).Decode(&user)
//
// 		_, err := db.Exec("INSERT INTO users (username, password, bio) VALUES (?, ?, ?)", user.Username, user.Password, user.Bio)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		w.WriteHeader(http.StatusCreated)
// 	}).Methods("POST")
//
// 	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
// 		var credentials struct {
// 			Username string `json:"username"`
// 			Password string `json:"password"`
// 		}
// 		json.NewDecoder(r.Body).Decode(&credentials)
//
// 		var id int
// 		err := db.QueryRow("SELECT id FROM users WHERE username = ? AND password = ?", credentials.Username, credentials.Password).Scan(&id)
// 		if err != nil {
// 			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
// 			return
// 		}
//
// 		json.NewEncoder(w).Encode(map[string]interface{}{
// 			"user_id": id,
// 		})
// 	}).Methods("POST")
//
// 	r.HandleFunc("/add_post", func(w http.ResponseWriter, r *http.Request) {
// 		var post struct {
// 			UserID  int    `json:"user_id"`
// 			Content string `json:"content"`
// 		}
// 		json.NewDecoder(r.Body).Decode(&post)
//
// 		_, err := db.Exec("INSERT INTO posts (user_id, content) VALUES (?, ?)", post.UserID, post.Content)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		w.WriteHeader(http.StatusCreated)
// 	}).Methods("POST")
//
// 	r.HandleFunc("/view_profile/{username}", func(w http.ResponseWriter, r *http.Request) {
// 		vars := mux.Vars(r)
// 		username := vars["username"]
//
// 		var user models.User
// 		err := db.QueryRow("SELECT id, username, bio FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Bio)
// 		if err != nil {
// 			http.Error(w, "User not found", http.StatusNotFound)
// 			return
// 		}
//
// 		json.NewEncoder(w).Encode(user)
// 	}).Methods("GET")
//
// 	return r
// }

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// AUTH HANDLERS
	// POST /auth/register
	mux.HandleFunc("POST /auth/register", RegisterHandler)
	// POST /auth/login
	mux.HandleFunc("POST /auth/login", RegisterHandler)
	// POST /auth/logout
	mux.HandleFunc("POST /auth/logout", RegisterHandler)

	// FOLLOWING HANDLERS
	// POST /follow/{id}
	mux.HandleFunc("POST /follow/{id}", FollowingHandler)
	// GET /follow/requests
	mux.HandleFunc("GET /follow/requests", FollowingRequestsHandler)
	// POST /follow/{id}/accept
	mux.HandleFunc("POST /follow/{id}/accept", FollowAcceptRequestHandler)
	// DELETE /follow/{id}
	mux.HandleFunc("DELETE /follow/{id}", FollowRejectRequestHandler)

	// GROUPS HANDLERS
	// POST /groups
	mux.HandleFunc("POST /groups", GroupCreateHandler)
	// POST /groups/{id}/invite
	mux.HandleFunc("POST /groups/{id}/invite", GroupInviteHandler)
	// GET /groups/search
	mux.HandleFunc("GET /groups/search", GroupSearchHandler)
	// POST /groups/{id}/join
	mux.HandleFunc("POST /groups/{id}/join", GroupRequestHandler)
	// GET /groups/{id}/requests
	mux.HandleFunc("GET /groups/{id}/requests", GroupListRequestsHandler)
	// POST /groups/{id}/requests/{userId}/accept
	mux.HandleFunc("POST /groups/{id}/requests/{userId}/accept", GroupAcceptRequestHandler)
	// POST /groups/{id}/events
	mux.HandleFunc("POST /groups/{id}/events", GroupCreateEventHandler)

	// POSTS HANDLERS
	// POST /posts
	mux.HandleFunc("POST /posts", PostCreateHandler)
	// GET /posts/{id}
	mux.HandleFunc("GET /posts/{id}", PostsGetHandler)
	// PUT /posts/{id}
	mux.HandleFunc("PUT /posts/{id}", PostEditHandler)
	// DELETE /posts/{id}
	mux.HandleFunc("DELETE /posts/{id}", PostDeleteHandler)
	// POST /posts/{id}/comments
	mux.HandleFunc("POST /posts/{id}/comments", PostCommentHandler)
	// POST /posts/{id}/likes
	mux.HandleFunc("POST /posts/{id}/likes", PostLikeHandler)

	// PROFILES HANDLERS
	// GET /profiles/{id}
	mux.HandleFunc("GET /profiles/{id}", ProfileGetHandler)
	// PUT /profiles/privacy
	mux.HandleFunc("PUT /profiles/privacy", ProfileSetPrivacyHandler)
	// GET /profiles/{id}/activity
	mux.HandleFunc("GET /profiles/{id}/activity", ProfileGetActivitiHandler)
	return mux
}
