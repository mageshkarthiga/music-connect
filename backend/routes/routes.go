package routes

import (
	"backend/controllers"
	"backend/services/chat"
	services "backend/services/spotify"

	"backend/services/maps"

	"github.com/labstack/echo/v4"
)

// RegisterRoutes sets up API endpoints for users, tracks, events and third party services.

func RegisterRoutes(e *echo.Echo, wsServer *chat.WsServer) {

	// projectID := "music-connect-608f6" // Replace with your actual project ID

	// User Routes
	e.GET("/users", controllers.GetUsers) // Fetch all users

	// e.GET("/users/:UserID", controllers.GetUser)          // Fetch a user by ID
	e.GET("/users/firebase/:uid", controllers.GetUserByFirebaseUID) // Fetch a user by Firebase UID
	e.POST("/users", controllers.CreateUser)                        // Create a new user
	e.PUT("/users/:UserID", controllers.UpdateUser)                 // Update an existing user by ID
	e.DELETE("/users/:UserID", controllers.DeleteUser)              // Delete a user by ID
	// e.POST("/auth/login", controllers.RegisterAuthRoutes)
	e.GET("/firebase/:uid", controllers.GetUserByFirebaseUID) // Fetch Firebase UID from token
	// e.GET("/users/:UserID", controllers.GetUser)          // Fetch a user by ID
	e.GET("/users/firebase/:uid", controllers.GetUserByFirebaseUID) // Fetch a user by Firebase UID
	e.POST("/users", controllers.CreateUser)                        // Create a new user
	e.PUT("/users/:UserID", controllers.UpdateUser)                 // Update an existing user by ID
	e.DELETE("/users/:UserID", controllers.DeleteUser)              // Delete a user by ID
	// e.POST("/auth/login", controllers.RegisterAuthRoutes)
	e.GET("/firebase/:uid", controllers.GetUserByFirebaseUID) // Fetch Firebase UID from token

    e.POST("/users/friends", controllers.AddFriend)
    e.GET("/users/:id/friends", controllers.GetFriends)

    //Friend Routes

    e.POST("/friend/:friend_id/request", controllers.SendFriendRequest)
    e.POST("/friend/:friend_id/accept", controllers.AcceptFriendRequest)
    e.POST("/friend/:friend_id/reject", controllers.RejectFriendRequest)
    e.POST("/friend/:friend_id/remove", controllers.RemoveFriend)
    e.GET("/friends", controllers.GetFriends)


    e.POST("/users/friends", controllers.AddFriend)
    e.GET("/users/:id/friends", controllers.GetFriends)

    //Friend Routes

    e.POST("/friend/:friend_id/request", controllers.SendFriendRequest)
    e.POST("/friend/:friend_id/accept", controllers.AcceptFriendRequest)
    e.POST("/friend/:friend_id/reject", controllers.RejectFriendRequest)
    e.POST("/friend/:friend_id/remove", controllers.RemoveFriend)
    e.GET("/friends", controllers.GetFriends)


	// Track Routes
	e.GET("/tracks", controllers.GetTracks)          // Fetch all tracks
	e.GET("/tracks/:id", controllers.GetTrackByID)   // Fetch a track by ID (Updated function name)
	e.POST("/tracks", controllers.CreateTrack)       // Create a new track
    
	e.PUT("/tracks/:id", controllers.UpdateTrack)    // Update an existing track by ID
	e.DELETE("/tracks/:id", controllers.DeleteTrack) // Delete a track by ID
    // e.GET("/tracks/:id/playlist", controllers.GetPlaylistByTrackID) // Fetch playlist by track ID
    e.GET("/likedTracks", controllers.GetLikedTracks) // Fetch liked tracks for a user
    e.POST("/likeTrack", controllers.LikeTrack) // Post request to like a track
    e.DELETE("/likeTrack", controllers.UnlikeTrack) // Delete request to unlike a track

	// Event Routes
	e.GET("/events", controllers.GetEvents)                      // Fetch all events
	e.POST("/events", controllers.CreateEvent)                   // Create a new event
	e.PUT("/events/:id", controllers.UpdateEvent)                // Update an existing event by ID
	e.DELETE("/events/:id", controllers.DeleteEvent)             // Delete an event by ID
	e.POST("/users/:userId/events", controllers.AddEventForUser) // Add an event for a specific user
    e.GET("/events/venues", controllers.GetEventVenues) // Fetch all events with venues
    e.GET("me/likedEvents", controllers.GetLikedEvents)  // Get liked events for a user
    e.POST("/likeEvent", controllers.LikeEvent)            // Post request to like an event
    e.DELETE("/likeEvent", controllers.UnlikeEvent) // Delete request to unlike an event

	//Playlist Routes
	e.GET("/playlists", controllers.GetPlaylists)          // Fetch all playlists
	e.GET("/playlists/:id", controllers.GetPlaylistByID)   // Fetch playlist by ID
	e.POST("/playlists", controllers.CreatePlaylist)       // Create a new playlist
	e.PUT("/playlists/:id", controllers.UpdatePlaylist)    // Update an existing playlist by ID
	e.DELETE("/playlists/:id", controllers.DeletePlaylist) // Delete a playlist by ID
    e.PUT("/playlists/:id", controllers.UpdatePlaylistDetails)
    e.PUT("/playlists/:id/tracks", controllers.ReplacePlaylistTracks)
    e.POST("/playlists/:id/tracks", controllers.AddTracksToPlaylist)


	e.GET("/me/playlists", controllers.GetPlaylistsForUser)           // Fetch playlists for a specific user
	e.POST("/me/playlists", controllers.AddPlaylistForUser)           // Add a playlist for a specific user
	e.GET("/playlists/:id/tracks", controllers.GetTracksByPlaylistID) // Fetch tracks in a playlist by ID
    e.PUT("/me/playlists/:playlistId/tracks", controllers.AddTracksToPlaylist) // Add tracks to a playlist

	// e.GET("/me/playlists", controllers.GetPlaylistsForUser) // Fetch playlists for a specific user
	// e.POST("/me/playlists", controllers.AddPlaylistForUser) // Add a playlist for a specific user

	// Service Routes
	e.GET("/spotify/token", services.GetSpotifyToken)

    e.GET("/calculateSimilarity", controllers.CalculateUserSimilarityHandler)
    e.GET("/userProfile", controllers.GetUserProfileHandler)

	// WebSocket Route
	e.GET("/ws", func(c echo.Context) error {
		chat.ServeWs(wsServer, c.Response().Writer, c.Request())
		return nil
	})
	// Message Retrieval Route
	e.GET("/rooms/:roomName/messages", controllers.GetMessagesForRoom)

	// Google Maps API Proxy Routes
	e.GET("/places/autocomplete", maps.AutocompleteHandler)
	e.GET("/maps", maps.MapsJSHandler)

}
