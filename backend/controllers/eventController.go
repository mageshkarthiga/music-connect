package controllers

import (
    "backend/config"
    "backend/models"
    "github.com/labstack/echo/v4"
    "net/http"
    "log"

)

// GetEvents fetches all events
func GetEvents(c echo.Context) error {
    var events []models.Event
    if err := config.DB.Find(&events).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, "Failed to fetch events")
    }
    return c.JSON(http.StatusOK, events)
}

func GetEventVenues(c echo.Context) error {
    var events []models.Event

    if err := config.DB.
        Preload("Venues"). // this loads the many-to-many relation
        Find(&events).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, "Failed to fetch event venues")
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


    //add events for a user
    func AddEventForUser(c echo.Context) error {
        userID := c.Param("id")
        if userID == "" {
            return c.JSON(http.StatusBadRequest, "UserID is missing in the URL")
        }
    
        // Check if the database connection is established
        if config.DB == nil {
            return c.JSON(http.StatusInternalServerError, "Database connection not established")
        }
    
        // Fetch the user from the database to check if the user exists
        var user models.User
        if err := config.DB.First(&user, "user_id = ?", userID).Error; err != nil {
            return c.JSON(http.StatusNotFound, "User not found")
        }
    
        // Bind the request body to the event struct
        var event models.Event
        if err := c.Bind(&event); err != nil {
            return c.JSON(http.StatusBadRequest, err.Error())
        }
    
        // Associate the event with the user
        if err := config.DB.Model(&user).Association("Events").Append(&event); err != nil {
            return c.JSON(http.StatusInternalServerError, "Failed to add event for user")
        }
    
        // Return the updated event as JSON
        return c.JSON(http.StatusOK, event)

    }

    // func GetEventsForUser(c echo.Context) error {
    //     // Get the user ID from the URL parameters
    //     userID := c.Param("userId")
    //     if userID == "" {
    //         return c.JSON(http.StatusBadRequest, "UserID is missing in the URL")
    //     }
    
    //     // Check if the database connection is established
    //     if config.DB == nil {
    //         return c.JSON(http.StatusInternalServerError, "Database connection not established")
    //     }
    
    //     // Fetch the user from the database to check if the user exists
    //     var user models.User
    //     if err := config.DB.First(&user, "user_id = ?", userID).Error; err != nil {
    //         return c.JSON(http.StatusNotFound, "User not found")
    //     }
    
    //     // Fetch all events for the user using the many-to-many relationship
    //     var events []models.Event
    //     if err := config.DB.Model(&user).Association("Events").Find(&events); err != nil {
    //         return c.JSON(http.StatusInternalServerError, "Failed to fetch events")
    //     }
    
    //     // If no events are found, return a message
    //     if len(events) == 0 {
    //         return c.JSON(http.StatusOK, "No events found for this user")
    //     }
    
    //     // Return the events as JSON
    //     return c.JSON(http.StatusOK, events)
    // }

    func GetEventsForUser(c echo.Context) error {
        // Get the user ID from the context (set by the AuthMiddleware)
        userID := c.Get("uid")
        log.Printf("User ID from context: %v\n", userID)
        // userID, ok := uidInterface.(string)

    
        // Check if the database connection is established
        if config.DB == nil {
            return c.JSON(http.StatusInternalServerError, "Database connection not established")
        }
    
        // Fetch the user from the database to check if the user exists
        var user models.User
        if err := config.DB.First(&user, "user_id = ?", userID).Error; err != nil {
            return c.JSON(http.StatusNotFound, "User not found")
        }
    
        // Fetch all events for the user using the many-to-many relationship
        var events []models.Event
        if err := config.DB.Model(&user).Association("Events").Find(&events); err != nil {
            return c.JSON(http.StatusInternalServerError, "Failed to fetch events")
        }
    
        // If no events are found, return a message
        if len(events) == 0 {
            return c.JSON(http.StatusOK, "No events found for this user")
        }
    
        // Return the events as JSON
        return c.JSON(http.StatusOK, events)
    }

    // GetEventsByUserID returns every event tied to the given user ID.
func GetEventsByUserID(c echo.Context) error {
	uid := c.Param("id")

	var events []models.Event
	if err := config.DB.
		Joins("JOIN user_events ue ON ue.event_id = events.event_id").
		Where("ue.user_id = ?", uid).
		Find(&events).Error; err != nil {
		log.Printf("Error fetching events for user %s: %v", uid, err)
		return c.JSON(http.StatusInternalServerError, "Failed to fetch events")
	}

	return c.JSON(http.StatusOK, events)
}

func GetLikedEvents(c echo.Context) error {
    // Get userID from the context set by middleware using "uid"
    uid, ok := c.Get("uid").(uint)
    if !ok {
        return c.JSON(http.StatusUnauthorized, "Missing or invalid user ID in context")
    }

    // Get events from user_events table
    var events []models.Event
    if err := config.DB.
        Joins("JOIN user_events ue ON ue.event_id = events.event_id").
        Where("ue.user_id = ?", uid). // Use `uid` directly here
        Find(&events).Error; err != nil {
        log.Printf("Error fetching events for user %d: %v", uid, err)
        return c.JSON(http.StatusInternalServerError, "Failed to fetch events")
    }

    return c.JSON(http.StatusOK, events)
}

func LikeEvent(c echo.Context) error {
    uid, ok := c.Get("uid").(uint)
    if !ok {
        return c.JSON(http.StatusUnauthorized, "User ID is required")
    }

    var req struct {
        EventID uint `json:"event_id"`
    }

    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, "Invalid request body")
    }

    var existingLike models.UserEvent
    if err := config.DB.
        Where("user_id = ? AND event_id = ?", uid, req.EventID).
        First(&existingLike).Error; err == nil {
        return c.JSON(http.StatusBadRequest, "Event already liked")
    }

    newLike := models.UserEvent{
        UserID:  uid,
        EventID: req.EventID,
    }

    if err := config.DB.Create(&newLike).Error; err != nil {
        log.Printf("Error liking event %d by user %d: %v", req.EventID, uid, err)
        return c.JSON(http.StatusInternalServerError, "Failed to like event")
    }

    return c.JSON(http.StatusOK, "Event liked successfully")
}

func UnlikeEvent(c echo.Context) error {
    uid, ok := c.Get("uid").(uint)
    if !ok {
        return c.JSON(http.StatusUnauthorized, "User ID is required")
    }

    var req struct {
        EventID uint `json:"event_id"`
    }

    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, "Invalid request body")
    }

    if err := config.DB.
        Where("user_id = ? AND event_id = ?", uid, req.EventID).
        Delete(&models.UserEvent{}).Error; err != nil {
        log.Printf("Error unliking event %d by user %d: %v", req.EventID, uid, err)
        return c.JSON(http.StatusInternalServerError, "Failed to unlike event")
    }

    return c.JSON(http.StatusOK, "Event unliked successfully")
}

