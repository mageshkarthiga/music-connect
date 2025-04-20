package server

import (
	"chat-service/firebase"
	"chat-service/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type Room struct {
	Name       string
	Clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *models.Message
}

func NewRoom(name string) *Room {
	return &Room{
		Name:       name,
		Clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *models.Message),
	}
}

// RunRoom to run the room and accept various requests
func (room *Room) Run() {
	fmt.Printf("RunRoom started for room: %s\n", room.Name)
	for {
		select {
		case client := <-room.register:
			fmt.Printf("Registering client %s in room %s\n", client.UserID, room.Name)
			room.registerClient(client)
		case client := <-room.unregister:
			fmt.Printf("Unregistering client %s from room %s\n", client.UserID, room.Name)
			room.unregisterClient(client)
		case message := <-room.broadcast:
			fmt.Printf("Broadcasting message in room %s: %s\n", room.Name, string(message.Encode()))
			room.broadcastMessage(message.Encode())
		}
	}
}

func (room *Room) registerClient(client *Client) {
	room.Clients[client] = true
	fmt.Printf("Client %s joined room %s\n", client.UserID, room.Name)
	fmt.Printf("Current clients in room %s: ", room.Name)
	for c := range room.Clients {
		fmt.Printf("%s ", c.UserID)
	}
	fmt.Println()
}

func (room *Room) unregisterClient(client *Client) {
	if _, ok := room.Clients[client]; ok {
		delete(room.Clients, client)
	}
}

func (room *Room) broadcastMessage(message []byte) {
	fmt.Printf("Broadcasting message to room %s: %s\n", room.Name, string(message))
	var msg models.Message
	if err := json.Unmarshal(message, &msg); err != nil {
		fmt.Printf("Error unmarshalling message: %v\n", err)
		return
	}

	// Generate a message ID if one doesn't exist
	if msg.MessageID == "" {
		msg.GenerateMessageID()
		// Re-encode with the message ID
		message = msg.Encode()
	}

	// Save to Firestore only once
	go func(msg models.Message) {
		_, _, err := firebase.FirestoreClient.Collection("rooms").
			Doc(room.Name).
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
			fmt.Printf("Message saved to Firestore for room %s\n", room.Name)
		}
	}(msg)

	// Broadcast to all clients
	for client := range room.Clients {
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

func GetMessagesForRoom(roomName string) ([]models.Message, error) {
	if firebase.FirestoreClient == nil {
		return nil, fmt.Errorf("Firestore client is not initialized")
	}
	ctx := context.Background()
	messages := []models.Message{}

	// Query the messages subcollection for the room
	iter := firebase.FirestoreClient.Collection("rooms").
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

		var msg models.Message
		if err := doc.DataTo(&msg); err != nil {
			return nil, fmt.Errorf("failed to parse message: %v", err)
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func GetUsersWithChatHistory(userID string) ([]string, error) {
	if firebase.FirestoreClient == nil {
		return nil, fmt.Errorf("Firestore client is not initialized")
	}
	ctx := context.Background()
	otherUsers := []string{}

	iter := firebase.FirestoreClient.Collection("rooms").
		Where("participants", "array-contains", userID).
		Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		data := doc.Data()
		participants, ok := data["participants"].([]interface{})
		if !ok {
			continue
		}

		for _, p := range participants {
			uid, ok := p.(string)
			if ok && uid != userID {
				otherUsers = append(otherUsers, uid)
			}
		}
	}

	return otherUsers, nil
}

func ParseParticipantsFromRoomName(name string) []string {
	parts := strings.Split(name, "-")
	if len(parts) != 2 {
		log.Printf("Unexpected room format: %s", name)
		return []string{}
	}
	log.Printf("Parsed participants from room name: %s", parts)
	return parts
}

func (room *Room) GetName() string {
	return room.Name
}
