package chat_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

// Contract:
// A WebSocket connection WITHOUT a name MUST be rejected.
func TestRejectAnonymousConnection(t *testing.T) {
	// Arrange: start test server
	server := httptest.NewServer(setupServer()) // to be implemented
	defer server.Close()

	// Convert http://127.0.0.1 -> ws://127.0.0.1/ws
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	// Act: attempt anonymous WebSocket connection
	_, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)

	// Assert: connection must not succeed
	if err == nil {
		t.Fatalf("expected anonymous websocket connection to be rejected, but it succeeded")
	}

	// If handshake failed before HTTP response, that's acceptable
	if resp == nil {
		return
	}

	// Explicitly forbid successful upgrade
	if resp.StatusCode == http.StatusSwitchingProtocols {
		t.Fatalf("expected websocket upgrade to be rejected, got status 101")
	}
}
