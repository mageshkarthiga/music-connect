package main

import (
	"backend/auth"
	"backend/config"
	"backend/middleware" 
	"backend/models"
	"backend/routes"
	"backend/services/spotify"
	"fmt"
	"log"
	"os"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load Firebase project ID from environment or hardcode for now
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	print("Project ID: ", projectID)

	// Initialize JWKS (needed before parsing Firebase JWTs)
	if err := middleware.InitJWKS(); err != nil {
		log.Fatalf("❌ Failed to initialize JWKS: %v", err)
	}
	log.Println("🔑 JWKS initialized successfully!")

	// Initialize DB
	if err := config.InitDB(); err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
		fmt.Println("Database connection status:", config.DB)
	}
	if config.DB == nil {
		log.Fatal("❌ Database connection is nil!")
	} else {
		log.Println("✅ Database connection initialized!")
	}

	// Run migrations
	if err := config.DB.AutoMigrate(
		&models.Event{},
		&models.UserEvent{},
		&models.Playlist{},
		&models.Track{},
		&models.PlaylistTrack{},
		&models.TrackArtist{},
	); err != nil {
		log.Fatal("❌ Failed to run migrations: ", err)
	} else {
		log.Println("✅ Migrations completed successfully!")
	}

	// Initialize Spotify services
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
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:8080", "https://music-connect-555448022527.us-central1.run.app", "https://music-connect-three.vercel.app","https://music-connect-five.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "x-xsrf-token"},
		AllowCredentials: true,
	}))
	log.Println("✅ CORS middleware applied")

	// Register routes 
	routes.RegisterRoutes(e)
	log.Println("✅ Routes registered")

	// Health check
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "✅ Server is live on port 8080!")
	})

	// Example secure route with JWT authentication
	e.GET("/secure", func(c echo.Context) error {
		uid := c.Get("uid").(string)
		return c.JSON(200, map[string]string{
			"message": "Authenticated ✅",
			"uid":     uid,
		})
	}, middleware.AuthMiddleware(projectID))

	auth.RegisterAuthRoutes(e, projectID)
	log.Println("✅ Auth routes registered")

	// Start server
	log.Println("🚀 Starting server on port 8080...")
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("❌ Error starting server: %v", err)
	}
}
