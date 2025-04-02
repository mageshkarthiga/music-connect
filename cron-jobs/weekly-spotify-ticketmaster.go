package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"github.com/joho/godotenv"
)

type SpotifyTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func getSpotifyToken() (string, error) {
	tokenURL := "https://accounts.spotify.com/api/token"
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("client_id", clientID)
	form.Add("client_secret", clientSecret)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("Response:", string(body))

	var tokenResp SpotifyTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}

func callSpotifyAPI(url string) {
	token, _ := getSpotifyToken()

    req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer " + token)
	client := &http.Client{}
	resp, err := client.Do(req)
	
    if err != nil {
        fmt.Printf("Error calling %s: %v\n", url, err)
        return
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)
    fmt.Printf("Response from %s: %s\n", url, body)
}

// func callTicketMasterAPI(url string) {
//     resp, err := http.Get(url)
//     if err != nil {
//         fmt.Printf("Error calling %s: %v\n", url, err)
//         return
//     }
//     defer resp.Body.Close()

//     body, _ := io.ReadAll(resp.Body)
//     fmt.Printf("Response from %s: %s\n", url, body)
// }

func main() {
	if err := godotenv.Load("../.env"); 
	err != nil {
		log.Println("Warning: .env file not found, relying on environment variables")
	}

    callSpotifyAPI("https://api.spotify.com/v1/tracks/11dFghVXANMlKmJXsNCbNl") // to be updated
    // callTicketMasterAPI("https://api.example.com/endpoint2")
}
