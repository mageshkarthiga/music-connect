package models

type Attraction struct {
	AttractionID       uint `json "attraction_id" gorm:"primaryKey;autoIncrement"`
	AttractionName     string `json:"attraction_name"`
	AttractionType     string `json:"attraction_type"`
	AttractionURL      string `json:"attraction_url"`
	AttractionImageURL string `json:"attraction_image_url"`
	Events             []*Event `json "events" gorm:"many2many:event_attractions;joinForeignKey:AttractionID;joinReferences:EventID"`
}
