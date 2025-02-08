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
			log.Printf("Signup message received: %v", msg)
		case "login":
			handleLogin(msg)
			log.Printf("Login message received: %v", msg)
		case "logout":
			handleLogout(msg)
			log.Printf("Logout message received: %v", msg)
		case "create_post":
			handlePostCreation(msg)
			log.Printf("Post message received: %v", msg)
		case "get_posts":
			handleGetPosts(msg)
			log.Printf("Get posts message received: %v", msg)
		case "follow_request":
			handleFollowRequest(msg)
			log.Printf("Follow request message received: %v", msg)
		case "follow_response":
			handleFollowResponse(msg)
			log.Printf("Follow response message received: %v", msg)
		case "unfollow":
			handleUnfollow(msg)
			log.Printf("Unfollow message received: %v", msg)
		case "get_followers":
			handleGetFollowers(msg)
			log.Printf("Get followers message received: %v", msg)
			
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
		}
	}
	log.Printf("Response sent to clients: %v", response.Content)
}

