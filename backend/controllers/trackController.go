package controllers

import (
	"backend/config"
	"backend/models"
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
