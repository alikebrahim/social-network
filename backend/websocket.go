package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Consider restricting this in production
	},
}

// WebSocketClient represents a connected websocket client
type WebSocketClient struct {
	UserID   int64
	Conn     *websocket.Conn
	Send     chan []byte
	ClientID string // For testing to track client identity
}

// WebSocketMessage represents a message that can be sent over a WebSocket
type WebSocketMessage struct {
	Type      string      `json:"type"`
	Payload   interface{} `json:"payload"`
	Client_ID string      `json:"client_id,omitempty"` // Used for test client to track message origin
}

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	// Registered clients
	clients map[int64]*WebSocketClient

	// Register requests from clients
	register chan *WebSocketClient

	// Unregister requests from clients
	unregister chan *WebSocketClient

	// Inbound messages to broadcast
	broadcast chan []byte

	// Mutex for thread safety
	mutex sync.Mutex

	// Group chat tracking - maps group ID to set of user IDs
	groupClients map[int64]map[int64]*WebSocketClient
	
	// Client mapping by client ID (for test client)
	clientIDMap map[string]*WebSocketClient
}

// Initialize a new hub instance
var hub = Hub{
	broadcast:    make(chan []byte),
	register:     make(chan *WebSocketClient),
	unregister:   make(chan *WebSocketClient),
	clients:      make(map[int64]*WebSocketClient),
	groupClients: make(map[int64]map[int64]*WebSocketClient),
	clientIDMap:  make(map[string]*WebSocketClient),
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client.UserID] = client
			
			// Add to client ID mapping if available
			if client.ClientID != "" {
				h.clientIDMap[client.ClientID] = client
				log.Printf("Registered client with ID: %s (total clients: %d)", 
					client.ClientID, len(h.clientIDMap))
			}
			
			h.mutex.Unlock()

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
				
				// Remove from client ID mapping
				if client.ClientID != "" {
					delete(h.clientIDMap, client.ClientID)
					log.Printf("Unregistered client with ID: %s", client.ClientID)
				}
				
				// Remove from group chats
				for groupID, users := range h.groupClients {
					if _, ok := users[client.UserID]; ok {
						delete(h.groupClients[groupID], client.UserID)
					}
				}
			}
			h.mutex.Unlock()

		case message := <-h.broadcast:
			// Process the message to determine recipients
			var wsMsg WebSocketMessage
			if err := json.Unmarshal(message, &wsMsg); err != nil {
				log.Printf("Error unmarshaling broadcast message: %v", err)
				continue
			}

			// Extract chat message
			var chatMsg ChatMessage
			payloadBytes, err := json.Marshal(wsMsg.Payload)
			if err != nil {
				log.Printf("Error re-marshaling payload: %v", err)
				continue
			}
			
			if err := json.Unmarshal(payloadBytes, &chatMsg); err != nil {
				log.Printf("Error unmarshaling chat message: %v", err)
				continue
			}

			h.mutex.Lock()

			// For test client, we need special handling
			if wsMsg.Client_ID != "" {
				log.Printf("Broadcasting message from client %s", wsMsg.Client_ID)
				
				// Get all registered client IDs for debugging
				var clientIDs []string
				for id := range h.clientIDMap {
					clientIDs = append(clientIDs, id)
				}
				log.Printf("All registered client IDs: %v", clientIDs)
				
				// Find the sender and recipient
				senderClient, senderExists := h.clientIDMap[wsMsg.Client_ID]
				
				// Send to all clients EXCEPT the sender
				for id, client := range h.clientIDMap {
					if id != wsMsg.Client_ID && client != nil && client.Conn != nil {
						log.Printf("Attempting to send message to client %s", id)
						
						select {
						case client.Send <- message:
							log.Printf("Successfully sent message to client %s", id)
						default:
							log.Printf("Failed to send to client %s (buffer full)", id)
							close(client.Send)
							delete(h.clients, client.UserID)
							delete(h.clientIDMap, id)
						}
					}
				}
				
				// Always echo back to sender for UI update
				if senderExists && senderClient != nil {
					select {
					case senderClient.Send <- message:
						log.Printf("Echoed message back to sender %s", wsMsg.Client_ID)
					default:
						log.Printf("Failed to echo to sender %s (buffer full)", wsMsg.Client_ID)
						close(senderClient.Send)
						delete(h.clients, senderClient.UserID)
						delete(h.clientIDMap, wsMsg.Client_ID)
					}
				}
			} else {
				// Normal operation (production code path)
				if chatMsg.Type == "private" {
					// Send to specific recipient
					if client, ok := h.clients[chatMsg.ReceiverID]; ok {
						select {
						case client.Send <- message:
						default:
							close(client.Send)
							delete(h.clients, client.UserID)
						}
					}
					
					// Also send to sender for their own display
					if client, ok := h.clients[chatMsg.SenderID]; ok {
						select {
						case client.Send <- message:
						default:
							close(client.Send)
							delete(h.clients, client.UserID)
						}
					}
				} else if chatMsg.Type == "group" {
					// Send to all members of the group
					if groupMembers, ok := h.groupClients[chatMsg.GroupID]; ok {
						for _, client := range groupMembers {
							select {
							case client.Send <- message:
							default:
								close(client.Send)
								delete(h.clients, client.UserID)
								delete(groupMembers, client.UserID)
							}
						}
					}
				}
			}
			
			h.mutex.Unlock()
		}
	}
}

// Register a client with a group chat
func (h *Hub) RegisterGroupClient(groupID int64, client *WebSocketClient) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	if _, ok := h.groupClients[groupID]; !ok {
		h.groupClients[groupID] = make(map[int64]*WebSocketClient)
	}
	
	h.groupClients[groupID][client.UserID] = client
}

// Create a WebSocket message
func createMessage(messageType string, payload interface{}) []byte {
	message := WebSocketMessage{
		Type:    messageType,
		Payload: payload,
	}
	
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling WebSocket message: %v", err)
		return nil
	}
	
	return data
}

// Start the client's goroutines for writing and reading messages
func (c *WebSocketClient) Start() {
	// Start the write pump in a goroutine
	go c.writePump()
	
	// Start the read pump in the current goroutine
	c.readPump()
}

// This method is used when we want to handle reading messages elsewhere
func (c *WebSocketClient) StartWritePumpOnly() {
	c.writePump()
}

// writePump pumps messages from the hub to the websocket connection
func (c *WebSocketClient) writePump() {
	ticker := time.NewTicker(60 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// Hub closed the channel
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			
			// Write a single message at a time instead of batching
			// This ensures each message is a complete JSON object
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("Error writing to client %s: %v", c.ClientID, err)
				return
			}
			
			// Process any queued messages individually
			n := len(c.Send)
			for i := 0; i < n; i++ {
				msg := <-c.Send
				if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
					log.Printf("Error writing queued message to client %s: %v", c.ClientID, err)
					return
				}
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// readPump pumps messages from the websocket connection to the hub
func (c *WebSocketClient) readPump() {
	defer func() {
		hub.unregister <- c
		c.Conn.Close()
	}()
	
	c.Conn.SetReadLimit(512 * 1024) // 512KB max message size
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error { 
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil 
	})
	
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error for client %s: %v", c.ClientID, err)
			}
			break
		}
		
		// Try to preserve client ID for message routing
		var wsMsg WebSocketMessage
		if err := json.Unmarshal(message, &wsMsg); err == nil && wsMsg.Client_ID == "" {
			// Add client ID if missing
			if c.ClientID != "" {
				wsMsg.Client_ID = c.ClientID
				if newMessage, err := json.Marshal(wsMsg); err == nil {
					message = newMessage
				}
			}
		}
		
		// Process the received message
		hub.broadcast <- message
	}
}