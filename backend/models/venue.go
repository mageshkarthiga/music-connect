package models

type Venue struct {
	VenueID       uint `json:"venue_id" gorm:"primaryKey;autoIncrement"`
	VenueName     string `json:"venue_name"`
	Location      string `json:"location"`
	SeatMap       string `json:"seat_map"` // URL to the seat map image
	Accessibility string `json:"accessibility"` // Information about accessibility options
	Events        []*Event `json:"events" gorm:"many2many:event_venues;joinForeignKey:VenueID;joinReferences:EventID"`
}