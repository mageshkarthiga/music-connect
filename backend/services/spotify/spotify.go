package spotify

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

var (
    authConfig *clientcredentials.Config
)

// InitSpotify initializes the Spotify client and registers the /token route
func InitSpotify(e *echo.Echo) error {
    // Load environment variables
    err := godotenv.Load("../../.env")
    if err != nil {
        fmt.Println("Error loading .env file")
        return err
    }

    clientId := os.Getenv("SPOTIFY_CLIENT_ID")
    clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
    if clientId == "" || clientSecret == "" {
        fmt.Println("Spotify client ID or client secret not found")
        return fmt.Errorf("missing Spotify client credentials")
    }

    // Configure Spotify client credentials
    authConfig = &clientcredentials.Config{
        ClientID:     clientId,
        ClientSecret: clientSecret,
        TokenURL:     spotify.TokenURL,
        Scopes:       []string{"user-read-playback-state", "user-modify-playback-state", "streaming"},
    }

    // Register the /token route
    e.GET("/spotify/token", getTokenHandler)

    fmt.Println("Spotify service initialized successfully")
    return nil
}

// getTokenHandler handles requests to /spotify/token
func getTokenHandler(c echo.Context) error {
    // Get the Spotify token
    token, err := authConfig.Token(context.Background())
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "Failed to get token",
        })
    }

    // Return the token as JSON
    return c.JSON(http.StatusOK, token)
}