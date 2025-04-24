package main

import (
	"chat-service/firebase"
	"chat-service/routes"
	"chat-service/server"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"os"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables.")
	}

	// Initialize Firebase
	firebase.InitFirebase()

	// Initialize Echo
	e := echo.New()

	// CORS middleware
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:8080", "http://localhost:8002", "https://music-connect-three.vercel.app","https://music-connect-five.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "x-xsrf-token"},
		AllowCredentials: true,
	}))
	log.Println("✅ CORS middleware applied")

	// Initialize the WebSocket server
	wsServer := server.NewWsServer()
	go wsServer.Run() // Start the WebSocket server in a separate goroutine

	// Register routes
	routes.RegisterRoutes(e, wsServer)

	// Health check route
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// Start the server on the specified port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to 8080 if not set
	}
	log.Printf("🚀 Chat service running on %s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatal("❌ Echo server error:", err)
	}
}

