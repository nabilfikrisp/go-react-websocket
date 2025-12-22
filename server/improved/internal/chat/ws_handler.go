package chat

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // irrelevant for now
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract name from query
	name := r.URL.Query().Get("name")

	// Reject anonymous connections
	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	// Upgrade only if name exists
	_, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// IMPORTANT:
	// Do nothing else yet.
	// Next tests will force structure to appear.
}
