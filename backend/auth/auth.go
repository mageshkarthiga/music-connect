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

	// e.GET("/users/:id", middleware.AuthMiddleware(projectID)(controllers.GetUser))        // Fetch a user by ID
	e.PUT("/users/:id", middleware.AuthMiddleware(projectID)(controllers.UpdateUser))     // Update an existing user by ID
	e.DELETE("/users/:id", middleware.AuthMiddleware(projectID)(controllers.DeleteUser))  // Delete a user by ID
    e.GET("/firebase/:uid", middleware.AuthMiddleware(projectID)(controllers.GetUserByFirebaseUID)) // Fetch Firebase UID from token

	// Track Routes (Protected)
	e.POST("/tracks", middleware.AuthMiddleware(projectID)(controllers.CreateTrack))       // Create a new track
	e.PUT("/tracks/:id", middleware.AuthMiddleware(projectID)(controllers.UpdateTrack))    // Update an existing track by ID
	e.DELETE("/tracks/:id", middleware.AuthMiddleware(projectID)(controllers.DeleteTrack)) // Delete a track by ID

	// Event Routes (Protected)
	e.POST("/events", middleware.AuthMiddleware(projectID)(controllers.CreateEvent))       // Create a new event
	e.PUT("/events/:id", middleware.AuthMiddleware(projectID)(controllers.UpdateEvent))    // Update an existing event by ID
	e.DELETE("/events/:id", middleware.AuthMiddleware(projectID)(controllers.DeleteEvent)) // Delete an event by ID
	e.GET("/users/:userId/events", middleware.AuthMiddleware(projectID) (controllers.GetEventsForUser)) // Fetch events for a specific user
    e.POST("/users/:userId/events", middleware.AuthMiddleware(projectID) (controllers.AddEventForUser)) // Add an event for a specific user

}