package models

import (
	"time"
)

type Friendship struct {
	UserID    uint      `json:"user_id" gorm:"primaryKey"`
	FriendID  uint      `json:"friend_id" gorm:"primaryKey"`
	Status    string    `json:"status"` // pending, accepted, etc.
	CreatedAt time.Time `json:"created_at"`

	User   User `json:"user" gorm:"foreignKey:UserID;references:UserID"`
	Friend User `json:"friend" gorm:"foreignKey:FriendID;references:UserID"`
}

func (Friendship) TableName() string {
	return "friendships"
}
