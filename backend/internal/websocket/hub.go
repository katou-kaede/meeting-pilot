package websocket

import (
	"sync"

	gorilla "github.com/gorilla/websocket"
)

type Event struct {
	Type      string `json:"type"`
	MeetingID int64  `json:"meeting_id"`
}

type Hub struct {
	mu      sync.RWMutex
	clients map[int64]map[*gorilla.Conn]bool
}

func NewHub() *Hub {
	return &Hub{
		clients: make(
			map[int64]map[*gorilla.Conn]bool,
		),
	}
}

func (h *Hub) AddClient(meetingID int64, conn *gorilla.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[meetingID] == nil {
		h.clients[meetingID] = make(map[*gorilla.Conn]bool)
	}

	h.clients[meetingID][conn] = true
}

func (h *Hub) RemoveClient(meetingID int64, conn *gorilla.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	clients := h.clients[meetingID]
	if clients == nil {
		return
	}

	delete(clients, conn)

	if len(clients) == 0 {
		delete(h.clients, meetingID)
	}

	conn.Close()
}

func (h *Hub) Broadcast(meetingID int64, message any) {
	h.mu.RLock()

	clients := make(
		[]*gorilla.Conn,
		0,
		len(h.clients[meetingID]),
	)

	for conn := range h.clients[meetingID] {
		clients = append(clients, conn)
	}

	h.mu.RUnlock()

	for _, conn := range clients {
		if err := conn.WriteJSON(message); err != nil {
			h.RemoveClient(meetingID, conn)
		}
	}
}
