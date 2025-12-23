package chat_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

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
	if err := aliceConn.WriteJSON(outgoingMessage{Content: "Hello"}); err != nil {
		t.Fatalf("failed to send message from Alice: %v", err)
	}

	// Bob should receive a message event (skip system events)
	deadline := time.Now().Add(200 * time.Millisecond)
	_ = bobConn.SetReadDeadline(deadline)

	var evt chatEvent
	for {
		if err := bobConn.ReadJSON(&evt); err != nil {
			t.Fatalf("Bob did not receive message event: %v", err)
		}

		if evt.Type == "message" {
			break
		}
	}

	if evt.Sender != "Alice" {
		t.Fatalf("expected sender=Alice, got %s", evt.Sender)
	}

	if evt.Content != "Hello" {
		t.Fatalf("expected content=Hello, got %s", evt.Content)
	}

	if evt.Timestamp == 0 {
		t.Fatalf("expected non-zero timestamp")
	}
}

// Contract:
// Empty messages MUST NOT be broadcast to other clients.
func TestRejectEmptyMessage(t *testing.T) {
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

	// Alice sends empty message
	if err := aliceConn.WriteJSON(map[string]string{
		"content": "",
	}); err != nil {
		t.Fatalf("failed to send empty message: %v", err)
	}

	// Bob must NOT receive a message event
	deadline := time.Now().Add(200 * time.Millisecond)
	_ = bobConn.SetReadDeadline(deadline)

	for {
		var evt chatEvent
		err := bobConn.ReadJSON(&evt)

		if err != nil {
			// timeout or close is acceptable → no message event observed
			return
		}

		if evt.Type == "message" {
			t.Fatalf("expected no message event for empty input, but received one")
		}
		// system events are ignored
	}
}

// Contract:
// Disconnected clients MUST NOT receive messages,
// and server MUST remain stable after disconnect.
func TestDisconnectCleanup(t *testing.T) {
	server := httptest.NewServer(setupServer())
	defer server.Close()

	baseWSURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	// Connect Alice
	aliceConn, _, err := websocket.DefaultDialer.Dial(baseWSURL+"?name=Alice", nil)
	if err != nil {
		t.Fatalf("failed to connect Alice: %v", err)
	}

	// Connect Bob
	bobConn, _, err := websocket.DefaultDialer.Dial(baseWSURL+"?name=Bob", nil)
	if err != nil {
		t.Fatalf("failed to connect Bob: %v", err)
	}
	defer bobConn.Close()

	// Alice disconnects
	aliceConn.Close()

	// Give server a moment to process disconnect
	time.Sleep(50 * time.Millisecond)

	// Bob sends message
	err = bobConn.WriteJSON(map[string]string{
		"content": "Hi",
	})
	if err != nil {
		t.Fatalf("Bob failed to send message: %v", err)
	}

	// Bob SHOULD receive his own message (current contract)
	_ = bobConn.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	_, _, err = bobConn.ReadMessage()
	if err != nil {
		t.Fatalf("Bob did not receive message after Alice disconnected: %v", err)
	}

	// Alice MUST receive nothing (already closed, but ensure no panic / write)
	// If server panics or writes to closed conn, test suite will fail
}

func TestBroadcastJoinEvent(t *testing.T) {
	server := httptest.NewServer(setupServer())
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws?name=Alice"

	// Connect Alice
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("failed to connect websocket: %v", err)
	}
	defer conn.Close()

	// Read first event
	_, data, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("failed to read message: %v", err)
	}

	// Parse JSON
	var event map[string]interface{}
	if err := json.Unmarshal(data, &event); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	// Assertions
	if event["type"] != "system" {
		t.Fatalf("expected type=system, got %v", event["type"])
	}

	if event["event"] != "join" {
		t.Fatalf("expected event=join, got %v", event["event"])
	}

	if event["user"] != "Alice" {
		t.Fatalf("expected user=Alice, got %v", event["user"])
	}

	if _, ok := event["timestamp"].(float64); !ok {
		t.Fatalf("expected numeric timestamp, got %T", event["timestamp"])
	}
}

// Contract:
// When a user disconnects, remaining clients receive a leave system event.
func TestBroadcastLeaveEvent(t *testing.T) {
	server := httptest.NewServer(setupServer())
	defer server.Close()

	baseWSURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	// Connect Alice
	aliceConn, _, err := websocket.DefaultDialer.Dial(baseWSURL+"?name=Alice", nil)
	if err != nil {
		t.Fatalf("failed to connect Alice: %v", err)
	}

	// Connect Bob
	bobConn, _, err := websocket.DefaultDialer.Dial(baseWSURL+"?name=Bob", nil)
	if err != nil {
		t.Fatalf("failed to connect Bob: %v", err)
	}
	defer bobConn.Close()

	// Alice disconnects
	aliceConn.Close()

	// Bob should receive a leave system event
	deadline := time.Now().Add(200 * time.Millisecond)
	_ = bobConn.SetReadDeadline(deadline)

	for {
		var evt chatEvent
		err := bobConn.ReadJSON(&evt)
		if err != nil {
			t.Fatalf("Bob did not receive leave event: %v", err)
		}

		if evt.Type != "system" {
			continue
		}

		if evt.Event != "leave" {
			continue
		}

		if evt.User != "Alice" {
			t.Fatalf("expected user=Alice, got %s", evt.User)
		}

		if evt.Timestamp == 0 {
			t.Fatalf("expected non-zero timestamp")
		}

		// Correct leave event received
		return
	}
}

// Contract:
// All broadcast events MUST include a server-generated timestamp.
func TestTimestampIntegrity(t *testing.T) {
	server := httptest.NewServer(setupServer())
	defer server.Close()

	baseWSURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	// Connect Alice
	aliceConn, _, err := websocket.DefaultDialer.Dial(baseWSURL+"?name=Alice", nil)
	if err != nil {
		t.Fatalf("failed to connect Alice: %v", err)
	}
	defer aliceConn.Close()

	// Read first event (join)
	var evt chatEvent
	_ = aliceConn.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	if err := aliceConn.ReadJSON(&evt); err != nil {
		t.Fatalf("failed to read event: %v", err)
	}

	// timestamp must exist
	if evt.Timestamp == 0 {
		t.Fatalf("expected timestamp to be present")
	}

	// timestamp must be numeric (already enforced by int64 decode)
	// Now assert it is server-generated and sane:
	now := time.Now().UnixMilli()

	// Allow small clock skew
	if evt.Timestamp > now+1000 {
		t.Fatalf("timestamp appears to be in the future: %d > %d", evt.Timestamp, now)
	}

	if evt.Timestamp < now-10_000 {
		t.Fatalf("timestamp appears too old to be server-generated: %d", evt.Timestamp)
	}
}
