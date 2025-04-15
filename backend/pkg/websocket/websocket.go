package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Client represents a connected websocket client
type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
	UserID   int64
	TargetID int64 // Can be either user ID or group ID
	IsGroup  bool
}

// Hub manages WebSocket connections
type Hub struct {
	clients    map[*Client]bool
	Register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	mutex      sync.Mutex
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Run starts the WebSocket hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mutex.Lock()
			h.clients[client] = true
			h.mutex.Unlock()
		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			h.mutex.Unlock()
		case message := <-h.broadcast:
			h.mutex.Lock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
			h.mutex.Unlock()
		}
	}
}

// SendMessageToUser sends a message to a specific user
func (h *Hub) SendMessageToUser(targetUserID int64, message []byte) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for client := range h.clients {
		if !client.IsGroup && client.UserID == targetUserID {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(h.clients, client)
			}
		}
	}
}

// SendMessageToGroup sends a message to all users in a group
func (h *Hub) SendMessageToGroup(groupID int64, message []byte) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for client := range h.clients {
		if client.IsGroup && client.TargetID == groupID {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(h.clients, client)
			}
		}
	}
}

// ReadPump handles reading messages from the WebSocket
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Register <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(1024 * 1024) // 1MB
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Process the message
		var messageObj map[string]interface{}
		if err := json.Unmarshal(message, &messageObj); err != nil {
			log.Printf("error unmarshaling message: %v", err)
			continue
		}

		// Add timestamp if not present
		if _, ok := messageObj["timestamp"]; !ok {
			messageObj["timestamp"] = time.Now()
			message, _ = json.Marshal(messageObj)
		}

		// Send to appropriate recipients based on the message type
		if c.IsGroup {
			c.Hub.SendMessageToGroup(c.TargetID, message)
		} else {
			// Regular chat - send to both the sender and receiver
			c.Hub.SendMessageToUser(c.UserID, message)      // Echo to sender
			c.Hub.SendMessageToUser(c.TargetID, message)    // Send to recipient
		}
	}
}

// WritePump handles sending messages to the WebSocket
func (c *Client) WritePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// The hub closed the channel
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}