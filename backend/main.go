package main

import (
	"log"
	"backend/config"
	"backend/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"backend/services/spotify" 
)

func main() {
	// Initialize the database connection
	err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	} else {
		log.Println("Database connection initialized successfully! ✅")
	}
	services.SpotifyAuth()
	token, err := services.GetSpotifyTokenRaw()
	if err != nil {
		log.Fatalf("❌ Failed to get Spotify token: %v", err)
	}
	log.Println("🎧 Spotify token retrieved successfully!")
	log.Printf("🔐 Access token: %s", token.AccessToken)


	// Initialize Echo
	e := echo.New()

	// Register routes
	routes.RegisterRoutes(e)
	log.Println("Routes registered successfully ✅")

	// Middleware: CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // Allow all origins (change it to more restrictive rules as needed)
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))
	log.Println("CORS middleware applied ✅")

	// Define a simple route to check if the server is running
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Database connection is live and the server is running on port 8080! ✅")
	})

	// Start the server on port 8080
	log.Printf("Starting server on port 8080... 🚀")
	err = e.Start(":8080")
	if err != nil {
		log.Fatal("Error starting server ⚠️", err)
	} else {
		log.Println("Server started successfully on port 8080 ✅")
	}
}
