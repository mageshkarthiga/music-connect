package routes

import (
	"backend/controllers"
	"backend/services/spotify"
	"github.com/labstack/echo/v4"
    "backend/services/chat"
)

// RegisterRoutes sets up API endpoints for users, tracks, events and third party services.
func RegisterRoutes(e *echo.Echo) {
    // User Routes
    e.GET("/users", controllers.GetUsers)             // Fetch all users
    e.GET("/users/:id", controllers.GetUser)          // Fetch a user by ID
    e.GET("/users/firebase/:uid", controllers.GetUserByFirebaseUID) // Fetch a user by Firebase UID
    e.POST("/users", controllers.CreateUser)          // Create a new user
    e.PUT("/users/:id", controllers.UpdateUser)       // Update an existing user by ID
    e.DELETE("/users/:id", controllers.DeleteUser)    // Delete a user by ID

    // Track Routes
    e.GET("/tracks", controllers.GetTracks)           // Fetch all tracks
    e.GET("/tracks/:id", controllers.GetTrackByID)        // Fetch a track by ID (Updated function name)
    e.POST("/tracks", controllers.CreateTrack)        // Create a new track
    e.PUT("/tracks/:id", controllers.UpdateTrack)     // Update an existing track by ID
    e.DELETE("/tracks/:id", controllers.DeleteTrack)  // Delete a track by ID

    // Event Routes
    e.GET("/events", controllers.GetEvents)           // Fetch all events
    e.GET("/events/:id", controllers.GetEventByID)        // Fetch an event by ID (Updated function name)
    e.POST("/events", controllers.CreateEvent)        // Create a new event
    e.PUT("/events/:id", controllers.UpdateEvent)     // Update an existing event by ID
    e.DELETE("/events/:id", controllers.DeleteEvent)  // Delete an event by ID

    // Service Routes
    e.GET("/spotify/token", services.GetSpotifyToken)
    // WebSocket Route
    wsServer := chat.NewWsServer()
    go wsServer.Run()
    e.GET("/ws", func(c echo.Context) error {
        chat.ServeWs(wsServer, c.Response().Writer, c.Request())
        return nil
    })
}
