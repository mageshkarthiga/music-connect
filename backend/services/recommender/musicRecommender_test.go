package recommender

import (
	"testing"
)

func TestComputeJaccardScore(t *testing.T) {
	tests := []struct {
		name     string
		a        TrackSet
		b        TrackSet
		expected float64
	}{
		{
			name:     "No overlap",
			a:        TrackSet{1: true, 2: true},
			b:        TrackSet{3: true, 4: true},
			expected: 0.0,
		},
		{
			name:     "Some overlap",
			a:        TrackSet{1: true, 2: true, 3: true},
			b:        TrackSet{2: true, 3: true, 4: true},
			expected: 2.0 / 4.0, // intersection: 2 (2,3); union: 4 (1,2,3,4)
		},
		{
			name:     "Complete overlap",
			a:        TrackSet{1: true, 2: true},
			b:        TrackSet{1: true, 2: true},
			expected: 1.0,
		},
		{
			name:     "Empty sets",
			a:        TrackSet{},
			b:        TrackSet{},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := computeJaccardScore(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

