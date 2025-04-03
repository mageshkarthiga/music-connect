package models

type Attraction struct {
	AttractionID       uint `gorm:"primaryKey;autoIncrement"`
	AttractionName     string
	AttractionType     string
	AttractionURL      string
	AttractionImageURL string
	Events             []*Event `gorm:"many2many:event_attractions;joinForeignKey:AttractionID;joinReferences:EventID"`
}
