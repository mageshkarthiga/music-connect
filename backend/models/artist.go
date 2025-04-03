package models


type Artist struct {
	ArtistID       uint `gorm:"primaryKey;autoIncrement"`
	ArtistName     string `json:"artist_name"`
	ArtistImageURL string `json:"artist_image_url"`
	Tracks         []*Track `gorm:"many2many:track_artists;joinForeignKey:ArtistID;joinReferences:TrackID"`
}
