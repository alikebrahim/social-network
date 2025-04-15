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

	return &APIServer{
		ListenAddr: listenAddr,
		Store:      store,
		Hub:        hub,
	}
}