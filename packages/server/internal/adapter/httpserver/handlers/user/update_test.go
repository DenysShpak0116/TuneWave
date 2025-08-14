package user

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/service/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUserHandler_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	mockUserFollowerService := mocks.NewMockUserFollowerService(ctrl)
	mockUserReactionService := mocks.NewMockUserReactionService(ctrl)
	mockMessageService := mocks.NewMockMessageService(ctrl)

	dtoBuilder := dto.NewDTOBuilder(mockUserService, mockUserReactionService)

	handler := NewUserHandler(
		mockUserService,
		mockUserFollowerService,
		mockUserReactionService,
		mockMessageService,
		dtoBuilder,
	)

	userID := uuid.New()
	httpHandler := handlers.MakeHandler(handler.Update)

	tests := []struct {
		name           string
		idParam        string
		body           string
		setupMocks     func()
		expectedStatus int
	}{
		{
			name:    "success",
			idParam: userID.String(),
			body:    `{"username":"newuser", "profileInfo":"updated bio"}`,
			setupMocks: func() {
				mockUserService.EXPECT().Update(gomock.Any(), &models.User{
					BaseModel:   models.BaseModel{ID: userID},
					Username:    "newuser",
					ProfileInfo: "updated bio",
				}).Return(nil)

				mockUserService.EXPECT().GetUserFollowersCount(gomock.Any(), userID).Return(int64(0))
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid UUID",
			idParam:        "invalid-uuid",
			body:           `{"username":"user"}`,
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid JSON payload",
			idParam:        userID.String(),
			body:           `{"username": "missing-quote}`,
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "update fails",
			idParam: userID.String(),
			body:    `{"username":"failuser", "profileInfo":"test"}`,
			setupMocks: func() {
				mockUserService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("update error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodPut, "/users/"+tt.idParam, bytes.NewBufferString(tt.body))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
				URLParams: chi.RouteParams{
					Keys:   []string{"id"},
					Values: []string{tt.idParam},
				},
			}))
			rr := httptest.NewRecorder()

			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			var actualBody map[string]interface{}
			err := json.NewDecoder(rr.Body).Decode(&actualBody)
			require.NoError(t, err)
		})
	}
}

func TestUserHandler_UpdateAvatar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	mockUserFollowerService := mocks.NewMockUserFollowerService(ctrl)
	mockUserReactionService := mocks.NewMockUserReactionService(ctrl)
	mockMessageService := mocks.NewMockMessageService(ctrl)

	dtoBuilder := dto.NewDTOBuilder(mockUserService, mockUserReactionService)

	handler := NewUserHandler(
		mockUserService,
		mockUserFollowerService,
		mockUserReactionService,
		mockMessageService,
		dtoBuilder,
	)

	userID := uuid.New()
	httpHandler := handlers.MakeHandler(handler.UpdateAvatar)

	tests := []struct {
		name           string
		setupRequest   func() *http.Request
		setupMocks     func()
		expectedStatus int
	}{
		{
			name: "success",
			setupRequest: func() *http.Request {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)

				part, _ := writer.CreateFormFile("file", "avatar.png")
				part.Write([]byte("image data"))

				writer.Close()

				req := httptest.NewRequest(http.MethodPost, "/users/avatar", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())

				ctx := context.WithValue(req.Context(), "userID", userID.String())
				return req.WithContext(ctx)
			},
			setupMocks: func() {
				mockUserService.EXPECT().UpdateUserPfp(gomock.Any(), gomock.AssignableToTypeOf(services.UpdatePfpParams{})).
					Return(nil).AnyTimes()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid form",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/users/avatar", bytes.NewBufferString("not a form"))
				req.Header.Set("Content-Type", "text/plain")
				ctx := context.WithValue(req.Context(), "userID", userID.String())
				return req.WithContext(ctx)
			},
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid form file",
			setupRequest: func() *http.Request {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)

				writer.Close()

				req := httptest.NewRequest(http.MethodPost, "/users/avatar", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				ctx := context.WithValue(req.Context(), "userID", userID.String())
				return req.WithContext(ctx)
			},
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid user ID",
			setupRequest: func() *http.Request {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("file", "avatar.png")
				part.Write([]byte("image"))
				writer.Close()

				req := httptest.NewRequest(http.MethodPost, "/users/avatar", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())

				return req
			},
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := tt.setupRequest()
			rr := httptest.NewRecorder()

			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			var actual map[string]string
			err := json.NewDecoder(rr.Body).Decode(&actual)
			require.NoError(t, err)
		})
	}
}
