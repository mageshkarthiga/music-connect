package main

import (
	"backend/config"
	"backend/middleware" // Make sure this is the correct path
	"backend/routes"
	"backend/services/spotify"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load Firebase project ID from environment or hardcode for now
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	if projectID == "" {
		projectID = "your-firebase-project-id" // TODO: Replace with your real Firebase project ID
	}

	// Initialize JWKS (needed before parsing Firebase JWTs)
	if err := middleware.InitJWKS(); err != nil {
		log.Fatalf("❌ Failed to initialize JWKS: %v", err)
	}
	log.Println("🔑 JWKS initialized successfully!")

	// Initialize DB
	if err := config.InitDB(); err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}
	log.Println("✅ Database connection initialized!")

	// Authenticate with Spotify
	services.SpotifyAuth()
	token, err := services.GetSpotifyTokenRaw()
	if err != nil {
		log.Fatalf("❌ Failed to get Spotify token: %v", err)
	}
	log.Println("🎧 Spotify token retrieved successfully!")
	log.Printf("🔐 Access token: %s", token.AccessToken)

	// Initialize Echo
	e := echo.New()

	// CORS middleware
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))
	log.Println("✅ CORS middleware applied")

	// Register routes
	routes.RegisterRoutes(e)
	log.Println("✅ Routes registered")

	// Health check
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "✅ Server is live on port 8080!")
	})

	// Example secure route with JWT
	e.GET("/secure", func(c echo.Context) error {
		uid := c.Get("uid").(string)
		return c.JSON(200, map[string]string{
			"message": "Authenticated ✅",
			"uid":     uid,
		})
	}, middleware.AuthMiddleware(projectID))

	// Start server
	log.Println("🚀 Starting server on port 8080...")
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("❌ Error starting server: %v", err)
	}
}
