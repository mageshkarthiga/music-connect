package chat

import (
	"fmt"
)

type Room struct {
	name       string
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
}

func NewRoom(name string) *Room {
	return &Room{
		name:       name,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
	}
}

const welcomeMessage = "%s joined the room"

func (room *Room) notifyClientJoined(client *Client) {
	message := &Message{
		Action:  SendMessageAction,
		Target:  room.name,
		Message: fmt.Sprintf(welcomeMessage, client.UserID),
	}

	room.broadcastMessageToClients(message.encode())
}

// RunRoom to run the room and accept various requests
func (room *Room) RunRoom() {
    fmt.Printf("RunRoom started for room: %s\n", room.name)
    for {
        select {
        case client := <-room.register:
            fmt.Printf("Registering client %s in room %s\n", client.UserID, room.name)
            room.registerClientInRoom(client)
        case client := <-room.unregister:
            fmt.Printf("Unregistering client %s from room %s\n", client.UserID, room.name)
            room.unregisterClientInRoom(client)
        case message := <-room.broadcast:
            fmt.Printf("Broadcasting message in room %s: %s\n", room.name, string(message.encode()))
            room.broadcastMessageToClients(message.encode())
        }
    }
}

func (room *Room) registerClientInRoom(client *Client) {
    room.notifyClientJoined(client)
    room.clients[client] = true
    fmt.Printf("Client %s joined room %s\n", client.UserID, room.name)
    fmt.Printf("Current clients in room %s: ", room.name)
    for c := range room.clients {
        fmt.Printf("%s ", c.UserID)
    }
    fmt.Println()
}

func (room *Room) unregisterClientInRoom(client *Client) {
	if _, ok := room.clients[client]; ok {
		delete(room.clients, client)
	}
}

func (room *Room) broadcastMessageToClients(message []byte) {
    fmt.Printf("Broadcasting message to room %s: %s\n", room.name, string(message))
    for client := range room.clients {
        fmt.Printf("Sending message to client %s\n", client.UserID)
        select {
        case client.send <- message:
            fmt.Printf("Message sent to client %s\n", client.UserID)
        default:
            fmt.Printf("Failed to send message to client %s (channel blocked)\n", client.UserID)
        }
    }
}

func (room *Room) GetName() string {
	return room.name
}
