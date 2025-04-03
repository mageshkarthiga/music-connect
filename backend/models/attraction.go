package models

type Attraction struct {
	AttractionID       uint `gorm:"primaryKey;autoIncrement"`
	AttractionName     string `json:"attraction_name"`
	AttractionType     string `json:"attraction_type"`
	AttractionURL      string `json:"attraction_url"`
	AttractionImageURL string `json:"attraction_image_url"`
	Events             []*Event `gorm:"many2many:event_attractions;joinForeignKey:AttractionID;joinReferences:EventID"`
}
