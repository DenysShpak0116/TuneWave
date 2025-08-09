package comment

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/service/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCommentHandler_CreateComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommentService := mocks.NewMockCommentService(ctrl)
	mockUserService := mocks.NewMockUserService(ctrl)
	dtoBuilder := dto.NewDTOBuilder(mockUserService, nil)

	handler := NewCommentHandler(mockCommentService, dtoBuilder)
	httpHandler := handlers.MakeHandler(handler.CreateComment)

	validSongID := uuid.New()
	validUserID := uuid.New()
	commentContent := "Great song!"

	testComment := &models.Comment{
		BaseModel: models.BaseModel{ID: uuid.New()},
		SongID:    validSongID,
		UserID:    validUserID,
		Content:   commentContent,
	}
	commentWithUser := *testComment
	commentWithUser.User = models.User{
		BaseModel: models.BaseModel{ID: validUserID},
		Username:  "testuser",
	}

	testCases := []struct {
		name           string
		expectedStatus int
		setupMocks     func()
		requestBody    map[string]any
	}{
		{
			name:           "success",
			expectedStatus: http.StatusCreated,
			setupMocks: func() {
				mockCommentService.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil)
				mockCommentService.EXPECT().
					GetByID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&commentWithUser, nil)
				mockUserService.EXPECT().
					GetUserFollowersCount(gomock.Any(), gomock.Any()).
					Return(int64(0))
			},
			requestBody: map[string]any{
				"songId":  validSongID,
				"userId":  validUserID,
				"content": commentContent,
			},
		},
		{
			name:           "invalid json",
			expectedStatus: http.StatusBadRequest,
			setupMocks:     func() {},
			requestBody: map[string]any{
				"songId": "missing_comma",
				"userId": "abc",
			},
		},
		{
			name:           "invalid song ID",
			expectedStatus: http.StatusBadRequest,
			setupMocks:     func() {},
			requestBody: map[string]any{
				"songId":  "invalid",
				"userId":  validUserID,
				"content": "test",
			},
		},
		{
			name:           "invalid user ID",
			expectedStatus: http.StatusBadRequest,
			setupMocks:     func() {},
			requestBody: map[string]any{
				"songId":  validSongID,
				"userId":  "invalid",
				"content": "test",
			},
		},
		{
			name:           "create error",
			expectedStatus: http.StatusInternalServerError,
			setupMocks: func() {
				mockCommentService.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(errors.New("db error"))
			},
			requestBody: map[string]any{
				"songId":  validSongID,
				"userId":  validUserID,
				"content": "fail",
			},
		},
		{
			name:           "get by ID error",
			expectedStatus: http.StatusInternalServerError,
			setupMocks: func() {
				mockCommentService.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil)
				mockCommentService.EXPECT().
					GetByID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("not found"))
			},
			requestBody: map[string]any{
				"songId":  validSongID,
				"userId":  validUserID,
				"content": "test",
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/comments", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedStatus == http.StatusCreated {
				var resp dto.CommentDTO
				err := json.NewDecoder(rr.Body).Decode(&resp)
				assert.NoError(t, err)
				assert.Equal(t, commentContent, resp.Content)
				assert.Equal(t, "testuser", resp.Author.Username)
			}
		})
	}
}

func TestCommentHandler_DeleteComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommentService := mocks.NewMockCommentService(ctrl)
	dtoBuilder := dto.NewDTOBuilder(nil, nil)

	handler := NewCommentHandler(mockCommentService, dtoBuilder)
	httpHandler := handlers.MakeHandler(handler.DeleteComment)

	validCommentID := uuid.New()

	testCases := []struct {
		name           string
		commentID      string
		expectedStatus int
		setupMocks     func()
	}{
		{
			name:           "success",
			commentID:      validCommentID.String(),
			expectedStatus: http.StatusNoContent,
			setupMocks: func() {
				mockCommentService.EXPECT().
					Delete(gomock.Any(), validCommentID).
					Return(nil)
			},
		},
		{
			name:           "invalid comment ID format",
			commentID:      "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
			setupMocks:     func() {},
		},
		{
			name:           "delete error",
			commentID:      validCommentID.String(),
			expectedStatus: http.StatusInternalServerError,
			setupMocks: func() {
				mockCommentService.EXPECT().
					Delete(gomock.Any(), validCommentID).
					Return(errors.New("delete failed"))
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodDelete, "/comments/"+tt.commentID, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.commentID)

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
