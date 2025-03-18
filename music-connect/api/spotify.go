package main

import (
	"context"
	"fmt"
	"os"
	"github.com/joho/godotenv"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	clientId := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	if clientId == "" || clientSecret == "" {
		fmt.Println("Spotify client id or client secret not found")
		return
	}
	// to authenticate and get token
	authConfig := &clientcredentials.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		TokenURL:     spotify.TokenURL,
	}
	accessToken, err := authConfig.Token(context.Background())
	if err != nil {
		fmt.Println("error retrieving access token: %v", err)
	}

	// connect to Spotify API
	spotifyClient := spotify.Authenticator{}.NewClient(accessToken)
	artistID := spotify.ID("4zCH9qm4R2DADamUHMCa6O")
	artist, err := spotifyClient.GetArtist(artistID)

	if err != nil {
        fmt.Printf("error retrieving artist: %v\n", err)
        return
    }


	fmt.Println("artist id:", artist.ID)
	fmt.Println("artist name:", artist.Name)
}
