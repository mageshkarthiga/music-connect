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
	fmt.Println("CHECKPOINT MUSIC RECOMMENDER 3: ", playlists)
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

	fmt.Println("CHECKPOINT MUSIC RECOMMENDER 4: ", trackDetails)
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

	fmt.Println("CHECKPOINT MUSIC RECOMMENDER 5: ", artists[0].ArtistName)
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

	// get user's playlists ids
	userPlaylistIds, err := getUserPlaylistIds(userId, apiKey)
	if err != nil {
		return nil, err
	}

	// if user has playlists, get all the tracks in the playlist and append to a set
	var userTrackSet TrackSet
	if len(userPlaylistIds) != 0 {
		userTrackSet, err = getTracksInPlaylists(userPlaylistIds, apiKey)
		if err != nil {
			return nil, err
		}
	}

	// get all other playlists that exist in db
	otherPlaylists, err := getAllPlaylistsExceptUserPlaylists(userId, apiKey)
	if err != nil {
		return nil, err
	}

	// get all the tracks in the other playlists and compute jaccard score of playlist similarity to all the songs the user has
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

	// sort the similarities by score in descending order (highest score first)
	// then get the top 10 tracks which the user does not have in their playlist
	sort.Slice(similarities, func(i, j int) bool {
		return similarities[i].Score > similarities[j].Score
	})

	var recommendations []int
	for _, playlists := range similarities {
		fmt.Println("CHECKPOINT MUSIC RECOMMENDER playlists: ", playlists)
		for trackID := range playlists.Tracks {
			fmt.Println("CHECKPOINT MUSIC RECOMMENDER trackID: ", trackID)
			if len(recommendations) == 10 {
				break
			}
			if !userTrackSet[trackID] {
				recommendations = append(recommendations, trackID)
			}
		}
	}

	fmt.Println("CHECKPOINT MUSIC RECOMMENDER recommendations: ", recommendations)
	var tracks []RecommendationResponse
	for _, trackID := range recommendations {
		track, err := getTrackDetails(trackID, apiKey)
		if err != nil {
			log.Fatalf("Error fetching track details for ID %d: %v\n", trackID, err)
		}
		tracks = append(tracks, track)
	}

	fmt.Println("CHECKPOINT MUSIC RECOMMENDER: ", tracks)
	return tracks, nil
}
