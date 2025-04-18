package controllers

import (
	"backend/config"
	"backend/models"
	"backend/services/recommender"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetTracks fetches all tracks
func GetTracks(c echo.Context) error {
	var tracks []models.Track

	// Get pagination parameters from the query string
	page := c.QueryParam("page")
	limit := c.QueryParam("limit")

	// Default values for pagination
	pageNumber := 1
	pageSize := 25

	if page != "" {
		fmt.Sscanf(page, "%d", &pageNumber)
	}
	if limit != "" {
		fmt.Sscanf(limit, "%d", &pageSize)
	}

	offset := (pageNumber - 1) * pageSize

	// Fetch tracks with pagination
	if err := config.DB.Limit(pageSize).Offset(offset).Find(&tracks).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to fetch tracks")
	}

	return c.JSON(http.StatusOK, tracks)
}

// GetTrackByID fetches a single track by ID
func GetTrackByID(c echo.Context) error {
	id := c.Param("id")
	var track models.Track
	if err := config.DB.First(&track, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, "Track not found")
	}
	return c.JSON(http.StatusOK, track)
}

// CreateTrack creates a new track
func CreateTrack(c echo.Context) error {
	var track models.Track
	if err := c.Bind(&track); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := config.DB.Create(&track).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to create track")
	}
	return c.JSON(http.StatusCreated, track)
}

// UpdateTrack updates an existing track by ID
func UpdateTrack(c echo.Context) error {
	id := c.Param("id")
	var track models.Track
	if err := config.DB.First(&track, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, "Track not found")
	}

	if err := c.Bind(&track); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := config.DB.Save(&track).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to update track")
	}
	return c.JSON(http.StatusOK, track)
}

// DeleteTrack deletes a track by ID
func DeleteTrack(c echo.Context) error {
	id := c.Param("id")
	var track models.Track
	if err := config.DB.First(&track, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, "Track not found")
	}

	if err := config.DB.Delete(&track).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to delete track")
	}
	return c.JSON(http.StatusOK, "Track deleted successfully")
}

func GetTrackRecommendation(c echo.Context) error {
    uid, ok := c.Get("uid").(uint)
    if !ok {
		return c.JSON(http.StatusUnauthorized, "Missing or invalid user ID in context")
	}

    uidStr := fmt.Sprintf("%d", uid)
    listOfSongs, err  := recommender.GetTrackRecommendation(uidStr)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, "Failed to fetch track recommendations")
    }

    return c.JSON(http.StatusOK, listOfSongs)
}

// GetTracksForUser fetches all tracks for a specific user based on their preferences
func GetTracksForUser(c echo.Context) error {
    userID := c.Get("uid").(uint)  // Assuming the user ID is stored in the context

    // Find the track IDs that the user has selected in their MusicPreferences
    var musicPreferences []models.MusicPreference
    if err := config.DB.Where("user_id = ?", userID).Find(&musicPreferences).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, "Failed to fetch music preferences")
    }

    // Extract TrackIDs from MusicPreferences
    var trackIDs []uint
    for _, mp := range musicPreferences {
        trackIDs = append(trackIDs, mp.TrackID)
    }

    // Fetch tracks that are associated with the user through MusicPreferences
    var tracks []models.Track
    if err := config.DB.Where("track_id IN ?", trackIDs).Find(&tracks).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, "Failed to fetch tracks for user")
    }

    return c.JSON(http.StatusOK, tracks)
}

func AddTracksForUser(c echo.Context) error {
    userID := c.Get("uid").(uint)  // Assuming the user ID is stored in the context

    // Bind the request data to a slice of MusicPreference objects
    var musicPreferences []models.MusicPreference
    if err := c.Bind(&musicPreferences); err != nil {
        return c.JSON(http.StatusBadRequest, "Invalid request data")
    }

    // Iterate over each MusicPreference and ensure the TrackID is valid
    for _, musicPreference := range musicPreferences {
        var track models.Track
        if err := config.DB.First(&track, musicPreference.TrackID).Error; err != nil {
            return c.JSON(http.StatusNotFound, "Track not found")
        }

        // Set the UserID for each MusicPreference
        musicPreference.UserID = userID

        // Insert the MusicPreference into the database
        if err := config.DB.Create(&musicPreference).Error; err != nil {
            return c.JSON(http.StatusInternalServerError, "Failed to add track for user")
        }
    }

    return c.JSON(http.StatusCreated, musicPreferences)
}

// returns tracks for the given user ID.
func GetUserTracksByID(c echo.Context) error {
    uid := c.Param("id")

    var prefs []models.MusicPreference
    if err := config.DB.Select("track_id").Where("user_id = ?", uid).Find(&prefs).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, "DB error")
    }
    if len(prefs) == 0 {
        return c.JSON(http.StatusOK, []models.Track{})
    }

    ids := make([]uint, len(prefs))
    for i, p := range prefs {
        ids[i] = p.TrackID
    }

    var tracks []models.Track
    if err := config.DB.Where("track_id IN ?", ids).Find(&tracks).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, "DB error")
    }
    return c.JSON(http.StatusOK, tracks)
}


func LikeTrack(c echo.Context) error {
	userID := c.Get("uid").(uint)  // Assuming the user ID is stored in the context
	trackID := c.Param("track_id")

	// Check if the track exists
	var track models.Track
	if err := config.DB.First(&track, trackID).Error; err != nil {
		return c.JSON(http.StatusNotFound, "Track not found")
	}

	// Create a new MusicPreference for the "like" action
	musicPreference := models.MusicPreference{
		UserID:  userID,
		TrackID: track.TrackID,
	}

	// Check if the user already liked the track
	var existingPreference models.MusicPreference
	if err := config.DB.Where("user_id = ? AND track_id = ?", userID, trackID).First(&existingPreference).Error; err == nil {
		return c.JSON(http.StatusConflict, "Track already liked")
	}
	// Insert the MusicPreference into the database
	if err := config.DB.Create(&musicPreference).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to like track")
	}
	return c.JSON(http.StatusCreated, "Track liked successfully")

}

func UnlikeTrack(c echo.Context) error {
	userID := c.Get("uid").(uint)  // Assuming the user ID is stored in the context
	trackID := c.Param("track_id")

	// Delete the MusicPreference for the "unlike" action
	if err := config.DB.Where("user_id = ? AND track_id = ?", userID, trackID).Delete(&models.MusicPreference{}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to unlike track")
	}

	return c.JSON(http.StatusOK, "Track unliked successfully")
}

func GetLikedTracks(c echo.Context) error {
	userID := c.Get("uid").(uint)  // Assuming the user ID is stored in the context

	// Fetch all liked tracks for the user
	var likedTracks []models.MusicPreference
	if err := config.DB.Where("user_id = ? AND preference_type = ?", userID, "like").Find(&likedTracks).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to fetch liked tracks")
	}

	// Extract TrackIDs from MusicPreferences
	var trackIDs []uint
	for _, mp := range likedTracks {
		trackIDs = append(trackIDs, mp.TrackID)
	}

	// Fetch tracks that are associated with the user through MusicPreferences
	var tracks []models.Track
	if err := config.DB.Where("track_id IN ?", trackIDs).Find(&tracks).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to fetch tracks for user")
	}

	return c.JSON(http.StatusOK, tracks)
}