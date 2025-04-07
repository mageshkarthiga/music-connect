package chat

import (
	"encoding/json"
	"log"
)

const SendMessageAction = "send-message"
const JoinRoomAction = "join-room"
const LeaveRoomAction = "leave-room"

type Message struct {
	Action  string `json:"action"`
	Message string `json:"message"`
	Target  string `json:"target"`
	Sender  string `json:"sender"` // Use user_id instead of Client
}

func (message *Message) encode() []byte {
	type MessageDTO struct {
		Action  string `json:"action"`
		Message string `json:"message"`
		Target  string `json:"target"`
		Sender  string `json:"sender"`
	}

	dto := MessageDTO{
		Action:  message.Action,
		Message: message.Message,
		Target:  message.Target,
		Sender:  message.Sender, // Use user_id
	}

	jsonMessage, err := json.Marshal(dto)
	if err != nil {
		log.Println(err)
	}
	return jsonMessage
}