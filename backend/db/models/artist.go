package models

type Artist struct {
	ArtistID       uint `gorm:"primaryKey;autoIncrement"`
	ArtistName     string
	ArtistImageURL string
	Tracks         []*Track `gorm:"many2many:track_artists;joinForeignKey:ArtistID;joinReferences:TrackID"`
}