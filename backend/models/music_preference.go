package models

type MusicPreference struct {
    UserID  uint `json:"user_id"`
    TrackID uint `json:"track_id"`
    IsLiked bool `json:"is_liked"`
    PlayCount int `json:"play_count"`
}
