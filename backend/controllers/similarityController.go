package controllers

import (
	"backend/models"
	"backend/config"
	"fmt"
	"net/http"
	"strconv"
	"log"

	"github.com/labstack/echo/v4"
)

// CreateUserProfile creates a user profile by fetching from existing data in the tables
func CreateUserProfile(userID uint) (models.UserProfile, error) {
	var musicPreferences []models.MusicPreference
	var playlistTracks []models.PlaylistTrack

	// Fetch music preferences for the user using GORM
	if err := config.DB.Where("user_id = ?", userID).Find(&musicPreferences).Error; err != nil {
		return models.UserProfile{}, fmt.Errorf("error fetching music preferences: %w", err)
	}

	// Fetch playlist tracks for the user using raw SQL query with GORM
	if err := config.DB.Raw(`
		SELECT pt.track_id
		FROM playlist_tracks pt
		JOIN playlists up ON pt.playlist_id = up.playlist_id
		WHERE up.user_id = ?`, userID).Scan(&playlistTracks).Error; err != nil {
		return models.UserProfile{}, fmt.Errorf("error fetching playlist tracks: %w", err)
	}

	// If no data is found, return an empty profile (not an error)
	if len(musicPreferences) == 0 && len(playlistTracks) == 0 {
		return models.UserProfile{
			ID:             userID,
			LikedTracks:    make(map[uint]int),
			PlayedTracks:   make(map[uint]int),
			PlaylistTracks: make(map[uint]int),
		}, nil // Return empty profile with no data
	}

	// Construct the user profile from the fetched data
	userProfile := models.UserProfile{
		ID:             userID,
		LikedTracks:    make(map[uint]int),
		PlayedTracks:   make(map[uint]int),
		PlaylistTracks: make(map[uint]int),
	}

	// Fill in the user profile data from music preferences
	for _, pref := range musicPreferences {
		if pref.IsLiked {
			userProfile.LikedTracks[pref.TrackID] = pref.PlayCount
		}
		userProfile.PlayedTracks[pref.TrackID] = pref.PlayCount
	}

	// Fill in the user profile data from playlist tracks
	for _, track := range playlistTracks {
		userProfile.PlaylistTracks[track.TrackID] = 1 // Assuming one entry per playlist
	}

	return userProfile, nil
}

// CosineSimilarity calculates cosine similarity between two user profiles
func CosineSimilarity(userProfile1, userProfile2 models.UserProfile) (float64, error) {
	similarity, err := models.CosineSimilarity(userProfile1, userProfile2)
	if err != nil {
		return 0, err
	}
	return similarity, nil
}

// CalculateUserSimilarity calculates the similarity between two users by fetching their profiles
func CalculateUserSimilarity(userID1, userID2 uint) (float64, error) {
	log.Printf("Creating profile for user %d", userID1)
	userProfile1, err := CreateUserProfile(userID1)
	if err != nil {
		log.Printf("Error in CreateUserProfile (user %d): %v", userID1, err)
		return 0, err
	}

	log.Printf("Creating profile for user %d", userID2)
	userProfile2, err := CreateUserProfile(userID2)
	if err != nil {
		log.Printf("Error in CreateUserProfile (user %d): %v", userID2, err)
		return 0, err
	}

	log.Printf("Calculating cosine similarity")
	similarity, err := CosineSimilarity(userProfile1, userProfile2)
	if err != nil {
		log.Printf("Error in CosineSimilarity: %v", err)
		return 0, err
	}

	log.Printf("Similarity calculated: %f", similarity)
	return similarity, nil
}


// GetUserProfileHandler handles the HTTP request to fetch a user's profile
func GetUserProfileHandler(c echo.Context) error {
	userID, err := strconv.ParseUint(c.QueryParam("user_id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user_id"})
	}

	userProfile, err := CreateUserProfile(uint(userID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error fetching user profile"})
	}

	return c.JSON(http.StatusOK, userProfile)
}

// CalculateUserSimilarityHandler handles the HTTP request to calculate similarity between two users
func CalculateUserSimilarityHandler(c echo.Context) error {
	userID1, err := strconv.ParseUint(c.QueryParam("user_id1"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user_id1"})
	}

	userID2, err := strconv.ParseUint(c.QueryParam("user_id2"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user_id2"})
	}

	similarity, err := CalculateUserSimilarity(uint(userID1), uint(userID2))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error calculating similarity"})
	}

	return c.JSON(http.StatusOK, map[string]float64{"similarity": similarity})
}
