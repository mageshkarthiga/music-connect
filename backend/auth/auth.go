package auth

import (
	"backend/controllers" 
	"backend/middleware"
	"github.com/labstack/echo/v4"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
	"fmt"
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
e.GET("/users/:id",          middleware.AuthMiddleware(projectID)(controllers.GetUserByUserID))
e.GET("/users",              middleware.AuthMiddleware(projectID)(controllers.GetUsers))
e.GET("/me",                 middleware.AuthMiddleware(projectID)(controllers.GetMe))
e.PUT("/users/:id",          middleware.AuthMiddleware(projectID)(controllers.UpdateUser))
e.DELETE("/users/:id",       middleware.AuthMiddleware(projectID)(controllers.DeleteUser))
e.GET("/users/firebase/:uid",middleware.AuthMiddleware(projectID)(controllers.GetUserByFirebaseUID))

// ---- tracks ----
e.GET("/tracks",             middleware.AuthMiddleware(projectID)(controllers.GetTracks))
e.GET("/me/tracks",          middleware.AuthMiddleware(projectID)(controllers.GetTracksForUser))
e.POST("/me/tracks",         middleware.AuthMiddleware(projectID)(controllers.AddTracksForUser))
e.GET("/users/:id/tracks",   middleware.AuthMiddleware(projectID)(controllers.GetUserTracksByID))

// ---- playlists ----
e.GET("/me/playlists",       middleware.AuthMiddleware(projectID)(controllers.GetPlaylistsForUser))
e.POST("/me/playlists",      middleware.AuthMiddleware(projectID)(controllers.AddPlaylistForUser))
e.GET("/users/:id/playlists",middleware.AuthMiddleware(projectID)(controllers.GetPlaylistByUserID))
e.PUT("/me/playlists/:id/tracks",middleware.AuthMiddleware(projectID) (controllers.AddTracksToPlaylist)) 

// ---- events ----
e.GET("/events",             middleware.AuthMiddleware(projectID)(controllers.GetEvents))
e.POST("/events",            middleware.AuthMiddleware(projectID)(controllers.CreateEvent))
e.PUT("/events/:id",         middleware.AuthMiddleware(projectID)(controllers.UpdateEvent))
e.DELETE("/events/:id",      middleware.AuthMiddleware(projectID)(controllers.DeleteEvent))
e.GET("/me/events",          middleware.AuthMiddleware(projectID)(controllers.GetEventsForUser))
e.GET("/users/:id/events",   middleware.AuthMiddleware(projectID)(controllers.GetEventsByUserID))
e.POST("/users/:id/events",  middleware.AuthMiddleware(projectID)(controllers.AddEventForUser))
e.GET("/me/likedEvents",     middleware.AuthMiddleware(projectID)(controllers.GetLikedEvents))
e.POST("/likeEvent",         middleware.AuthMiddleware(projectID)(controllers.LikeEvent))
e.DELETE("/likeEvent",      middleware.AuthMiddleware(projectID)(controllers.UnlikeEvent))

// ---- recommendations ----
e.GET("/tracks/recommendations", middleware.AuthMiddleware(projectID)(controllers.GetTrackRecommendation))
}
