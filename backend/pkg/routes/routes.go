package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"backend/pkg/models"
)

func SetupRoutes(db *sql.DB) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		var user struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Bio      string `json:"bio"`
		}
		json.NewDecoder(r.Body).Decode(&user)

		_, err := db.Exec("INSERT INTO users (username, password, bio) VALUES (?, ?, ?)", user.Username, user.Password, user.Bio)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}).Methods("POST")

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		var credentials struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		json.NewDecoder(r.Body).Decode(&credentials)

		var id int
		err := db.QueryRow("SELECT id FROM users WHERE username = ? AND password = ?", credentials.Username, credentials.Password).Scan(&id)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"user_id": id,
		})
	}).Methods("POST")

	r.HandleFunc("/add_post", func(w http.ResponseWriter, r *http.Request) {
		var post struct {
			UserID  int    `json:"user_id"`
			Content string `json:"content"`
		}
		json.NewDecoder(r.Body).Decode(&post)

		_, err := db.Exec("INSERT INTO posts (user_id, content) VALUES (?, ?)", post.UserID, post.Content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}).Methods("POST")

	r.HandleFunc("/view_profile/{username}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]

		var user models.User
		err := db.QueryRow("SELECT id, username, bio FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Bio)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(user)
	}).Methods("GET")

	return r
}