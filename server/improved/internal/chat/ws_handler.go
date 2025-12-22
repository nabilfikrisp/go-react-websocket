package chat

import (
	"net/http"
	"sync"

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

type incomingMessage struct {
	Content string `json:"content"`
}

type outgoingMessage struct {
	Sender  string `json:"sender"`
	Content string `json:"content"`
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
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

	// Cleanup on disconnect
	defer func() {
		mu.Lock()
		delete(clients, conn)
		mu.Unlock()
	}()

	for {
		var msg incomingMessage
		if err := conn.ReadJSON(&msg); err != nil {
			return
		}

		// Ignore empty messages (will be enforced by Test 4)
		if msg.Content == "" {
			continue
		}

		out := outgoingMessage{
			Sender:  name,
			Content: msg.Content,
		}

		// Broadcast to all clients
		mu.Lock()
		for cConn, c := range clients {
			if err := c.conn.WriteJSON(out); err != nil {
				// Remove dead client immediately
				delete(clients, cConn)
				cConn.Close()
			}
		}
		mu.Unlock()
	}
}
