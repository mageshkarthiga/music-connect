package models

type Playlist struct {
	PlaylistID   uint `json:"playlist_id" gorm:"primaryKey;autoIncrement"`
	PlaylistName string `json:"playlist_name"`
	Tracks       []*Track `json:"tracks" gorm:"many2many:playlist_tracks;joinForeignKey:PlaylistID;joinReferences:TrackID"`
}

