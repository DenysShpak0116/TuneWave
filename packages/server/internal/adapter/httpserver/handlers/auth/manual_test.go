package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/config"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/service"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/service/mocks"
	"github.com/go-chi/render"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthService(ctrl)
	mockUserService := mocks.NewMockUserService(ctrl)
	dtoBuilder := dto.NewDTOBuilder(mockUserService, nil)

	cfg := &config.Config{
		JwtSecret: "test-secret",
		Google: config.GoogleConfig{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
		},
	}

	handler := NewAuthHandler(mockAuthService, mockUserService, dtoBuilder, cfg)

	tests := []struct {
		name           string
		body           map[string]any
		mockSetup      func()
		expectedStatus int
	}{
		{
			name: "successful registration",
			body: map[string]any{
				"email":    "example@gmail.com",
				"password": "password123",
				"username": "testuser",
			},
			mockSetup: func() {
				mockUserService.EXPECT().
					First(gomock.Any(), &models.User{Email: "example@gmail.com"}).
					Return(nil, service.ErrNotFound)

				mockUserService.EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(&models.User{})).
					Return(nil)

				mockUserService.EXPECT().
					GetUserFollowersCount(gomock.Any(), gomock.Any()).
					Return(0)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "user already exists",
			body: map[string]any{
				"email":    "taken@example.com",
				"password": "password123",
				"username": "testuser",
			},
			mockSetup: func() {
				mockUserService.EXPECT().
					First(gomock.Any(), &models.User{Email: "taken@example.com"}).
					Return(&models.User{Email: "taken@example.com"}, nil)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid JSON body",
			body:           nil,
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "error during create",
			body: map[string]any{
				"email":    "newuser@example.com",
				"password": "password123",
				"username": "newuser",
			},
			mockSetup: func() {
				mockUserService.EXPECT().
					First(gomock.Any(), &models.User{Email: "newuser@example.com"}).
					Return(nil, service.ErrNotFound)

				mockUserService.EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(&models.User{})).
					Return(errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			var reqBody []byte
			var req *http.Request

			if tt.body != nil {
				reqBody, _ = json.Marshal(tt.body)
				req = httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBody))
			} else {
				req = httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader([]byte("invalid-json")))
			}

			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			render.Status(req, 0)

			httpHandler := handlers.MakeHandler(handler.Register)
			httpHandler.ServeHTTP(rr, req)

			res := rr.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)
		})
	}
}
