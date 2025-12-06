package helpers_test

import (
	"context"
	"testing"
	"time"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/golang-jwt/jwt/v5"
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

func TestParseToken(t *testing.T) {
	jwtSecret := "secret_key"
	validUserID := "123e4567-e89b-12d3-a456-426614174000"

	validToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": validUserID,
		"exp":    time.Now().Add(time.Hour).Unix(),
	})
	validTokenStr, _ := validToken.SignedString([]byte(jwtSecret))

	invalidMethodToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"userId": validUserID,
	})
	invalidMethodTokenStr, _ := invalidMethodToken.SignedString([]byte(jwtSecret))

	noUserIDToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	noUserIDTokenStr, _ := noUserIDToken.SignedString([]byte(jwtSecret))

	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": validUserID,
		"exp":    time.Now().Add(-time.Hour).Unix(),
	})
	expiredTokenStr, _ := expiredToken.SignedString([]byte(jwtSecret))

	tests := []struct {
		name        string
		tokenStr    string
		expectedID  string
		expectError bool
	}{
		{
			name:        "valid token",
			tokenStr:    validTokenStr,
			expectedID:  validUserID,
			expectError: false,
		},
		{
			name:        "invalid signing method",
			tokenStr:    invalidMethodTokenStr,
			expectedID:  "",
			expectError: true,
		},
		{
			name:        "no userId in claims",
			tokenStr:    noUserIDTokenStr,
			expectedID:  "",
			expectError: true,
		},
		{
			name:        "expired token",
			tokenStr:    expiredTokenStr,
			expectedID:  "",
			expectError: true,
		},
		{
			name:        "completely invalid token string",
			tokenStr:    "not-a-jwt",
			expectedID:  "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := helpers.ParseToken(jwtSecret, tt.tokenStr)

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
