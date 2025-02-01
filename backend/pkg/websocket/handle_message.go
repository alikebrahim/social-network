package websocket

import (
	"log"
)

func HandleMessages() {
	for {
		msg := <-broadcast
		switch msg.Type {
		case "register":
			handleSignup(msg)
		case "login":
			handleLogin(msg)
			log.Printf("Login message received: %v", msg)
		case "logout":
			handleLogout(msg)
		default:
			broadcastMessageToClients(msg)
		}
	}
}

func broadcastMessageToClients(msg Message) {
	mutex.Lock()
	defer mutex.Unlock()

	for client := range clients {
		err := client.WriteJSON(msg)
		if err != nil {
			log.Printf("WebSocket write error: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}
func sendResponseToClients(response Message) {
	mutex.Lock()
	defer mutex.Unlock()

	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Printf("WebSocket write error: %v", err)
			client.Close()
			delete(clients, client)
		} else {
			sender.WriteJSON(response)
			log.Printf("Response sent to client: %v", response.Content)
		}
	}
}
