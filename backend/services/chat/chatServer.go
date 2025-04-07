package chat

import (
	"log"
)

type WsServer struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	rooms      map[*Room]bool
}

// Create new instance of WsServer
func NewWsServer() *WsServer {
	return &WsServer{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		rooms: make(map[*Room]bool),
	}
}

// Run WebSocket server and accept various requests
func (s *WsServer) Run() {
	for {
		select {
		case client := <-s.register:
			s.registerClient(client)
		case client := <-s.unregister:
			s.unregisterClient(client)
		case message := <-s.broadcast:
			s.broadcastToClients(message)
		}
	}
}

func (s *WsServer) registerClient(client *Client) {
	s.clients[client] = true
}

func (s *WsServer) unregisterClient(client *Client) {
	if _, ok := s.clients[client]; ok {
		delete(s.clients, client)
	}
}

// Listens for messages sent by client readPump and writes them to the broadcast channel
func (s *WsServer) broadcastToClients(message []byte) {
	for client := range s.clients {
		client.send <- message
	}
}

func (s *WsServer) createRoom(name string) *Room {
    room := NewRoom(name)
    log.Printf("Starting RunRoom for room: %s\n", name)
    go room.RunRoom()
    s.rooms[room] = true
    return room
}

func (s *WsServer) findRoom(name string) *Room{
	var foundRoom *Room
	for room := range s.rooms {
		if room.GetName() == name{
			foundRoom = room
			break
		}
	}
	return foundRoom
}
