package main

import (
	"context"
	"fmt"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

func main() {
	// to authenticate and get token
	authConfig := &clientcredentials.Config{
		ClientID:     "f2726120b049453a9af9dbf0551add00",
		ClientSecret: "47f81105840d46c4b9938e33cb976a3d",
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
