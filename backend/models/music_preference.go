package models

type MusicPreference struct {
    UserID   uint `json:"user_id" gorm:"primaryKey;index:idx_track_user,unique"`
    TrackID  uint `json:"track_id" gorm:"primaryKey;index:idx_track_user,unique"`
    IsLiked  bool `json:"is_liked" gorm:"default:false"`
    PlayCount int `json:"play_count" gorm:"default:0"`
}
