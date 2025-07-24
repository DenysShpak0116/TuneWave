package helpers_test

import (
	"context"
	"testing"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetUserID(t *testing.T) {
	validUUID := uuid.New()

	tests := []struct {
		name        string
		ctx         context.Context
		expectedID  uuid.UUID
		expectError bool
	}{
		{
			name:        "success",
			ctx:         context.WithValue(context.Background(), "userID", validUUID.String()),
			expectedID:  validUUID,
			expectError: false,
		},
		{
			name:        "userID not found in context",
			ctx:         context.Background(),
			expectedID:  uuid.UUID{},
			expectError: true,
		},
		{
			name:        "invalid UUID format",
			ctx:         context.WithValue(context.Background(), "userID", "not-a-uuid"),
			expectedID:  uuid.UUID{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := helpers.GetUserID(tt.ctx)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedID, userID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, userID)
			}
		})
	}
}
