package services

import (
    "net/http"
)

func CallSpotifyAPI() (*http.Response, error) {
    resp, err := http.Get("https://api.spotify.com/v1/albums")
    if err != nil {
        return nil, err
    }
    return resp, nil
}
