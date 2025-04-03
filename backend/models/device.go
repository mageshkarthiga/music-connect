
package models

type Device struct {
	DeviceID   uint `gorm:"primaryKey;autoIncrement"`
	DeviceName string
	DeviceUUID string `gorm:"uniqueIndex"`
	UserID     uint 
}