package auth

import (
	"backend/controllers"
	"backend/middleware"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func RegisterAuthRoutes(e *echo.Echo, projectID string) {
	e.POST("/auth/login", func(c echo.Context) error {

		projectID := "music-connect-608f6"

		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, "Missing Authorization header")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, "Invalid Authorization format")
		}

		tokenStr := parts[1]
		fmt.Println("Token String: ", tokenStr)

		token, err := jwt.Parse(tokenStr, middleware.GetJWKS().Keyfunc)
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, fmt.Sprintf("Invalid token: %v", err))
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, "Invalid token claims")
		}

		// Debug: Print the audience claim to check its value
		fmt.Printf("Token Audience: %s\n", claims["aud"])
		fmt.Printf("projectID: %s\n", projectID)
		fmt.Printf("Token Issuer: %s\n", claims["iss"])
		fmt.Printf("Token UID: %s\n", claims["user_id"])

		aud, ok := claims["aud"].(string)
		if !ok || aud != projectID {
			return c.JSON(http.StatusUnauthorized, "Invalid audience")
		}

		uid, ok := claims["user_id"].(string)
		if !ok {
			return c.JSON(http.StatusUnauthorized, "UID missing in token")
		}

		// Set cookie
		if err := SetAuthCookie(c, tokenStr); err != nil {
			return c.JSON(http.StatusInternalServerError, "Error setting cookie")
		}

		return c.JSON(http.StatusOK, map[string]string{
			"message": "Login successful",
			"uid":     uid,
		})
	})
	// ---- user ----
	e.GET("/users/:id", middleware.AuthMiddleware(projectID)(controllers.GetUserByUserID))
	e.GET("/users", middleware.AuthMiddleware(projectID)(controllers.GetUsers))
	e.GET("/me", middleware.AuthMiddleware(projectID)(controllers.GetMe))
	e.PUT("/users/:id", middleware.AuthMiddleware(projectID)(controllers.UpdateUser))
	e.DELETE("/users/:id", middleware.AuthMiddleware(projectID)(controllers.DeleteUser))
	e.GET("/users/firebase/:uid", middleware.AuthMiddleware(projectID)(controllers.GetUserByFirebaseUID))

	// ---- friends ----

	e.POST("/friends/:friend_id/request", middleware.AuthMiddleware(projectID)(controllers.SendFriendRequest))
	e.POST("/friends/:friend_id/accept", middleware.AuthMiddleware(projectID)(controllers.AcceptFriendRequest))
	e.POST("/friends/:friend_id/reject", middleware.AuthMiddleware(projectID)(controllers.RejectFriendRequest))
	e.POST("/friends/:friend_id/remove", middleware.AuthMiddleware(projectID)(controllers.RemoveFriend))
	e.GET("/friends", middleware.AuthMiddleware(projectID)(controllers.GetFriends))
	e.GET("/friendship/:friend_id/status", middleware.AuthMiddleware(projectID)(controllers.GetFriendshipStatus))
	e.GET("/friends/pending", middleware.AuthMiddleware(projectID)(controllers.GetPendingFriendRequests))
	e.POST("/users/friends", middleware.AuthMiddleware(projectID)(controllers.AddFriend))
	e.GET("/users/:id/friends", middleware.AuthMiddleware(projectID)(controllers.GetFriends))

	// ---- tracks ----
	e.GET("/tracks", middleware.AuthMiddleware(projectID)(controllers.GetTracks))
	e.GET("/me/tracks", middleware.AuthMiddleware(projectID)(controllers.GetTracksForUser))
	e.GET("/me/favtracks", middleware.AuthMiddleware(projectID)(controllers.GetFavTracksForUser))
	e.POST("/me/tracks", middleware.AuthMiddleware(projectID)(controllers.AddTracksForUser))

	e.GET("/users/:id/tracks", middleware.AuthMiddleware(projectID)(controllers.GetUserTracksByID))
	e.GET("/users/:id/favtracks", middleware.AuthMiddleware(projectID)(controllers.GetFavUserTracksByID))
	e.GET("/likedTracks", middleware.AuthMiddleware(projectID)(controllers.GetLikedTracks))
	e.PUT("/likeTrack/:track_id", middleware.AuthMiddleware(projectID)(controllers.LikeTrack))   // PUT request to like a track
	e.PUT("/unlikeTrack/:track_id", middleware.AuthMiddleware(projectID)(controllers.UnlikeTrack)) // PUT request to unlike a track
	
	e.GET("/tracks/top", middleware.AuthMiddleware(projectID)(controllers.GetTopPlayedTracks))
	e.PUT("/tracks/:track_id/increment", middleware.AuthMiddleware(projectID)(controllers.IncrementTrackPlayCount))

	// ---- playlists ----
	e.GET("/me/playlists", middleware.AuthMiddleware(projectID)(controllers.GetPlaylistsForUser))
	e.POST("/me/playlists", middleware.AuthMiddleware(projectID)(controllers.AddPlaylistForUser))
	e.GET("/users/:id/playlists", middleware.AuthMiddleware(projectID)(controllers.GetPlaylistByUserID))
	e.POST("/playlists/:id/tracks", middleware.AuthMiddleware(projectID)(controllers.AddTracksToPlaylist))

	// ---- events ----
	e.GET("/events", middleware.AuthMiddleware(projectID)(controllers.GetEvents))
	e.POST("/events", middleware.AuthMiddleware(projectID)(controllers.CreateEvent))
	e.PUT("/events/:id", middleware.AuthMiddleware(projectID)(controllers.UpdateEvent))
	e.DELETE("/events/:id", middleware.AuthMiddleware(projectID)(controllers.DeleteEvent))
	e.GET("/me/events", middleware.AuthMiddleware(projectID)(controllers.GetEventsForUser))
	e.GET("/me/favevents", middleware.AuthMiddleware(projectID)(controllers.GetFavEventsForUser))
	e.GET("/users/:id/events", middleware.AuthMiddleware(projectID)(controllers.GetEventsByUserID))
	e.GET("/users/:id/favevents", middleware.AuthMiddleware(projectID)(controllers.GetFavEventsByUserID))
	e.POST("/users/:id/events", middleware.AuthMiddleware(projectID)(controllers.AddEventForUser))
	e.GET("/me/likedEvents", middleware.AuthMiddleware(projectID)(controllers.GetLikedEvents))
	e.POST("/likeEvent", middleware.AuthMiddleware(projectID)(controllers.LikeEvent))
	e.DELETE("/likeEvent", middleware.AuthMiddleware(projectID)(controllers.UnlikeEvent))

	// ---- recommendations ----
	e.GET("/tracks/recommendations", middleware.AuthMiddleware(projectID)(controllers.GetTrackRecommendation))

	// ---- chat ----
	e.GET("/rooms/:roomName/messages", middleware.AuthMiddleware(projectID)(controllers.GetMessagesForRoom))
	e.GET("/users/:userID/chat-history", middleware.AuthMiddleware(projectID)(controllers.GetUsersWithChatHistory))

}
