package models

type UserEvent struct {
    UserID  uint `json:"user_id"`
    EventID uint `json:"event_id"`
}
