package controllers

import (
    "backend/config"
    "backend/models"
    "github.com/labstack/echo/v4"
    "net/http"
)

// GetEvents fetches all events
func GetEvents(c echo.Context) error {
    var events []models.Event
    if err := config.DB.Find(&events).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, "Failed to fetch events")
    }
    return c.JSON(http.StatusOK, events)
}

// GetEventByID fetches a single event by ID
func GetEventByID(c echo.Context) error {
    id := c.Param("id")
    var event models.Event
    if err := config.DB.First(&event, id).Error; err != nil {
        return c.JSON(http.StatusNotFound, "Event not found")
    }
    return c.JSON(http.StatusOK, event)
}

// CreateEvent creates a new event
func CreateEvent(c echo.Context) error {
    var event models.Event
    if err := c.Bind(&event); err != nil {
        return c.JSON(http.StatusBadRequest, err.Error())
    }
    if err := config.DB.Create(&event).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, "Failed to create event")
    }
    return c.JSON(http.StatusCreated, event)
}

// UpdateEvent updates an existing event by ID
func UpdateEvent(c echo.Context) error {
    id := c.Param("id")
    var event models.Event
    if err := config.DB.First(&event, id).Error; err != nil {
        return c.JSON(http.StatusNotFound, "Event not found")
    }

    if err := c.Bind(&event); err != nil {
        return c.JSON(http.StatusBadRequest, err.Error())
    }

    if err := config.DB.Save(&event).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, "Failed to update event")
    }
    return c.JSON(http.StatusOK, event)
}

// DeleteEvent deletes an event by ID
func DeleteEvent(c echo.Context) error {
    id := c.Param("id")
    var event models.Event
    if err := config.DB.First(&event, id).Error; err != nil {
        return c.JSON(http.StatusNotFound, "Event not found")
    }

    if err := config.DB.Delete(&event).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, "Failed to delete event")
    }
    return c.JSON(http.StatusOK, "Event deleted successfully")
}
