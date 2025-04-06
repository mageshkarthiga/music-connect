package chat

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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
	Name     string `json:"name"`
}

func newClient(conn *websocket.Conn, wsServer *WsServer, name string) *Client {
	return &Client{
		conn:     conn,
		wsServer: wsServer,
		send:     make(chan []byte, 256),
		rooms:    make(map[*Room]bool),
		Name: name,
	}
}

func (client *Client) GetName() string {
	return client.Name
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

	testMessage := &Message{
        Action:  SendMessageAction,
        Message: "Hello from the server!",
        Target:  "Alice", // Replace with the target user
        Sender:  client,
    }
    client.wsServer.broadcast <- testMessage.encode()


	for {
		_, jsonMessage, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println("Error while reading message:", err)
			}
			break
		}
		client.handleNewMessage(jsonMessage) // Handle the new message
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
			client.conn.SetWriteDeadline(time.Now().Add(writeWait)) // Set the write deadline
			if !ok {
				// The WebSocket connection is closed.
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
	}

	message.Sender = client
	switch message.Action {
	case SendMessageAction:
		// send the message to a specific room, which depends on message target
		roomName := message.Target
		if room := client.wsServer.findRoom(roomName); room != nil {
			room.broadcast <- &message
		} else {
			log.Printf("Room %s not found", roomName)
		}
	case JoinRoomAction:
		client.handleJoinRoom(message)
	case LeaveRoomAction:
		client.handleLeaveRoom(message)
	}
}

// client.go
func (client *Client) handleJoinRoom(message Message) {
	roomName := message.Message

	room := client.wsServer.findRoom(roomName)
	if room == nil {
		room = client.wsServer.createRoom(roomName)
	}

	client.rooms[room] = true

	room.register <- client
}

func (client *Client) handleLeaveRoom(message Message) {
	room := client.wsServer.findRoom(message.Message)
	if _, ok := client.rooms[room]; ok {
		delete(client.rooms, room)
	}

	room.unregister <- client
}

func ServeWs(wsServer *WsServer, w http.ResponseWriter, r *http.Request) {
	name, ok := r.URL.Query()["name"]

	if !ok || len(name[0]) < 1 {
		log.Println("Url Param 'name' is missing")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error while upgrading connection:", err)
		return
	}

	client := newClient(conn, wsServer,name[0])

	go client.writePump()
	go client.readPump()
	wsServer.register <- client
}
