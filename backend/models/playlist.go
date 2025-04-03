package models

type Playlist struct {
	PlaylistID   uint `gorm:"primaryKey;autoIncrement"`
	PlaylistName string
	Tracks       []*Track `gorm:"many2many:playlist_tracks;joinForeignKey:PlaylistID;joinReferences:TrackID"`
}

