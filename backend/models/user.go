package models


type User struct {
	UserID       uint   `gorm:"primaryKey;autoIncrement"`
	PhoneNumber  string `gorm:"size:11;unique"`
	EmailAddress string
	Location     string
	UserName     string
	Events       []*Event `gorm:"many2many:user_events;joinForeignKey:UserID;joinReferences:EventID"`
	Tracks       []*Track `gorm:"many2many:music_preferences;joinForeignKey:UserID;joinReferences:TrackID"`
	Devices      []Device
}