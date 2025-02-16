package main

import (
	"log"
	"net/http"

	db "backend/pkg/db/sqlite"
	routes "backend/pkg/routes"
	"backend/pkg/websocket"
)

var (
	port = ":8080"
)

func init() {
	db.InitDB()
	db.RunMigrations()
}

func main() {
	// Assign the initialized DB to the websocket package
	websocket.DB = db.DB

	go websocket.HandleMessages()

	mux := routes.SetupRoutes(db.DB)

	mux.HandleFunc("/ws", websocket.HandleConnections)

	log.Println("Server started on port", port)
	log.Fatal(http.ListenAndServe(port, mux))
}
