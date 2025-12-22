package main

import (
	"log"
	"net/http"

	"go-websocket-improved/internal/chat"
)

func main() {
	mux := http.NewServeMux()

	// Base endpoint
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","message":"hello"}`))
	})

	// WebSocket endpoint
	mux.HandleFunc("/ws", chat.WSHandler)

	addr := ":3001"
	log.Printf("server listening on %s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
