package models

type Track struct {
	TrackID    uint `json:"track_id" gorm:"primaryKey;autoIncrement"`
	TrackTitle string `json:"track_title"`
	ArtistID   uint `json:"artist_id" gorm:"foreignKey:artist_id"`
	Genre      string `json:"genre"`
	TrackURI   string `json:"track_uri"`
	TrackImageUrl string `json:"track_image_url"`
	Artists    []*Artist   `json:"artists" gorm:"many2many:track_artists;joinForeignKey:TrackID;joinReferences:ArtistID"`
	Playlists  []*Playlist `json:"playlists" gorm:"many2many:playlist_tracks;joinForeignKey:TrackID;joinReferences:PlaylistID"`
	TrackSpotifyID string `json:"track_spotify_id"`
	TrackImageURL string `json:"track_image_url"`
}