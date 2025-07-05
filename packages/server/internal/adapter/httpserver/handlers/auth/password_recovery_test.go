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
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/service/mocks"
	"github.com/go-chi/render"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestForgotPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthService(ctrl)
	mockUserService := mocks.NewMockUserService(ctrl)
	dtoBuilder := dto.NewDTOBuilder(mockUserService, nil)

	cfg := &config.Config{
		JwtSecret: "test-secret",
		Google:    config.GoogleConfig{},
	}

	handler := NewAuthHandler(mockAuthService, mockUserService, dtoBuilder, cfg)

	tests := []struct {
		name           string
		body           map[string]any
		mockSetup      func()
		expectedStatus int
	}{
		{
			name: "success",
			body: map[string]any{
				"email": "user@example.com",
			},
			mockSetup: func() {
				mockAuthService.EXPECT().
					HandleForgotPassword("user@example.com").
					Return("reset-token", nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid JSON body",
			body:           nil,
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service error",
			body: map[string]any{
				"email": "fail@example.com",
			},
			mockSetup: func() {
				mockAuthService.EXPECT().
					HandleForgotPassword("fail@example.com").
					Return("", errors.New("email error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			var req *http.Request
			if tt.body != nil {
				jsonBody, _ := json.Marshal(tt.body)
				req = httptest.NewRequest(http.MethodPost, "/auth/forgot-password", bytes.NewReader(jsonBody))
			} else {
				req = httptest.NewRequest(http.MethodPost, "/auth/forgot-password", bytes.NewReader([]byte("bad-json")))
			}

			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			render.Status(req, 0)

			httpHandler := handlers.MakeHandler(handler.ForgotPassword)
			httpHandler.ServeHTTP(rr, req)

			res := rr.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)
		})
	}
}

func TestResetPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthService(ctrl)
	mockUserService := mocks.NewMockUserService(ctrl)
	dtoBuilder := dto.NewDTOBuilder(mockUserService, nil)

	cfg := &config.Config{
		JwtSecret: "test-secret",
		Google:    config.GoogleConfig{},
	}

	handler := NewAuthHandler(mockAuthService, mockUserService, dtoBuilder, cfg)

	tests := []struct {
		name           string
		body           map[string]any
		mockSetup      func()
		expectedStatus int
	}{
		{
			name: "success",
			body: map[string]any{
				"token":       "valid-token",
				"newPassword": "newpass123",
			},
			mockSetup: func() {
				mockAuthService.EXPECT().
					HandleResetPassword("valid-token", "newpass123").
					Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid JSON",
			body:           nil,
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service failure",
			body: map[string]any{
				"token":       "bad-token",
				"newPassword": "newpass123",
			},
			mockSetup: func() {
				mockAuthService.EXPECT().
					HandleResetPassword("bad-token", "newpass123").
					Return(errors.New("invalid token"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			var req *http.Request
			if tt.body != nil {
				jsonBody, _ := json.Marshal(tt.body)
				req = httptest.NewRequest(http.MethodPost, "/auth/reset-password", bytes.NewReader(jsonBody))
			} else {
				req = httptest.NewRequest(http.MethodPost, "/auth/reset-password", bytes.NewReader([]byte("bad-json")))
			}

			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			render.Status(req, 0)

			httpHandler := handlers.MakeHandler(handler.ResetPassword)
			httpHandler.ServeHTTP(rr, req)

			res := rr.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)
		})
	}
}
