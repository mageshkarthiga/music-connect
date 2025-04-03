package models


type Event struct {
	EventID          uint `json:"event_id" gorm:"primaryKey;autoIncrement"`
	EventName        string `json:"event_name"`
	EventDescription string `json:"event_description"`
	EventURL         string `json:"event_url"`
	EventImageURL    string `json:"event_image_url"`
	Users            []*User       `json:"users" gorm:"many2many:user_events;joinForeignKey:EventID;joinReferences:UserID"`
	Attractions      []*Attraction `json:"attractions" gorm:"many2many:event_attractions;joinForeignKey:EventID;joinReferences:AttractionID"`
	Venues           []*Venue      `json:"venues" gorm:"many2many:event_venues;joinForeignKey:EventID;joinReferences:VenueID"`
}
