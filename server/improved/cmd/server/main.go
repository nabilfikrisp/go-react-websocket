package main

import (
	"log"
	"net/http"
	"os"

	"go-websocket-improved/internal/chat"
)

func main() {
	mux := http.NewServeMux()

	// Static client (SPA)
	fs := http.FileServer(http.Dir("static"))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Serve index.html for SPA routes
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "static/index.html")
			return
		}

		// Try static file
		if _, err := os.Stat("static" + r.URL.Path); err == nil {
			fs.ServeHTTP(w, r)
			return
		}

		// Fallback to SPA entry
		http.ServeFile(w, r, "static/index.html")
	})

	// WebSocket endpoint
	mux.HandleFunc("/ws", chat.WSHandler)

	addr := ":3001"
	log.Printf("server listening on %s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
