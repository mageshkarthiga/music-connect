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
	"golang.org/x/oauth2"

)

// Global Spotify Auth Config
var authConfig *clientcredentials.Config

// SpotifyAuth initializes the Spotify authentication configuration
func SpotifyAuth() {
	if os.Getenv("ENV") != "prod" {
		err := godotenv.Load()
		if err != nil {
			fmt.Println("Error loading .env file:", err)
		}
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

// GetSpotifyToken is the Echo handler for /spotify/token
func GetSpotifyToken(c echo.Context) error {
	token, err := GetSpotifyTokenRaw()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get token"})
	}
	return c.JSON(http.StatusOK, token)
}

// GetSpotifyTokenRaw returns the token directly (for use in main.go etc.)
func GetSpotifyTokenRaw() (*oauth2.Token, error) {

	if authConfig == nil {
		SpotifyAuth()
	}
	if authConfig == nil {
		return nil, fmt.Errorf("Spotify auth config not initialized")
	}
	return authConfig.Token(context.Background())
}
