package chat_test

import (
	"go-websocket-improved/internal/chat"
	"net/http"
)

func setupServer() http.Handler {
	return chat.NewServer()
}
