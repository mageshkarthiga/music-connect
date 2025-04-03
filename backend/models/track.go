package models

type Track struct {
	TrackID    uint `gorm:"primaryKey;autoIncrement"`
	TrackTitle string `json:"track_title"`
	ArtistID   uint `gorm:"foreignKey:artist_id"`
	Genre      string `json:"genre"`
	TrackURI   string `json:"track_uri"`
	Artists    []*Artist   `gorm:"many2many:track_artists;joinForeignKey:TrackID;joinReferences:ArtistID"`
	Playlists  []*Playlist `gorm:"many2many:playlist_tracks;joinForeignKey:TrackID;joinReferences:PlaylistID"`
}