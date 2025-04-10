package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const SPOTIFY_BASE_URL = "https://api.spotify.com/v1";

type SpotifyTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type AlbumWithImage struct {
	ID    string
	Image string
}

type Artist struct {
	ArtistID         int    `json:"artist_id"`
	ArtistName       string `json:"artist_name"`
	ArtistSpotifyID  string `json:"artist_spotify_id"`
}

type Track struct {
	TrackID 	   int    `json:"track_id"`
	TrackTitle	   string `json:"track_title"`
	ArtistID	   int    `json:"artist_id"`
	TrackURI	   string `json:"track_uri"`
	TrackSpotifyID string `json:"track_spotify_id"`
	TrackImageUrl  string `json:"track_image_url"`
}

func getSpotifyToken() (string, error) {
	tokenURL := "https://accounts.spotify.com/api/token"
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("client_id", clientID)
	form.Add("client_secret", clientSecret)

	req, err := http.NewRequest(http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
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

func getArtistsInDB() ([]string, []int, error) {
	apiKey := os.Getenv("publicApiKey")

	req, err := http.NewRequest(http.MethodGet, SUPABASE_URL+"/rest/v1/artists", nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("apikey", apiKey)
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var artists []Artist
	if err := json.Unmarshal(body, &artists); err != nil {
		return nil, nil, err
	}

	var artistSpotifyId []string
	var artistId []int
	for _, a := range artists {
		artistSpotifyId = append(artistSpotifyId, a.ArtistSpotifyID)
		artistId = append(artistId, a.ArtistID)
	}

	return artistSpotifyId, artistId, nil
}

func getAllTrackSpotifyIdsInDB() ([]string, error) {
	apiKey := os.Getenv("publicApiKey")

	req, err := http.NewRequest(http.MethodGet, SUPABASE_URL+"/rest/v1/tracks", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("apikey", apiKey)
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response []Track
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	var trackSpotifyIdList []string
	for _, t := range response {
		trackSpotifyIdList = append(trackSpotifyIdList, t.TrackSpotifyID)
	}

	return trackSpotifyIdList, nil
}

func getArtistAlbums(token string, artistSpotifyId string) ([]AlbumWithImage, error) {
	url := fmt.Sprintf(SPOTIFY_BASE_URL + "/artists/%s/albums?include_groups=album,single&limit=5", artistSpotifyId) // to update limit later
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer " + token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var albums struct {
		Items []struct {
			ID string `json:"id"`
			Images []struct {
				URL    string `json:"url"`
			} `json:"images"`
		} `json:"items"`
	}
	if err := json.Unmarshal(body, &albums); err != nil {
		return nil, err
	}

	var albumDetails []AlbumWithImage
	for _, item := range albums.Items {
		imageURL := ""
		if len(item.Images) > 0 {
			imageURL = item.Images[0].URL
		}

		albumDetails = append(albumDetails, AlbumWithImage{
			ID:    item.ID,
			Image: imageURL,
		})
	}

	return albumDetails, nil
}

func getAlbumTracks(token, albumID string, artistID int) ([]Track, error) {
	url := fmt.Sprintf(SPOTIFY_BASE_URL + "/albums/%s/tracks", albumID)

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var trackList struct {
		Items []struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			URI     string `json:"uri"`
			Artists []struct {
				ID string `json:"id"`
			} `json:"artists"`
		} `json:"items"`
	}
	if err := json.Unmarshal(body, &trackList); err != nil {
		return nil, err
	}

	var tracks []Track

	for _, t := range trackList.Items {
		tracks = append(tracks, Track{
			TrackTitle:     t.Name,
			TrackSpotifyID: t.ID,
			TrackURI:       t.URI,
			ArtistID:       artistID,
		})
	}

	return tracks, nil
}

func insertTrackIntoSupabase(t Track, image string) (int, error) {
	apiKey := os.Getenv("publicApiKey")
	url := SUPABASE_URL + "/rest/v1/tracks"

	data := map[string]interface{}{
		"track_spotify_id": t.TrackSpotifyID,
		"track_title":      t.TrackTitle,
		"artist_id":        t.ArtistID,
		"track_uri":        t.TrackURI,
		"track_image_url":  image,
	}

	payload, _ := json.Marshal(data)
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", apiKey)
	req.Header.Set("Authorization", "Bearer " + apiKey)
	req.Header.Set("Prefer", "return=representation")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	
	var insertedTrack []Track
	if err := json.Unmarshal(body, &insertedTrack); err != nil {
		return -1, fmt.Errorf("unmarshal error: %w", err)
	}
	return insertedTrack[0].TrackID, nil
}

func insertTrackArtistIntoSupabase(trackID int, artistID int) error {
	apiKey := os.Getenv("publicApiKey")
	url := SUPABASE_URL + "/rest/v1/track_artists"

	data := map[string]interface{}{
		"track_id":  trackID,
		"artist_id": artistID,
	}

	payload, _ := json.Marshal(data)
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", apiKey)
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func callSpotifyAPI() {
	artistSpotifyIds, artistIds, err := getArtistsInDB()
	if err != nil {
		log.Fatal("Error getting artists:", err)
	}

	trackSpotifyIds, err := getAllTrackSpotifyIdsInDB()
	if err != nil {
		log.Fatal("Error getting tracks:", err)
	}
	trackSpotifyIDSet := toSet(trackSpotifyIds)

	token, _ := getSpotifyToken()

	for index, artistSpotifyId := range artistSpotifyIds {
		albumsDetails, err := getArtistAlbums(token, artistSpotifyId)
		if err != nil {
			log.Fatalln("Failed to fetch albums:", err)
		}

		for _, album := range albumsDetails {
			tracks, err := getAlbumTracks(token, album.ID, artistIds[index])
			if err != nil {
				log.Printf("Skipping album %s due to error: %v\n", album, err)
				continue
			}

			for _, track := range tracks {
				if trackSpotifyIDSet[track.TrackSpotifyID] {
					continue
				}				

				trackId, err := insertTrackIntoSupabase(track, album.Image)
				if err != nil {
					log.Fatalln("Insert track error: ", err)
					continue
				}

				if err := insertTrackArtistIntoSupabase(trackId, artistIds[index]); err != nil {
					log.Fatalln("Insert track_artist error: ", err)
				}
			}
			
		}
	}
	log.Println("Spotify API call completed successfully.")
}
