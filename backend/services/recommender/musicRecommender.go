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
	TrackID    int `json:"track_id"`
	TrackTitle string `json:"track_title"`
	TrackURI   string `json:"track_uri"`
	TrackImageUrl string `json:"track_image_url"`
	ArtistID   int `json:"artist_id"`
}

func getUserPlaylistIds(userId string, apiKey string) ([]int, error) {
	url := fmt.Sprintf("%s/rest/v1/playlists?user_id=eq.%s&select=playlist_id", SUPABASE_URL, userId)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("apikey", apiKey)
	req.Header.Set("Authorization", "Bearer "+ apiKey)
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
	req.Header.Set("Authorization", "Bearer "+ apiKey)
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

func getTrackDetails(trackID int, apiKey string) (Track, error) {
	url := fmt.Sprintf("%s/rest/v1/tracks?track_id=eq.%d", SUPABASE_URL, trackID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Track{}, err
	}

	req.Header.Set("apikey", apiKey)
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Track{}, err
	}
	defer resp.Body.Close()

	var tracks []Track
	if err := json.NewDecoder(resp.Body).Decode(&tracks); err != nil {
		return Track{}, err
	}

	if len(tracks) == 0 {
		return Track{}, fmt.Errorf("no track found with ID %d", trackID)
	}

	fmt.Printf("Track Details: %+v\n", tracks[0]) // HERE
	return tracks[0], nil
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


func GetTrackRecommendation(userId string) ([]Track, error) {
	apiKey := os.Getenv("publicApiKey")

	userPlaylistIds, err := getUserPlaylistIds(userId, apiKey)
	if err != nil {
		return nil, err
	}
	fmt.Println("CHECKPOINT 1: User Playlist IDs:", userPlaylistIds) // HERE

	userTrackSet, err := getTracksInPlaylists(userPlaylistIds, apiKey)
	if err != nil {
		return nil, err
	}
	fmt.Println("CHECKPOINT 2: User Track Set:", userTrackSet) // HERE

	otherPlaylists, err := getAllPlaylistsExceptUserPlaylists(userId, apiKey)
	if err != nil {
		return nil, err
	}
	fmt.Printf("CHECKPOINT 3: Other Playlists: %+v\n", otherPlaylists) // HERE

	var similarities []PlaylistSimilarity
	for _, playlist := range otherPlaylists {
		trackSet, err := getTracksInPlaylists([]int{playlist.PlaylistID}, apiKey)
		if err != nil {
			continue
		}

		sim := computeJaccardScore(userTrackSet, trackSet)
		similarities = append(similarities, PlaylistSimilarity{
			PlaylistID: playlist.PlaylistID,
			Score:      sim,
			Tracks:     trackSet,
		})
	}

	fmt.Printf("CHECKPOINT 4: Similarities: %+v\n", similarities) // HERE
	sort.Slice(similarities, func(i, j int) bool {
		return similarities[i].Score > similarities[j].Score
	})

	if len(similarities) == 0 {
		return []Track{}, nil
	}

	top := similarities[0]
	var recommendations []int
	for trackID := range top.Tracks {
		if !userTrackSet[trackID] {
			recommendations = append(recommendations, trackID)
		}
		if len(recommendations) == 10 {
			break
		}
	}
	fmt.Println("CHECKPOINT 5: Recommendations:", recommendations)

	var tracks []Track
	for _, trackID := range recommendations {
		track, err := getTrackDetails(trackID, apiKey)
		if err != nil {
			log.Fatalf("Error fetching track details for ID %d: %v\n", trackID, err)
		}
		tracks = append(tracks, track)
	}

	return tracks, nil
}
