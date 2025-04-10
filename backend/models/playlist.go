package models

type Playlist struct {
	PlaylistID       uint    `json:"playlist_id" gorm:"primaryKey;autoIncrement"`
	PlaylistName     string  `json:"playlist_name"`
	PlaylistImageURL string  `json:"playlist_image_url"`
	UserID           uint    `json:"user_id"`
	Tracks           []Track `json:"tracks" gorm:"many2many:playlist_tracks;joinForeignKey:PlaylistID;joinReferences:TrackID"`
}
