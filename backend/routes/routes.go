package routes

import (
	"backend/controllers"
	"backend/services/spotify"
	"github.com/labstack/echo/v4"
    "backend/services/chat"
)

// RegisterRoutes sets up API endpoints for users, tracks, events and third party services.
func RegisterRoutes(e *echo.Echo, wsServer *chat.WsServer) {

    // projectID := "music-connect-608f6" // Replace with your actual project ID

    // User Routes
    e.GET("/users", controllers.GetUsers)             // Fetch all users
    
    // e.GET("/users/:UserID", controllers.GetUser)          // Fetch a user by ID
    e.GET("/users/firebase/:uid", controllers.GetUserByFirebaseUID) // Fetch a user by Firebase UID
    e.POST("/users", controllers.CreateUser)          // Create a new user
    e.PUT("/users/:UserID", controllers.UpdateUser)       // Update an existing user by ID
    e.DELETE("/users/:UserID", controllers.DeleteUser)    // Delete a user by ID
    // e.POST("/auth/login", controllers.RegisterAuthRoutes)
    e.GET("/firebase/:uid", controllers.GetUserByFirebaseUID) // Fetch Firebase UID from token


    // Track Routes
    
    e.GET("/tracks", controllers.GetTracks)           // Fetch all tracks
    e.GET("/tracks/:id", controllers.GetTrackByID)        // Fetch a track by ID (Updated function name)
    e.POST("/tracks", controllers.CreateTrack)        // Create a new track
    e.PUT("/tracks/:id", controllers.UpdateTrack)     // Update an existing track by ID
    e.DELETE("/tracks/:id", controllers.DeleteTrack)  // Delete a track by ID


    // Event Routes
    e.GET("/events", controllers.GetEvents)           // Fetch all events    
    e.POST("/events", controllers.CreateEvent)        // Create a new event
    e.PUT("/events/:id", controllers.UpdateEvent)     // Update an existing event by ID
    e.DELETE("/events/:id", controllers.DeleteEvent)  // Delete an event by ID    
    e.POST("/users/:userId/events", controllers.AddEventForUser) // Add an event for a specific user

    //Playlist Routes
    e.GET("/playlists", controllers.GetPlaylists)           // Fetch all playlists
    e.GET("/playlists/:id", controllers.GetPlaylistByID)        // Fetch a playlist by ID (Updated function name)
    e.POST("/playlists", controllers.CreatePlaylist)        // Create a new playlist
    e.PUT("/playlists/:id", controllers.UpdatePlaylist)     // Update an existing playlist by ID
    e.DELETE("/playlists/:id", controllers.DeletePlaylist)  // Delete a playlist by ID

    e.GET("/me/playlists", controllers.GetPlaylistsForUser) // Fetch playlists for a specific user
    e.POST("/me/playlists", controllers.AddPlaylistForUser) // Add a playlist for a specific user

    // Service Routes
    e.GET("/spotify/token", services.GetSpotifyToken)
    // WebSocket Route
    e.GET("/ws", func(c echo.Context) error {
        chat.ServeWs(wsServer, c.Response().Writer, c.Request())
        return nil
    })
    // Message Retrieval Route
    e.GET("/rooms/:roomName/messages", controllers.GetMessagesForRoom)
    
}

