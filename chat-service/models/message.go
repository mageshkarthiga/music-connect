package models

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

const SendMessageAction = "send-message"
const JoinRoomAction = "join-room"
const LeaveRoomAction = "leave-room"

type Message struct {
	Action    string `json:"action"`
	Message   string `json:"message"`
	Target    string `json:"target"`
	Sender    string `json:"sender"`
	MessageID string `json:"message_id,omitempty"`
}

func (message *Message) Encode() []byte {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}
	return jsonMessage
}

func (message *Message) GenerateMessageID() {
	message.MessageID = uuid.New().String()
}
