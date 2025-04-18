package api

import (
	"socialNetwork/pkg/storage"
	"socialNetwork/pkg/websocket"
)

type APIServer struct {
	ListenAddr string
	Store      storage.Storage
	Hub        *websocket.Hub
}

func NewAPIServer(listenAddr string, store storage.Storage) *APIServer {
	hub := websocket.NewHub()
	
	// Set the hub in the storage layer for real-time notifications
	store.SetHub(hub)

	return &APIServer{
		ListenAddr: listenAddr,
		Store:      store,
		Hub:        hub,
	}
}