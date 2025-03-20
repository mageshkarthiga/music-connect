package models

type Event struct {
	EventID          uint   `gorm:"primaryKey;autoIncrement"`
	EventName        string
	EventDescription string
	EventURL         string
	EventImageURL    string
	Users            []*User       `gorm:"many2many:user_events;joinForeignKey:EventID;joinReferences:UserID"`
	Attractions      []*Attraction `gorm:"many2many:event_attractions;joinForeignKey:EventID;joinReferences:AttractionID"`
	Venues           []*Venue      `gorm:"many2many:event_venues;joinForeignKey:EventID;joinReferences:VenueID"`
}
