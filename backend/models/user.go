package models


type User struct {
	UserID       uint   `gorm:"primaryKey;autoIncrement"`
	PhoneNumber  string `json:"phone_number" gorm:"size:11;unique"`
	EmailAddress string `json:"email_address"`
	Location     string `json:"location"`
	UserName     string `json:"user_name"`
	Events       []*Event `gorm:"many2many:user_events;joinForeignKey:UserID;joinReferences:EventID"`
	Tracks       []*Track `gorm:"many2many:music_preferences;joinForeignKey:UserID;joinReferences:TrackID"`
	Devices      []Device `json:"devices"`
}

