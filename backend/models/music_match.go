// models/cosine_similarity.go
package models

import (
	"math"
	"strconv"
)


func CosineSimilarity(userProfile1, userProfile2 UserProfile) (float64, error) {
	// Helper function to aggregate track data with custom weights
	aggregateTracks := func(user UserProfile) map[string]float64 {
		trackVector := make(map[string]float64)

		// Weight factors for different track types
		const likedWeight = 3.0
		const playedWeight = 1.0
		const playlistWeight = 1.5
		
		for id, count := range user.LikedTracks {
			trackVector[strconv.Itoa(int(id))] += float64(count) * likedWeight
		}
		for id, count := range user.PlayedTracks {
			trackVector[strconv.Itoa(int(id))] += float64(count) * playedWeight
		}
		for id, count := range user.PlaylistTracks {
			trackVector[strconv.Itoa(int(id))] += float64(count) * playlistWeight
		}

		return trackVector
	}

	vec1 := aggregateTracks(userProfile1)
	vec2 := aggregateTracks(userProfile2)

	// Calculate dot product and magnitudes
	var dotProduct, magnitude1, magnitude2 float64

	// Dot product
	for trackID, val1 := range vec1 {
		if val2, exists := vec2[trackID]; exists {
			dotProduct += val1 * val2
		}
	}

	// Magnitude for user 1
	for _, val := range vec1 {
		magnitude1 += val * val
	}

	// Magnitude for user 2
	for _, val := range vec2 {
		magnitude2 += val * val
	}

	// Avoid division by zero
	if magnitude1 == 0 || magnitude2 == 0 {
		return 0, nil
	}

	// Cosine similarity
	return dotProduct / (math.Sqrt(magnitude1) * math.Sqrt(magnitude2)), nil
}
