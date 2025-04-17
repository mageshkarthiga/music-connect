package controllers

import (
	"backend/config"
	"backend/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"log"
	"strconv"
	"fmt"
)

// GetPlaylists fetches all playlists
func GetPlaylists(c echo.Context) error {
	var playlists []models.Playlist

	// Preload the tracks for each playlist
	if err := config.DB.Preload("Tracks").Find(&playlists).Error; err != nil {
		log.Printf("Error fetching playlists: %v", err)
		return c.JSON(http.StatusInternalServerError, "Failed to fetch playlists")
	}

	return c.JSON(http.StatusOK, playlists)
}

// GetPlaylistByID fetches a single playlist by ID
func GetPlaylistByID(c echo.Context) error {
    id := c.Param("id")
    var playlist models.Playlist
    if err := config.DB.Preload("Tracks").Find(&playlist, "playlist_id = ?", id).Error; err != nil {
        log.Printf("Error fetching playlist with ID %s: %v", id, err)
        return c.JSON(http.StatusNotFound, "Playlist not found")
    }
    return c.JSON(http.StatusOK, playlist)
}



// CreatePlaylist creates a new playlist with the user's existing tracks
func CreatePlaylist(c echo.Context) error {
	var playlistRequest struct {
		UserID   uint   `json:"user_id"`
		Name     string `json:"name"`
		TrackIDs []uint `json:"track_ids"` // IDs of the tracks in the playlist
	}

	// Bind the request body to playlistRequest
	if err := c.Bind(&playlistRequest); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request format")
	}

	// Find the user by ID
	var user models.User
	if err := config.DB.First(&user, playlistRequest.UserID).Error; err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	// Get the tracks from the database based on TrackIDs
	var tracks []models.Track
	if err := config.DB.Where("track_id IN ?", playlistRequest.TrackIDs).Find(&tracks).Error; err != nil {
		log.Printf("Error fetching tracks: %v", err)
		return c.JSON(http.StatusInternalServerError, "Error fetching tracks")
	}

	// Create the playlist
	playlist := models.Playlist{
		PlaylistName: playlistRequest.Name,
		UserID:       playlistRequest.UserID, // Ensure this is the correct type
		Tracks:       tracks, // Store tracks as slice of models.Track (not pointers)
	}

	// Save the playlist to the database
	if err := config.DB.Create(&playlist).Error; err != nil {
		log.Printf("Error saving playlist: %v", err)
		return c.JSON(http.StatusInternalServerError, "Failed to create playlist")
	}

	// Return the created playlist
	return c.JSON(http.StatusOK, playlist)
}


// UpdatePlaylist updates an existing playlist by ID
func UpdatePlaylist(c echo.Context) error {
	id := c.Param("id")
	var playlist models.Playlist
	if err := config.DB.First(&playlist, id).Error; err != nil {
		log.Printf("Error fetching playlist for update with ID %s: %v", id, err)
		return c.JSON(http.StatusNotFound, "Playlist not found")
	}

	// Bind the request data to the playlist
	if err := c.Bind(&playlist); err != nil {
		log.Printf("Error binding playlist data for update: %v", err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Save the updated playlist
	if err := config.DB.Save(&playlist).Error; err != nil {
		log.Printf("Error updating playlist: %v", err)
		return c.JSON(http.StatusInternalServerError, "Failed to update playlist")
	}

	return c.JSON(http.StatusOK, playlist)
}

// DeletePlaylist deletes a playlist by ID
func DeletePlaylist(c echo.Context) error {
	id := c.Param("id")
	var playlist models.Playlist
	if err := config.DB.First(&playlist, id).Error; err != nil {
		log.Printf("Error fetching playlist for deletion with ID %s: %v", id, err)
		return c.JSON(http.StatusNotFound, "Playlist not found")
	}

	// Delete the playlist
	if err := config.DB.Delete(&playlist).Error; err != nil {
		log.Printf("Error deleting playlist: %v", err)
		return c.JSON(http.StatusInternalServerError, "Failed to delete playlist")
	}

	return c.NoContent(http.StatusNoContent)
}

// GetPlaylistsForUser fetches playlists for a specific user
func GetPlaylistsForUser(c echo.Context) error {
	userID := c.Get("uid")
	var playlists []models.Playlist

	// Preload the tracks for each playlist
	if err := config.DB.Where("user_id = ?", userID).Preload("Tracks").Find(&playlists).Error; err != nil {
		log.Printf("Error fetching playlists for user ID %s: %v", userID, err)
		return c.JSON(http.StatusInternalServerError, "Failed to fetch playlists for user")
	}

	return c.JSON(http.StatusOK, playlists)
}


// AddPlaylistForUser adds a playlist for the current authenticated user
func AddPlaylistForUser(c echo.Context) error {
    // Get the current user from the context (assuming the user is authenticated)
    var playlist models.Playlist
    if err := c.Bind(&playlist); err != nil {
        log.Printf("Error binding playlist data: %v", err)
        return c.JSON(http.StatusBadRequest, "Invalid input data")
    }

    // Log the received playlist data for debugging
    log.Printf("Received playlist data: %+v", playlist)

    // Save the playlist to the database
    if err := config.DB.Create(&playlist).Error; err != nil {
        log.Printf("Error creating playlist: %v", err)
        return c.JSON(http.StatusInternalServerError, "Failed to create playlist")
    }

    return c.JSON(http.StatusCreated, playlist)
}



// GetPlaylistByUserID returns a specific user‑owned playlist with its tracks.
func GetPlaylistByUserID(c echo.Context) error {
	uid := c.Param("id")

	var playlists []models.Playlist
	if err := config.DB.
		Where("user_id = ?", uid).
		Preload("Tracks").
		Find(&playlists).Error; err != nil {
		log.Printf("Error fetching playlists for user %s: %v", uid, err)
		return c.JSON(http.StatusInternalServerError, "Failed to fetch playlists")
	}

	return c.JSON(http.StatusOK, playlists)
}

//Add tracks to a playlist

// Define a struct for the request body
type AddTracksRequest struct {
	TrackIDs []uint `json:"track_ids"`
}

func AddTracksToPlaylist(c echo.Context) error {
	// Get the playlist ID from the route parameter
	playlistID := c.Param("id")
	log.Println("Extracted playlist ID:", playlistID)

	// Convert the playlist ID to uint (if necessary)
	playlistIDInt, err := strconv.Atoi(playlistID)
	if err != nil {
		log.Printf("Error parsing playlist ID: %v", err)
		return c.JSON(http.StatusBadRequest, "Invalid playlist ID")
	}

	// Bind the request body to the struct
	var requestBody AddTracksRequest
	if err := c.Bind(&requestBody); err != nil {
		log.Printf("Error binding track IDs: %v", err)
		return c.JSON(http.StatusBadRequest, "Invalid request format")
	}

	// Check if the playlist exists in the database
	var playlist models.Playlist
	if err := config.DB.First(&playlist, playlistIDInt).Error; err != nil {
		log.Printf("Error fetching playlist: %v", err)
		return c.JSON(http.StatusNotFound, "Playlist not found")
	}

	fmt.Printf("Playlist found: %+v\n", playlist)

	// Add each track to the playlist_tracks table
	for _, trackID := range requestBody.TrackIDs {
		// Insert into playlist_tracks table
		if err := config.DB.Create(&models.PlaylistTrack{
			PlaylistID: uint(playlistIDInt), // Correctly use PlaylistID as uint
			TrackID:    trackID,
		}).Error; err != nil {
			log.Printf("Error adding track %d to playlist %d: %v", trackID, playlistIDInt, err)
			return c.JSON(http.StatusInternalServerError, "Failed to add tracks")
		}
	}

	// Return success message
	return c.JSON(http.StatusOK, "Tracks added successfully")
}

// GetTracksByPlaylistID fetches the tracks associated with a specific playlist
func GetTracksByPlaylistID(c echo.Context) error {
    playlistID := c.Param("id")
	
    var playlist models.Playlist
    // Fetch the playlist and preload its tracks
    if err := config.DB.Preload("Tracks").First(&playlist, playlistID).Error; err != nil {
        log.Printf("Error fetching playlist with ID %s: %v", playlistID, err)
        return c.JSON(http.StatusNotFound, "Playlist not found")
    }

    // Return the tracks associated with the playlist
    return c.JSON(http.StatusOK, playlist.Tracks)