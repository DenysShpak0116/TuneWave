package song

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/service"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/service/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSongHandler_SetReaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSongService := mocks.NewMockSongService(ctrl)
	mockCollectionSongService := mocks.NewMockCollectionSongService(ctrl)
	mockUserReactionService := mocks.NewMockUserReactionService(ctrl)
	mockCommentService := mocks.NewMockCommentService(ctrl)
	mockUserService := mocks.NewMockUserService(ctrl)
	mockEventService := &service.EventService{}

	dtoBuilder := dto.NewDTOBuilder(mockUserService, mockUserReactionService)
	var logger *slog.Logger

	handler := NewSongHandler(
		mockSongService,
		mockCollectionSongService,
		mockUserReactionService,
		mockCommentService,
		mockEventService,
		dtoBuilder,
		logger,
	)

	httpHandler := handlers.MakeHandler(handler.SetReaction)

	validSongID := uuid.New()
	validUserID := uuid.New()

	testCases := []struct {
		name         string
		songID       string
		body         any
		expectedCode int
		setupMocks   func()
	}{
		{
			name:         "success",
			songID:       validSongID.String(),
			expectedCode: http.StatusOK,
			body: map[string]any{
				"userID":       validUserID.String(),
				"reactionType": "like",
			},
			setupMocks: func() {
				mockSongService.EXPECT().
					SetReaction(gomock.Any(), validSongID, validUserID, "like").
					Return(10, 2, "", nil)
			},
		},
		{
			name:         "invalid song ID",
			songID:       "invalid-id",
			body:         map[string]any{},
			expectedCode: http.StatusBadRequest,
			setupMocks:   func() {},
		},
		{
			name:         "invalid request body",
			songID:       validSongID.String(),
			body:         "invalid-json",
			expectedCode: http.StatusBadRequest,
			setupMocks:   func() {},
		},
		{
			name:   "invalid user ID format",
			songID: validSongID.String(),
			body: map[string]any{
				"userID":       "invalid-user-id",
				"reactionType": "like",
			},
			expectedCode: http.StatusBadRequest,
			setupMocks:   func() {},
		},
		{
			name:   "set reaction error",
			songID: validSongID.String(),
			body: map[string]any{
				"userID":       validUserID.String(),
				"reactionType": "dislike",
			},
			expectedCode: http.StatusInternalServerError,
			setupMocks: func() {
				mockSongService.EXPECT().
					SetReaction(gomock.Any(), validSongID, validUserID, "dislike").
					Return(0, 0, "", errors.New("some error"))
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			var body io.Reader
			switch v := tt.body.(type) {
			case string:
				body = strings.NewReader(v)
			default:
				b, _ := json.Marshal(v)
				body = bytes.NewReader(b)
			}

			req := httptest.NewRequest(http.MethodPost, "/songs/"+tt.songID+"/reaction", body)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.songID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
		})
	}
}

func TestSongHandler_CheckReaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSongService := mocks.NewMockSongService(ctrl)
	mockCollectionSongService := mocks.NewMockCollectionSongService(ctrl)
	mockUserReactionService := mocks.NewMockUserReactionService(ctrl)
	mockCommentService := mocks.NewMockCommentService(ctrl)
	mockUserService := mocks.NewMockUserService(ctrl)
	mockEventService := &service.EventService{}

	dtoBuilder := dto.NewDTOBuilder(mockUserService, mockUserReactionService)
	var logger *slog.Logger

	handler := NewSongHandler(
		mockSongService,
		mockCollectionSongService,
		mockUserReactionService,
		mockCommentService,
		mockEventService,
		dtoBuilder,
		logger,
	)
	httpHandler := handlers.MakeHandler(handler.CheckReaction)

	validSongID := uuid.New()
	validUserID := uuid.New()

	testCases := []struct {
		name         string
		songID       string
		userID       string
		expectedCode int
		expectedBody string
		setupMocks   func()
	}{
		{
			name:         "user id is 'undefined'",
			songID:       validSongID.String(),
			userID:       "undefined",
			expectedCode: http.StatusOK,
			expectedBody: `{"type":"none"}`,
			setupMocks:   func() {},
		},
		{
			name:         "invalid user id",
			songID:       validSongID.String(),
			userID:       "invalid",
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"invalid user id"}`,
			setupMocks:   func() {},
		},
		{
			name:         "invalid song id",
			songID:       "bad-song-id",
			userID:       validUserID.String(),
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"invalid song id"}`,
			setupMocks:   func() {},
		},
		{
			name:         "internal error from service",
			songID:       validSongID.String(),
			userID:       validUserID.String(),
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"error":"failed to check reaction"}`,
			setupMocks: func() {
				mockSongService.EXPECT().
					IsReactedByUser(gomock.Any(), validSongID, validUserID).
					Return("", errors.New("some error"))
			},
		},
		{
			name:         "success - liked",
			songID:       validSongID.String(),
			userID:       validUserID.String(),
			expectedCode: http.StatusOK,
			expectedBody: `{"type":"like"}`,
			setupMocks: func() {
				mockSongService.EXPECT().
					IsReactedByUser(gomock.Any(), validSongID, validUserID).
					Return("like", nil)
			},
		},
		{
			name:         "success - no reaction",
			songID:       validSongID.String(),
			userID:       validUserID.String(),
			expectedCode: http.StatusOK,
			expectedBody: `{"type":"none"}`,
			setupMocks: func() {
				mockSongService.EXPECT().
					IsReactedByUser(gomock.Any(), validSongID, validUserID).
					Return("none", nil)
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodGet, "/songs/"+tt.songID+"/reaction/"+tt.userID, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.songID)
			rctx.URLParams.Add("userId", tt.userID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func TestSongHandler_ListenSong(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSongService := mocks.NewMockSongService(ctrl)
	mockCollectionSongService := mocks.NewMockCollectionSongService(ctrl)
	mockUserReactionService := mocks.NewMockUserReactionService(ctrl)
	mockCommentService := mocks.NewMockCommentService(ctrl)
	mockUserService := mocks.NewMockUserService(ctrl)
	mockEventService := &service.EventService{}

	dtoBuilder := dto.NewDTOBuilder(mockUserService, mockUserReactionService)
	var logger *slog.Logger

	handler := NewSongHandler(
		mockSongService,
		mockCollectionSongService,
		mockUserReactionService,
		mockCommentService,
		mockEventService,
		dtoBuilder,
		logger,
	)
	httpHandler := handlers.MakeHandler(handler.ListenSong)

	validSongID := uuid.New()
	validUserID := uuid.New()
	existingSong := &models.Song{
		BaseModel:  models.BaseModel{ID: validSongID},
		Listenings: 5,
	}

	testCases := []struct {
		name         string
		songID       string
		userID       string
		expectedCode int
		setupMocks   func()
	}{
		{
			name:         "userID is undefined",
			songID:       validSongID.String(),
			userID:       "undefined",
			expectedCode: http.StatusNoContent,
			setupMocks:   func() {},
		},
		{
			name:         "invalid songID",
			songID:       "not-a-uuid",
			userID:       validUserID.String(),
			expectedCode: http.StatusBadRequest,
			setupMocks:   func() {},
		},
		{
			name:         "song does not exist",
			songID:       validSongID.String(),
			userID:       validUserID.String(),
			expectedCode: http.StatusBadRequest,
			setupMocks: func() {
				mockSongService.EXPECT().
					GetByID(gomock.Any(), validSongID).
					Return(nil, errors.New("not found"))
			},
		},
		{
			name:         "update error",
			songID:       validSongID.String(),
			userID:       validUserID.String(),
			expectedCode: http.StatusInternalServerError,
			setupMocks: func() {
				mockSongService.EXPECT().
					GetByID(gomock.Any(), validSongID).
					Return(existingSong, nil)
				mockSongService.EXPECT().
					Update(gomock.Any(), &models.Song{
						BaseModel:  models.BaseModel{ID: existingSong.ID},
						Listenings: existingSong.Listenings + 1,
					}).
					Return(errors.New("db error"))
			},
		},
		{
			name:         "success",
			songID:       validSongID.String(),
			userID:       validUserID.String(),
			expectedCode: http.StatusNoContent,
			setupMocks: func() {
				mockSongService.EXPECT().
					GetByID(gomock.Any(), validSongID).
					Return(existingSong, nil)
				mockSongService.EXPECT().
					Update(gomock.Any(), &models.Song{
						BaseModel:  models.BaseModel{ID: existingSong.ID},
						Listenings: existingSong.Listenings + 1,
					}).
					Return(nil)
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodPost, "/songs/"+tt.songID+"/listen/"+tt.userID, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.songID)
			rctx.URLParams.Add("userId", tt.userID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
		})
	}
}
