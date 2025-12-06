package song

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSongHandler_AddToCollection(t *testing.T) {
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

	httpHandler := handlers.MakeHandler(handler.AddToCollection)

	validSongID := uuid.New()
	validCollectionID := uuid.New()

	requestBody := func(collectionID string) io.Reader {
		body := map[string]string{"collectionId": collectionID}
		b, _ := json.Marshal(body)
		return bytes.NewReader(b)
	}

	testCases := []struct {
		name             string
		songID           string
		requestBody      io.Reader
		setupMocks       func()
		expectedCode     int
		expectedResponse string
	}{
		{
			name:        "success",
			songID:      validSongID.String(),
			requestBody: requestBody(validCollectionID.String()),
			setupMocks: func() {
				mockCollectionSongService.
					EXPECT().
					First(gomock.Any(), &models.CollectionSong{
						SongID:       validSongID,
						CollectionID: validCollectionID,
					}).
					Return(nil, service.ErrNotFound)

				mockSongService.
					EXPECT().
					AddToCollection(gomock.Any(), validSongID, validCollectionID).
					Return(nil)
			},
			expectedCode:     http.StatusOK,
			expectedResponse: `{"message":"Song added to collection"}`,
		},
		{
			name:        "already in collection",
			songID:      validSongID.String(),
			requestBody: requestBody(validCollectionID.String()),
			setupMocks: func() {
				mockCollectionSongService.
					EXPECT().
					First(gomock.Any(), &models.CollectionSong{
						SongID:       validSongID,
						CollectionID: validCollectionID,
					}).
					Return(&models.CollectionSong{}, nil)
			},
			expectedCode:     http.StatusConflict,
			expectedResponse: `{"error":"this song is already in collection"}`,
		},
		{
			name:        "first() internal error",
			songID:      validSongID.String(),
			requestBody: requestBody(validCollectionID.String()),
			setupMocks: func() {
				mockCollectionSongService.
					EXPECT().
					First(gomock.Any(), &models.CollectionSong{
						SongID:       validSongID,
						CollectionID: validCollectionID,
					}).
					Return(nil, errors.New("db error"))
			},
			expectedCode:     http.StatusInternalServerError,
			expectedResponse: `{"error":"could not add song to collection"}`,
		},
		{
			name:        "addToCollection error",
			songID:      validSongID.String(),
			requestBody: requestBody(validCollectionID.String()),
			setupMocks: func() {
				mockCollectionSongService.
					EXPECT().
					First(gomock.Any(), &models.CollectionSong{
						SongID:       validSongID,
						CollectionID: validCollectionID,
					}).
					Return(nil, service.ErrNotFound)

				mockSongService.
					EXPECT().
					AddToCollection(gomock.Any(), validSongID, validCollectionID).
					Return(errors.New("add error"))
			},
			expectedCode:     http.StatusInternalServerError,
			expectedResponse: `{"error":"failed to add song to collection"}`,
		},
		{
			name:             "invalid collection ID",
			songID:           validSongID.String(),
			requestBody:      requestBody("invalid-uuid"),
			setupMocks:       func() {},
			expectedCode:     http.StatusBadRequest,
			expectedResponse: `{"error":"invalid collection ID"}`,
		},
		{
			name:             "invalid song ID",
			songID:           "invalid-uuid",
			requestBody:      requestBody(validCollectionID.String()),
			setupMocks:       func() {},
			expectedCode:     http.StatusBadRequest,
			expectedResponse: `{"error":"invalid song ID"}`,
		},
		{
			name:             "invalid request body",
			songID:           validSongID.String(),
			requestBody:      strings.NewReader("bad-json"),
			setupMocks:       func() {},
			expectedCode:     http.StatusBadRequest,
			expectedResponse: `{"error":"invalid request body"}`,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			url := fmt.Sprintf("/songs/%s/collection", tt.songID)
			req := httptest.NewRequest(http.MethodPost, url, tt.requestBody)
			req = muxSetURLParam(req, "id", tt.songID)

			rr := httptest.NewRecorder()
			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)

			if tt.expectedResponse != "" {
				assert.JSONEq(t, tt.expectedResponse, rr.Body.String())
			}
		})
	}
}

func TestSongHandler_RemoveFromCollection(t *testing.T) {
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

	httpHandler := handlers.MakeHandler(handler.RemoveFromCollection)

	validSongID := uuid.New()
	validCollectionID := uuid.New()
	collectionSongID := uuid.New()

	requestBody := func(collectionID string) io.Reader {
		body := map[string]string{"collectionId": collectionID}
		b, _ := json.Marshal(body)
		return bytes.NewReader(b)
	}

	testCases := []struct {
		name             string
		songID           string
		requestBody      io.Reader
		setupMocks       func()
		expectedCode     int
		expectedResponse string
	}{
		{
			name:        "success",
			songID:      validSongID.String(),
			requestBody: requestBody(validCollectionID.String()),
			setupMocks: func() {
				mockCollectionSongService.
					EXPECT().
					First(gomock.Any(), &models.CollectionSong{
						CollectionID: validCollectionID,
						SongID:       validSongID,
					}).
					Return(&models.CollectionSong{
						BaseModel: models.BaseModel{
							ID: collectionSongID,
						},
						CollectionID: validCollectionID,
						SongID:       validSongID,
					}, nil)

				mockCollectionSongService.
					EXPECT().
					Delete(gomock.Any(), collectionSongID).
					Return(nil)
			},
			expectedCode:     http.StatusOK,
			expectedResponse: `{"message":"Song removed from collection"}`,
		},
		{
			name:             "invalid request body",
			songID:           validSongID.String(),
			requestBody:      strings.NewReader("invalid-json"),
			setupMocks:       func() {},
			expectedCode:     http.StatusBadRequest,
			expectedResponse: `{"error":"invalid request body"}`,
		},
		{
			name:             "invalid collection ID",
			songID:           validSongID.String(),
			requestBody:      requestBody("bad-uuid"),
			setupMocks:       func() {},
			expectedCode:     http.StatusBadRequest,
			expectedResponse: `{"error":"invalid collection ID"}`,
		},
		{
			name:             "invalid song ID",
			songID:           "bad-uuid",
			requestBody:      requestBody(validCollectionID.String()),
			setupMocks:       func() {},
			expectedCode:     http.StatusBadRequest,
			expectedResponse: `{"error":"invalid song ID"}`,
		},
		{
			name:        "collection song not found error",
			songID:      validSongID.String(),
			requestBody: requestBody(validCollectionID.String()),
			setupMocks: func() {
				mockCollectionSongService.
					EXPECT().
					First(gomock.Any(), &models.CollectionSong{
						CollectionID: validCollectionID,
						SongID:       validSongID,
					}).
					Return(nil, errors.New("not found"))
			},
			expectedCode:     http.StatusInternalServerError,
			expectedResponse: `{"error":"failed to find song in collection"}`,
		},
		{
			name:        "delete failed",
			songID:      validSongID.String(),
			requestBody: requestBody(validCollectionID.String()),
			setupMocks: func() {
				mockCollectionSongService.
					EXPECT().
					First(gomock.Any(), &models.CollectionSong{
						CollectionID: validCollectionID,
						SongID:       validSongID,
					}).
					Return(&models.CollectionSong{
						BaseModel: models.BaseModel{
							ID: collectionSongID,
						},
						CollectionID: validCollectionID,
						SongID:       validSongID,
					}, nil)

				mockCollectionSongService.
					EXPECT().
					Delete(gomock.Any(), collectionSongID).
					Return(errors.New("delete failed"))
			},
			expectedCode:     http.StatusInternalServerError,
			expectedResponse: `{"error":"failed to remove song from collection"}`,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			url := fmt.Sprintf("/songs/%s/collection", tt.songID)
			req := httptest.NewRequest(http.MethodDelete, url, tt.requestBody)
			req = muxSetURLParam(req, "id", tt.songID)

			rr := httptest.NewRecorder()
			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)

			if tt.expectedResponse != "" {
				assert.JSONEq(t, tt.expectedResponse, rr.Body.String())
			}
		})
	}
}
