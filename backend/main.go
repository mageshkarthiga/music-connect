package main

import (
	"log"
	"backend/config"
	"backend/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)


func main() {
	// Initialize the database connection
	err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize Echo
	e := echo.New()

	// Register routes
	routes.RegisterRoutes(e)

	// Middleware: CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // Allow all origins (change it to more restrictive rules as needed)
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	// Define a simple route to check if the server is running
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Database connection is live and the server is running on port 1323! ✅")
	})

	// Start the server on port 1323
	err = e.Start(":1323")
	if err != nil {
		log.Fatal("Error starting server ⚠️", err)
	}
}
