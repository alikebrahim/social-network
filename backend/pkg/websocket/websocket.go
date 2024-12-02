package websocket

import (
	"database/sql"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan Message)
	mutex     sync.Mutex
	DB        *sql.DB
)

type Message struct {
	Type     string `json:"type"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Content  string `json:"content"`
	Username string `json:"username"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
	Data     User   `json:"data"`
}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	log.Println("WebSocket connection established")

	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		println(msg.Type)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("WebSocket closed: %v", err)
			} else {
				log.Printf("WebSocket read error: %v", err)
			}
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			break
		}
		broadcast <- msg
	}
	log.Println("WebSocket connection closed")
}

func HandleMessages() {
	for {
		msg := <-broadcast
		switch msg.Type {
		case "register":
			handleSignup(msg)
		default:
			mutex.Lock()
			for client := range clients {
				err := client.WriteJSON(msg)
				if err != nil {
					log.Printf("WebSocket write error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
			mutex.Unlock()
		}
	}
}

func handleSignup(msg Message) {
	_, err := DB.Exec("INSERT INTO users (username, password) VALUES (?, ?)", msg.Data.Email, msg.Data.Password)
	if err != nil {
		log.Printf("Failed to sign up user: %v", err)
		return
	}
	log.Printf("User signed up: %s", msg.Username)
}
