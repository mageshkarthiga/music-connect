package chat

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"cloud.google.com/go/firestore"
	"strings"
	"context"
	"github.com/gorilla/websocket"
)

// for readPump goroutine
const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 10000
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Upgrader used to upgrade the HTTP connection to a WebSocket connection.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client represents a WebSocket client.
type Client struct {
	conn     *websocket.Conn
	wsServer *WsServer
	send     chan []byte
	rooms    map[*Room]bool
	UserID   string `json:"user_id"`
}

func newClient(conn *websocket.Conn, wsServer *WsServer, userID string) *Client {
	return &Client{
		conn:     conn,
		wsServer: wsServer,
		send:     make(chan []byte, 256),
		rooms:    make(map[*Room]bool),
		UserID:   userID,
	}
}

func (client *Client) disconnect() {
	client.wsServer.unregister <- client
	for room := range client.rooms {
		room.unregister <- client
	}
	client.conn.Close()
}

// Read new messages from the WebSocket connection and broadcast them to all clients.
func (client *Client) readPump() {
	defer func() {
		client.disconnect()
	}()
	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error {
		client.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, jsonMessage, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println("Error while reading message:", err)
			}
			break
		}
		client.handleNewMessage(jsonMessage)
	}
}

// Keeps connection alive by sending ping messages to the client.
func (client *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Write the message to the WebSocket connection.
			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-client.send)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (client *Client) handleNewMessage(jsonMessage []byte) {
	var message Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error while unmarshalling message: %v", err)
		return
	}

	message.Sender = client.UserID
	switch message.Action {
	case SendMessageAction:
		roomName := message.Target
		if room := client.wsServer.findRoom(roomName); room != nil {
			room.broadcast <- &message
		}
	case JoinRoomAction:
		log.Printf("Processing join-room action for room: %s", message.Message)
		client.handleJoinRoom(message)
	case LeaveRoomAction:
		client.handleLeaveRoom(message)
	}
}

func (client *Client) handleJoinRoom(message Message) {
	roomName := message.Message

	room := client.wsServer.findRoom(roomName)
	if room == nil {
		fmt.Printf("Creating new room: %s\n", roomName)
		room = client.wsServer.createRoom(roomName)

		// Save room to Firestore with participants
		participants := parseParticipantsFromRoomName(roomName) // new helper
		_, err := FirestoreClient.Collection("rooms").Doc(roomName).Set(context.Background(), map[string]interface{}{
			"participants": participants,
			"created_at":   firestore.ServerTimestamp,
		}, firestore.MergeAll)
		if err != nil {
			log.Printf("Failed to store room in Firestore: %v\n", err)
		} else {
			log.Printf("Room %s saved to Firestore\n", roomName)
		}
	} else {
		fmt.Printf("Room %s already exists\n", roomName)
	}

	client.rooms[room] = true

	room.register <- client
}

func parseParticipantsFromRoomName(name string) []string {
	// roomName is like "userA-userB"
	log.Printf("Room name before parsing: %s", name)
	parts := strings.Split(name, "-")
	if len(parts) != 2 {
		log.Printf("Unexpected room format: %s", name)
		return []string{}
	}
	return parts
}


func (client *Client) handleLeaveRoom(message Message) {
	room := client.wsServer.findRoom(message.Message)
	if _, ok := client.rooms[room]; ok {
		delete(client.rooms, room)
	}

	room.unregister <- client
}

func ServeWs(wsServer *WsServer, w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userid")
	room := r.URL.Query().Get("room")

	if userID == "" || room == "" {
		log.Println("Missing 'userID' or 'room' query parameter")
		http.Error(w, "Missing 'userID' or 'room' query parameter", http.StatusBadRequest)
		return
	}

	log.Printf("New WebSocket connection: userID=%s, room=%s", userID, room)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error while upgrading connection:", err)
		return
	}

	client := newClient(conn, wsServer, userID)

	go client.writePump()
	go client.readPump()
	wsServer.register <- client
}
