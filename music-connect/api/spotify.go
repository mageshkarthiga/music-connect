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
		ClientID:     "",
		ClientSecret: "",
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
