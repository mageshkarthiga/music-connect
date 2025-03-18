package main

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "os"

    "github.com/joho/godotenv"
    "github.com/zmb3/spotify"
    "golang.org/x/oauth2/clientcredentials"
)

var (
    authConfig *clientcredentials.Config
)

func main() {
    err := godotenv.Load("../.env")
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

    authConfig = &clientcredentials.Config{
        ClientID:     clientId,
        ClientSecret: clientSecret,
        TokenURL:     spotify.TokenURL,
        Scopes:       []string{"user-read-playback-state", "user-modify-playback-state", "streaming"},
    }

    http.HandleFunc("/token", getTokenHandler)
    fmt.Println("Server is running on port 8080")
    http.ListenAndServe(":8080", nil)
}

func getTokenHandler(w http.ResponseWriter, r *http.Request) {
    enableCors(&w)
    token, err := authConfig.Token(context.Background())
    if err != nil {
        http.Error(w, "Failed to get token", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(token)
}

func enableCors(w *http.ResponseWriter) {
    (*w).Header().Set("Access-Control-Allow-Origin", "*")
}