package chat

import (
	"net/http"
)

// NewServer returns an http.Handler with routes wired.
// This is the only thing tests are allowed to depend on.
func NewServer() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", wsHandler)
	return mux
}
