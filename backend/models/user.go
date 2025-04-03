package models


type User struct {
	UserID       uint   `json:"user_id" gorm:"primaryKey;autoIncrement"`
	PhoneNumber  string `json:"phone_number" gorm:"size:11;unique"`
	EmailAddress string `json:"email_address"`
	Location     string `json:"location"`
	UserName     string `json:"user_name"`
	Events       []*Event `json: "events" gorm:"many2many:user_events;joinForeignKey:UserID;joinReferences:EventID"`
	Tracks       []*Track `json: "tracks" gorm:"many2many:music_preferences;joinForeignKey:UserID;joinReferences:TrackID"`
}

