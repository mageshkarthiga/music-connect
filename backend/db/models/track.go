package models

type Track struct {
	TrackID    uint `gorm:"primaryKey;autoIncrement"`
	TrackTitle string
	ArtistID   uint `gorm:"foreignKey:artist_id"`
	Genre      string
	TrackURL   string
	Artists    []*Artist   `gorm:"many2many:track_artists;joinForeignKey:TrackID;joinReferences:ArtistID"`
	Playlists  []*Playlist `gorm:"many2many:playlist_tracks;joinForeignKey:TrackID;joinReferences:PlaylistID"`
}