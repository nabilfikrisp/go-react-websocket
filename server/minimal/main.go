package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

type M map[string]any

type SocketPayload struct {
	Message string
}

type SocketResponse struct {
	From    string
	Type    string
	Message string
}

type WebSocketConnection struct {
	*websocket.Conn
	Username string
}

const MESSAGE_NEW_USER = "New User"
const MESSAGE_CHAT = "Chat"
const MESSAGE_LEAVE = "Leave"

var connections = make([]*WebSocketConnection, 0)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		content, err := os.ReadFile("index.html")
		if err != nil {
			http.Error(w, "Could not open the requested file", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "%s", content)
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsUpgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}

		currentGorillaConn, err := wsUpgrader.Upgrade(w, r, w.Header())

		if err != nil {
			http.Error(w, "Could not open web socket connection", http.StatusBadRequest)
		}

		username := r.URL.Query().Get("username")
		currentConn := WebSocketConnection{Conn: currentGorillaConn, Username: username}
		connections = append(connections, &currentConn)

		go handleIO(&currentConn, connections)
	})

	fmt.Println("Server starting at :3000")
	http.ListenAndServe(":3000", nil)
}

func handleIO(currentConn *WebSocketConnection, connections []*WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR", fmt.Sprintf("%v", r))
		}
	}()

	broadcastMessage(currentConn, MESSAGE_NEW_USER, "")

	for {
		payload := SocketPayload{}
		err := currentConn.ReadJSON(&payload)

		if err != nil {
			if strings.Contains(err.Error(), "websocket: close") {
				broadcastMessage(currentConn, MESSAGE_LEAVE, "")
				ejectConnection(currentConn)
				return
			}

			log.Println("ERROR", err.Error())
			continue
		}

		broadcastMessage(currentConn, MESSAGE_CHAT, payload.Message)
	}
}

func ejectConnection(currentConn *WebSocketConnection) {
	filtered := make([]*WebSocketConnection, 0, len(connections))

	for _, conn := range connections {
		if conn != currentConn {
			filtered = append(filtered, conn)
		}
	}

	connections = filtered
}

func broadcastMessage(currentConn *WebSocketConnection, kind, message string) {
	for _, eachConn := range connections {
		if eachConn == currentConn {
			continue
		}

		eachConn.WriteJSON(SocketResponse{
			From:    currentConn.Username,
			Type:    kind,
			Message: message,
		})
	}
}
