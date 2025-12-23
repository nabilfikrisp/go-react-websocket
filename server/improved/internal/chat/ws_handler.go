package chat

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type client struct {
	name string
	conn *websocket.Conn
}

var (
	clients = make(map[*websocket.Conn]*client)
	mu      sync.Mutex
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

/*
Client → Server (unchanged)
*/
type incomingMessage struct {
	Content string `json:"content"`
}

/*
Server → Client (new contract)
*/
type chatEvent struct {
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp"`

	// message
	Sender  string `json:"sender,omitempty"`
	Content string `json:"content,omitempty"`

	// system
	Event string `json:"event,omitempty"`
	User  string `json:"user,omitempty"`
}

func WSHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// Register client
	mu.Lock()
	clients[conn] = &client{name: name, conn: conn}
	mu.Unlock()

	// Broadcast join event
	joinEvent := chatEvent{
		Type:      "system",
		Event:     "join",
		User:      name,
		Timestamp: time.Now().UnixMilli(),
	}

	mu.Lock()
	for cConn, c := range clients {
		if err := c.conn.WriteJSON(joinEvent); err != nil {
			delete(clients, cConn)
			cConn.Close()
		}
	}
	mu.Unlock()

	// Ensure leave event is broadcast on disconnect
	defer func() {
		leaveEvent := chatEvent{
			Type:      "system",
			Event:     "leave",
			User:      name,
			Timestamp: time.Now().UnixMilli(),
		}

		mu.Lock()
		delete(clients, conn)

		for cConn, c := range clients {
			if err := c.conn.WriteJSON(leaveEvent); err != nil {
				delete(clients, cConn)
				cConn.Close()
			}
		}
		mu.Unlock()
	}()

	for {
		var msg incomingMessage
		if err := conn.ReadJSON(&msg); err != nil {
			return
		}

		// Ignore empty messages
		if msg.Content == "" {
			continue
		}

		event := chatEvent{
			Type:      "message",
			Sender:    name,
			Content:   msg.Content,
			Timestamp: time.Now().UnixMilli(),
		}

		// Broadcast message event
		mu.Lock()
		for cConn, c := range clients {
			if err := c.conn.WriteJSON(event); err != nil {
				delete(clients, cConn)
				cConn.Close()
			}
		}
		mu.Unlock()
	}
}
