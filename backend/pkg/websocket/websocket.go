package websocket

import (
	"database/sql"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	clients      = make(map[*websocket.Conn]bool)
	logedInUsers = make(map[*websocket.Conn]bool)
	broadcast    = make(chan Message)
	mutex        sync.Mutex
	DB           *sql.DB
	sender       *websocket.Conn
)

type Message struct {
	Type         string `json:"type"`
	Sender       string `json:"sender"`
	Receiver     string `json:"receiver"`
	Content      string `json:"content"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Avatar       string `json:"avatar"`
	Nickname     string `json:"nickname"`
	Bio          string `json:"bio"`
	SessionToken string `json:"session_token,omitempty"` 
	Data         User   `json:"data"`
}

type User struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BirthDate string `json:"birth_date"`
	Avatar    string `json:"avatar"`
	Nickname  string `json:"nickname"`
	SessionToken string `json:"session_token,omitempty"` 
	Bio       string `json:"bio"`
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
	logedInUsers[conn] = false
	mutex.Unlock()

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			handleWebSocketError(conn, err)
			break
		}
		broadcast <- msg
		sender = conn
	}
	log.Println("WebSocket connection closed")
}

func handleWebSocketError(conn *websocket.Conn, err error) {
	if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
		log.Printf("WebSocket closed: %v", err)
	} else {
		log.Printf("WebSocket read error: %v", err)
	}
	mutex.Lock()
	delete(clients, conn)
	mutex.Unlock()
}
