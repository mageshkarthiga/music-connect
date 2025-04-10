package chat

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
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

// const welcomeMessage = "%s joined the room"

// func (room *Room) notifyClientJoined(client *Client) {
// 	message := &Message{
// 		Action:  SendMessageAction,
// 		Target:  room.name,
// 		Message: fmt.Sprintf(welcomeMessage, client.UserID),
// 	}

// 	room.broadcastMessageToClients(message.encode())
// }

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
	// room.notifyClientJoined(client)
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
	var msg Message
	if err := json.Unmarshal(message, &msg); err != nil {
		fmt.Printf("Error unmarshalling message: %v\n", err)
		return
	}

	// Generate a message ID if one doesn't exist
	if msg.MessageID == "" {
		msg.GenerateMessageID()
		// Re-encode with the message ID
		message = msg.encode()
	}

	// Save to Firestore only once
	go func(msg Message) {
		_, _, err := FirestoreClient.Collection("rooms").
			Doc(room.name).
			Collection("messages").
			Add(context.Background(), map[string]interface{}{
				"sender":     msg.Sender,
				"message":    msg.Message,
				"target":     msg.Target,
				"action":     msg.Action,
				"message_id": msg.MessageID,
				"timestamp":  firestore.ServerTimestamp,
			})
		if err != nil {
			fmt.Printf("Failed to save message to Firestore: %v\n", err)
		} else {
			fmt.Printf("Message saved to Firestore for room %s\n", room.name)
		}
	}(msg)

	// Broadcast to all clients
	for client := range room.clients {
		// Skip sending the message back to the sender
		if client.UserID == msg.Sender {
			continue
		}

		fmt.Printf("Sending message to client %s\n", client.UserID)
		select {
		case client.send <- message:
			fmt.Printf("Message sent to client %s\n", client.UserID)
		default:
			fmt.Printf("Failed to send message to client %s (channel blocked)\n", client.UserID)
		}
	}
}

func GetMessagesForRoom(roomName string) ([]Message, error) {
	ctx := context.Background()
	messages := []Message{}

	// Query the messages subcollection for the room
	iter := FirestoreClient.Collection("rooms").
		Doc(roomName).
		Collection("messages").
		OrderBy("timestamp", firestore.Asc).
		Documents(ctx)

	for {
		doc, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			return nil, fmt.Errorf("failed to retrieve messages: %v", err)
		}

		var msg Message
		if err := doc.DataTo(&msg); err != nil {
			return nil, fmt.Errorf("failed to parse message: %v", err)
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (room *Room) GetName() string {
	return room.name
}
