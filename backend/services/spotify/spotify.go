package services

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

// Global Spotify Auth Config
var authConfig *clientcredentials.Config

// SpotifyAuth initializes the Spotify authentication configuration
func SpotifyAuth() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	clientId := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	if clientId == "" || clientSecret == "" {
		fmt.Println("Spotify credentials missing. Ensure .env contains SPOTIFY_CLIENT_ID and SPOTIFY_CLIENT_SECRET")
		return
	}

	authConfig = &clientcredentials.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		TokenURL:     spotify.TokenURL,
		Scopes:       []string{"user-read-playback-state", "user-modify-playback-state", "streaming"},
	}

	fmt.Println("✅ Spotify authentication initialized successfully!")
}


// GetSpotifyToken is the handler function for the /spotify/token route
func GetSpotifyToken(c echo.Context) error {
    SpotifyAuth() // Ensure SpotifyAuth is called to initialize authConfig
	if authConfig == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Spotify authentication is not initialized"})
	}

	token, err := authConfig.Token(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get token"})
	}
	return c.JSON(http.StatusOK, token)
}

