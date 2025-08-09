package helpers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "less than a minute",
			duration: 45 * time.Second,
			expected: "00:45",
		},
		{
			name:     "one minute",
			duration: 1 * time.Minute,
			expected: "01:00",
		},
		{
			name:     "minutes and seconds",
			duration: 3*time.Minute + 7*time.Second,
			expected: "03:07",
		},
		{
			name:     "more than hour (only minutes and seconds counted)",
			duration: 1*time.Hour + 2*time.Minute + 5*time.Second,
			expected: "62:05",
		},
		{
			name:     "zero duration",
			duration: 0,
			expected: "00:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatDuration(tt.duration)
			assert.Equal(t, tt.expected, result)
		})
	}
}
