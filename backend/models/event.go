package models


type Event struct {
	EventID          uint `gorm:"primaryKey;autoIncrement"`
	EventName        string `json:"event_name"`
	EventDescription string `json:"event_description"`
	EventURL         string `json:"event_url"`
	EventImageURL    string `json:"event_image_url"`
	Users            []*User       `gorm:"many2many:user_events;joinForeignKey:EventID;joinReferences:UserID"`
	Attractions      []*Attraction `gorm:"many2many:event_attractions;joinForeignKey:EventID;joinReferences:AttractionID"`
	Venues           []*Venue      `gorm:"many2many:event_venues;joinForeignKey:EventID;joinReferences:VenueID"`
}
