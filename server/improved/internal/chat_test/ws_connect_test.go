package chat_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

type incomingMessage struct {
	Sender  string `json:"sender"`
	Content string `json:"content"`
}

type outgoingMessage struct {
	Content string `json:"content"`
}

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

// Contract:
// A WebSocket connection WITH a name MUST be accepted.
func TestAcceptNamedConnection(t *testing.T) {
	// Arrange
	server := httptest.NewServer(setupServer())
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws?name=Alice"

	// Act
	conn, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("expected websocket connection to succeed, got error: %v", err)
	}
	defer conn.Close()

	if resp.StatusCode != 101 {
		t.Fatalf("expected status 101 Switching Protocols, got %d", resp.StatusCode)
	}

	// Assert connection is not immediately closed
	_ = conn.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
	_, _, err = conn.ReadMessage()

	if err == nil {
		// No message is fine, but connection must still be alive
		return
	}

	if websocket.IsUnexpectedCloseError(err) {
		t.Fatalf("connection was closed unexpectedly after successful upgrade")
	}
}

// Contract:
// When Alice sends a message, Bob receives it with Alice as sender.
func TestBroadcastMessage(t *testing.T) {
	server := httptest.NewServer(setupServer())
	defer server.Close()

	baseWSURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	// Connect Alice
	aliceConn, _, err := websocket.DefaultDialer.Dial(baseWSURL+"?name=Alice", nil)
	if err != nil {
		t.Fatalf("failed to connect Alice: %v", err)
	}
	defer aliceConn.Close()

	// Connect Bob
	bobConn, _, err := websocket.DefaultDialer.Dial(baseWSURL+"?name=Bob", nil)
	if err != nil {
		t.Fatalf("failed to connect Bob: %v", err)
	}
	defer bobConn.Close()

	// Alice sends message
	err = aliceConn.WriteJSON(outgoingMessage{Content: "Hello"})
	if err != nil {
		t.Fatalf("failed to send message from Alice: %v", err)
	}

	// Bob should receive it
	_ = bobConn.SetReadDeadline(time.Now().Add(200 * time.Millisecond))

	var msg incomingMessage
	err = bobConn.ReadJSON(&msg)
	if err != nil {
		t.Fatalf("Bob did not receive message: %v", err)
	}

	if msg.Sender != "Alice" {
		t.Fatalf("expected sender 'Alice', got '%s'", msg.Sender)
	}

	if msg.Content != "Hello" {
		t.Fatalf("expected content 'Hello', got '%s'", msg.Content)
	}
}
