package main

import (
	"log"
	"net/http"

	"backend/pkg/db/sqlite"
	"backend/pkg/routes"
	"backend/pkg/websocket"
)

func main() {
	sqlite.InitDB()
	sqlite.RunMigrations()

	// Assign the initialized DB to the websocket package
	websocket.DB = sqlite.DB

	go websocket.HandleMessages()

	r := routes.SetupRoutes(sqlite.DB)

	r.HandleFunc("/ws", websocket.HandleConnections)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
