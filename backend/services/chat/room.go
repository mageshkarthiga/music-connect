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

// room.go
const welcomeMessage = "%s joined the room"

func (room *Room) notifyClientJoined(client *Client) {
	message := &Message{
		Action:  SendMessageAction,
        Target:  room.name,
		Message: fmt.Sprintf(welcomeMessage, client.GetName()),
	}

	room.broadcastMessageToClients(message.encode())
}

// RunRoom to run the room and accept various requests
func (room *Room) RunRoom() {
	for {
		select {
		case client := <-room.register:
			room.registerClientInRoom(client)
		case client := room.unregister:
			room.unregisterClientInRoom(client)
		case message := room.broadcast:
			room.broadcastMessageToClients(message.encode())
		}
	}
}

func (room *Room) registerClientInRoom(client *Client) {
	room.notifyClientJoined(client) // send message in the chat to indicate that a new user has joined
	room.clients[client] = true
	fmt.Printf("Client %s joined room %s\n", client.name, room.name)
}

func (room *Room) unregisterClientInRoom(client *Client) {
	if _,ok := room.clients[client]; ok {
		delete(room.clients,client)
	}
}

func (room *Room) broadcastMessageToClients(message []byte) {
	for client := range room.clients {
		client.send <- message
	}
}
