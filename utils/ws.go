package utils

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn   *websocket.Conn
	UserID string
	Rooms  map[string]bool
}

type Hub struct {
	Clients map[*Client]bool
	Rooms   map[string]map[*Client]bool
	mu      sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Clients: make(map[*Client]bool),
		Rooms:   make(map[string]map[*Client]bool),
	}
}

func (h *Hub) JoinRoom(client *Client, roomID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.Rooms[roomID] == nil {
		h.Rooms[roomID] = make(map[*Client]bool)
	}
	h.Rooms[roomID][client] = true
	client.Rooms[roomID] = true
}
func (h *Hub) BroadcastToRoom(roomID string, message WebSocketMessage) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for client := range h.Rooms[roomID] {
		client.Conn.WriteJSON(message)
	}
}
func (h *Hub) BroadcastToAll(message WebSocketMessage) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for client := range h.Clients {
		client.Conn.WriteJSON(message)
	}
}
