package models


type Artist struct {
	ArtistID       uint `json:"artist_id" gorm:"primaryKey;autoIncrement"`
	ArtistName     string `json:"artist_name"`
	ArtistImageURL string `json:"artist_image_url"`
	Tracks         []*Track `json:"tracks" gorm:"many2many:track_artists;joinForeignKey:ArtistID;joinReferences:TrackID"`
}
