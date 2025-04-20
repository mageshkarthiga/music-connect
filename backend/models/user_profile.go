package models

// UserProfile struct for holding user preferences and playlist data
type UserProfile struct {
    ID             uint
    LikedTracks    map[uint]int // TrackID -> Play count (only for liked tracks)
    PlayedTracks   map[uint]int // TrackID -> Play count
    PlaylistTracks map[uint]int // TrackID -> Playlist count
}
