package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractS3Key(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected string
	}{
		{
			name:     "valid S3 URL",
			inputURL: "https://tunewavebucket.s3.eu-west-3.amazonaws.com/uploads/audio/song.mp3",
			expected: "uploads/audio/song.mp3",
		},
		{
			name:     "URL without base prefix",
			inputURL: "https://otherbucket.s3.eu-west-3.amazonaws.com/uploads/audio/song.mp3",
			expected: "https://otherbucket.s3.eu-west-3.amazonaws.com/uploads/audio/song.mp3",
		},
		{
			name:     "empty string",
			inputURL: "",
			expected: "",
		},
		{
			name:     "only base URL",
			inputURL: "https://tunewavebucket.s3.eu-west-3.amazonaws.com/",
			expected: "",
		},
		{
			name:     "partial prefix match",
			inputURL: "https://tunewavebucket.s3.eu-west-3.amazonaws.comed/file.txt",
			expected: "https://tunewavebucket.s3.eu-west-3.amazonaws.comed/file.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractS3Key(tt.inputURL)
			assert.Equal(t, tt.expected, result)
		})
	}
}
