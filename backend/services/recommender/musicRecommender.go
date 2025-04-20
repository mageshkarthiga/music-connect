package recommender

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
)

const SUPABASE_URL = "https://kzxuobrnlppliqiwwgvu.supabase.co"

type Playlist struct {
	PlaylistID int `json:"playlist_id"`
	UserID     int `json:"user_id"`
}

type TrackSet map[int]bool

type PlaylistSimilarity struct {
	PlaylistID int
	Score      float64
	Tracks     TrackSet
}

type Track struct {
	TrackID       int    `json:"track_id"`
	TrackTitle    string `json:"track_title"`
	TrackURI      string `json:"track_uri"`
	TrackImageUrl string `json:"track_image_url"`
	ArtistID      int    `json:"artist_id"`
}

type RecommendationResponse struct {
	Track
	ArtistName string `json:"artist_name"`
}

func getUserPlaylistIds(userId string, apiKey string) ([]int, error) {
	url := fmt.Sprintf("%s/rest/v1/playlists?user_id=eq.%s&select=playlist_id", SUPABASE_URL, userId)
	req, err := http.NewRequest(http.MethodGet, url, nil)
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
		log.Fatalf("Failed to read response body: %v", err)
	}

	fmt.Println("Full JSON response:")
	fmt.Println(string(body))
	
	var playlistObjIds []struct {
		PlaylistID int `json:"playlist_id"`
	}
	if err := json.Unmarshal(body, &playlistObjIds); err != nil {
		return nil, err
	}

	var playlistIds []int
	for _, playlist := range playlistObjIds {
		playlistIds = append(playlistIds, playlist.PlaylistID)
	}

	return playlistIds, nil
}

func getTracksInPlaylists(playlistIDs []int, apiKey string) (TrackSet, error) {
	if len(playlistIDs) == 0 {
		return map[int]bool{}, nil
	}

	stringIds := make([]string, len(playlistIDs))
	for i, id := range playlistIDs {
		stringIds[i] = fmt.Sprintf(`"%d"`, id)
	}
	joined := strings.Join(stringIds, ",")

	url := fmt.Sprintf("%s/rest/v1/playlist_tracks?playlist_id=in.(%s)&select=track_id", SUPABASE_URL, joined)

	req, err := http.NewRequest(http.MethodGet, url, nil)
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

	var rows []struct {
		TrackID int `json:"track_id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&rows); err != nil {
		return nil, err
	}

	tracks := make(map[int]bool)
	for _, row := range rows {
		tracks[row.TrackID] = true
	}

	return tracks, nil
}

func getAllPlaylistsExceptUserPlaylists(userId string, apiKey string) ([]Playlist, error) {
	url := fmt.Sprintf("%s/rest/v1/playlists?user_id=neq.%s&select=playlist_id,user_id", SUPABASE_URL, userId)
	req, err := http.NewRequest(http.MethodGet, url, nil)
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

	var playlists []Playlist
	if err := json.NewDecoder(resp.Body).Decode(&playlists); err != nil {
		return nil, err
	}
	return playlists, nil
}

func getTrackDetails(trackID int, apiKey string) (RecommendationResponse, error) {
	url := fmt.Sprintf("%s/rest/v1/tracks?track_id=eq.%d", SUPABASE_URL, trackID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return RecommendationResponse{}, err
	}

	req.Header.Set("apikey", apiKey)
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return RecommendationResponse{}, err
	}
	defer resp.Body.Close()

	var tracks []Track
	if err := json.NewDecoder(resp.Body).Decode(&tracks); err != nil {
		return RecommendationResponse{}, err
	}

	if len(tracks) == 0 {
		return RecommendationResponse{}, fmt.Errorf("no track found with ID %d", trackID)
	}

	artistName := getArtistName(tracks[0].ArtistID, apiKey)
	if artistName == "" {
		return RecommendationResponse{}, fmt.Errorf("no artist found with ID %d", tracks[0].ArtistID)
	}

	trackDetails := RecommendationResponse{
		Track:      tracks[0],
		ArtistName: artistName,
	}
	return trackDetails, nil
}

func getArtistName(artistId int, apiKey string) string {
	url := fmt.Sprintf("%s/rest/v1/artists?artist_id=eq.%d", SUPABASE_URL, artistId)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return ""
	}

	req.Header.Set("apikey", apiKey)
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var artists []struct {
		ArtistName string `json:"artist_name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		return ""
	}

	if len(artists) == 0 {
		return ""
	}

	return artists[0].ArtistName
}

func computeJaccardScore(a, b TrackSet) float64 {
	intersection := 0
	union := make(map[int]bool)

	for track := range a {
		union[track] = true
		if b[track] {
			intersection++
		}
	}
	for track := range b {
		union[track] = true
	}

	if len(union) == 0 {
		return 0
	}
	return float64(intersection) / float64(len(union))
}

func GetTrackRecommendation(userId string) ([]RecommendationResponse, error) {
	apiKey := os.Getenv("publicApiKey")
	log.Printf("Starting track recommendation for user ID: %s", userId)

	// Get user's playlist IDs
	log.Println("Fetching user's playlist IDs...")
	userPlaylistIds, err := getUserPlaylistIds(userId, apiKey)
	if err != nil {
		log.Printf("Error fetching user's playlist IDs: %v", err)
		return nil, err
	}
	log.Printf("User's playlist IDs: %v", userPlaylistIds)

	// If user has playlists, get all the tracks in the playlists and append to a set
	var userTrackSet TrackSet
	if len(userPlaylistIds) != 0 {
		log.Println("Fetching tracks in user's playlists...")
		userTrackSet, err = getTracksInPlaylists(userPlaylistIds, apiKey)
		if err != nil {
			log.Printf("Error fetching tracks in user's playlists: %v", err)
			return nil, err
		}
		log.Printf("User's track set: %v", userTrackSet)
	} else {
		log.Println("User has no playlists.")
	}

	// Get all other playlists that exist in the database
	log.Println("Fetching all other playlists except user's playlists...")
	otherPlaylists, err := getAllPlaylistsExceptUserPlaylists(userId, apiKey)
	if err != nil {
		log.Printf("Error fetching other playlists: %v", err)
		return nil, err
	}
	log.Printf("Other playlists: %v", otherPlaylists)

	// Compute Jaccard similarity scores for other playlists
	log.Println("Computing Jaccard similarity scores...")
	var similarities []PlaylistSimilarity
	for _, playlist := range otherPlaylists {
		trackSet, err := getTracksInPlaylists([]int{playlist.PlaylistID}, apiKey)
		if err != nil {
			log.Printf("Error fetching tracks for playlist ID %d: %v", playlist.PlaylistID, err)
			continue
		}

		sim := computeJaccardScore(userTrackSet, trackSet)
		log.Printf("Computed similarity for playlist ID %d: %f", playlist.PlaylistID, sim)
		similarities = append(similarities, PlaylistSimilarity{
			PlaylistID: playlist.PlaylistID,
			Score:      sim,
			Tracks:     trackSet,
		})
	}

	// Sort the similarities by score in descending order
	log.Println("Sorting playlists by similarity scores...")
	sort.Slice(similarities, func(i, j int) bool {
		return similarities[i].Score > similarities[j].Score
	})

	// Get the top 10 tracks that the user does not have in their playlists
	log.Println("Selecting top 10 recommended tracks...")
	var recommendations []int
	for _, playlists := range similarities {
		for trackID := range playlists.Tracks {
			if len(recommendations) == 10 {
				break
			}
			if !userTrackSet[trackID] {
				recommendations = append(recommendations, trackID)
			}
		}
	}
	log.Printf("Recommended track IDs: %v", recommendations)

	// Fetch details for the recommended tracks
	log.Println("Fetching details for recommended tracks...")
	var tracks []RecommendationResponse
	for _, trackID := range recommendations {
		track, err := getTrackDetails(trackID, apiKey)
		if err != nil {
			log.Printf("Error fetching track details for ID %d: %v", trackID, err)
			continue
		}
		tracks = append(tracks, track)
	}
	log.Printf("Final recommended tracks: %v", tracks)

	log.Println("Track recommendation process completed.")
	return tracks, nil
}
