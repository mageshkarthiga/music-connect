package models

type Venue struct {
	VenueID       uint `gorm:"primaryKey;autoIncrement"`
	VenueName     string
	Location      string
	SeatMap       string
	Accessibility string
	Events        []*Event `gorm:"many2many:event_venues;joinForeignKey:VenueID;joinReferences:EventID"`
}