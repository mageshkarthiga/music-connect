package server

import (
	"log"
	"sync"
)

type WsServer struct {
	Clients    map[*Client]bool
	Rooms      map[string]*Room
	Register   chan *Client // Exported field
	Unregister chan *Client // Exported field
	Broadcast  chan []byte
	mu         sync.Mutex
}

// NewWsServer creates a new WebSocket server instance.
func NewWsServer() *WsServer {
	return &WsServer{
		Clients:    make(map[*Client]bool),
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}
}

// Run starts the WebSocket server and listens for events.
func (s *WsServer) Run() {
	for {
		select {
		case client := <-s.Register:
			s.registerClient(client)
		case client := <-s.Unregister:
			s.unregisterClient(client)
		case message := <-s.Broadcast:
			s.broadcastToClients(message)
		}
	}
}

// registerClient adds a client to the server's client list.
func (s *WsServer) registerClient(client *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Clients[client] = true
	log.Printf("Client registered: %s", client.UserID)
}

// unregisterClient removes a client from the server's client list.
func (s *WsServer) unregisterClient(client *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.Clients[client]; ok {
		delete(s.Clients, client)
		close(client.send)
		log.Printf("Client unregistered: %s", client.UserID)
	}
}

// broadcastToClients sends a message to all connected clients.
func (s *WsServer) broadcastToClients(message []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for client := range s.Clients {
		select {
		case client.send <- message:
		default:
			log.Printf("Failed to send message to client %s", client.UserID)
		}
	}
}

// CreateRoom creates a new room if it doesn't already exist.
func (s *WsServer) CreateRoom(name string) *Room {
	s.mu.Lock()
	defer s.mu.Unlock()
	room, exists := s.Rooms[name]
	if !exists {
		room = NewRoom(name)
		s.Rooms[name] = room
		go room.Run()
		log.Printf("Room created: %s", name)
	}
	return room
}

// FindRoom finds a room by its name.
func (s *WsServer) FindRoom(name string) *Room {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Rooms[name]
}
