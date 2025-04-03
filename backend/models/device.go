
package models

type Device struct {
	DeviceID   uint `gorm:"primaryKey;autoIncrement"`
	DeviceName string `json:"device_name"`
	DeviceUUID string `gorm:"uniqueIndex"`
	UserID     uint `json:"user_id"`
}