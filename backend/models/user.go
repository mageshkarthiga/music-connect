package models

type User struct {
	UserID          uint      `json:"user_id" gorm:"primaryKey;autoIncrement"`
	PhoneNumber     string    `json:"phone_number"`
	EmailAddress    string    `json:"email_address"`
	Location        string    `json:"location"`
	UserName        string    `json:"user_name"`
	ProfilePhotoUrl string    `json:"profile_photo_url"`
	FirebaseUID     string    `json:"firebase_uid" gorm:"uniqueIndex:uni_users_firebase_uid"`

	Events    []*Event    `json:"events" gorm:"many2many:user_events;joinForeignKey:UserID;joinReferences:EventID"`
	Tracks    []*Track    `json:"tracks" gorm:"many2many:music_preferences;joinForeignKey:UserID;joinReferences:TrackID"`

	Friends []*User `json:"friends" gorm:"many2many:user_friends;joinForeignKey:UserID;joinReferences:FriendID"`

}
