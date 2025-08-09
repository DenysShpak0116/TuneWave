package song

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/service/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSongHandler_GetSongs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSongService := mocks.NewMockSongService(ctrl)
	mockCollectionSongService := mocks.NewMockCollectionSongService(ctrl)
	mockUserReactionService := mocks.NewMockUserReactionService(ctrl)
	mockCommentService := mocks.NewMockCommentService(ctrl)
	mockUserService := mocks.NewMockUserService(ctrl)

	dtoBuilder := dto.NewDTOBuilder(mockUserService, mockUserReactionService)

	handler := NewSongHandler(
		mockSongService,
		mockCollectionSongService,
		mockUserReactionService,
		mockCommentService,
		dtoBuilder,
	)
	httpHandler := handlers.MakeHandler(handler.GetSongs)

	mockSongs := []models.Song{
		{
			BaseModel: models.BaseModel{ID: uuid.New()},
			Title:     "Test Song",
		},
	}

	testCases := []struct {
		name         string
		query        string
		expectedCode int
		setupMocks   func()
	}{
		{
			name:         "success",
			query:        "?search=rock&sortBy=name&order=asc&page=1&limit=5",
			expectedCode: http.StatusOK,
			setupMocks: func() {
				expectedParams := services.SearchSongsParams{
					Search: "rock",
					SortBy: "name",
					Order:  "asc",
					Page:   1,
					Limit:  5,
				}
				mockSongService.EXPECT().
					GetSongs(gomock.Any(), expectedParams, "Authors", "Authors.Author").
					Return(mockSongs, nil)

				mockUserReactionService.EXPECT().GetSongLikes(gomock.Any(), gomock.Any()).Return(int64(0))
				mockUserReactionService.EXPECT().GetSongDislikes(gomock.Any(), gomock.Any()).Return(int64(0))
			},
		},
		{
			name:         "invalid page and limit - fallback to defaults",
			query:        "?page=abc&limit=-1",
			expectedCode: http.StatusOK,
			setupMocks: func() {
				expectedParams := services.SearchSongsParams{
					Search: "",
					SortBy: "created_at",
					Order:  "desc",
					Page:   1,
					Limit:  10,
				}
				mockSongService.EXPECT().
					GetSongs(gomock.Any(), expectedParams, "Authors", "Authors.Author").
					Return(mockSongs, nil)
				mockUserReactionService.EXPECT().GetSongLikes(gomock.Any(), gomock.Any()).Return(int64(0))
				mockUserReactionService.EXPECT().GetSongDislikes(gomock.Any(), gomock.Any()).Return(int64(0))
			},
		},
		{
			name:         "internal error from service",
			query:        "",
			expectedCode: http.StatusInternalServerError,
			setupMocks: func() {
				expectedParams := services.SearchSongsParams{
					Search: "",
					SortBy: "created_at",
					Order:  "desc",
					Page:   1,
					Limit:  10,
				}
				mockSongService.EXPECT().
					GetSongs(gomock.Any(), expectedParams, "Authors", "Authors.Author").
					Return(nil, errors.New("db error"))
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodGet, "/songs"+tt.query, nil)
			rr := httptest.NewRecorder()

			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)

			if tt.expectedCode == http.StatusOK {
				var response []dto.SongPreviewDTO
				err := json.NewDecoder(rr.Body).Decode(&response)
				assert.NoError(t, err)
			}
		})
	}
}

func TestSongHandler_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSongService := mocks.NewMockSongService(ctrl)
	mockCollectionSongService := mocks.NewMockCollectionSongService(ctrl)
	mockUserReactionService := mocks.NewMockUserReactionService(ctrl)
	mockCommentService := mocks.NewMockCommentService(ctrl)
	mockUserService := mocks.NewMockUserService(ctrl)

	dtoBuilder := dto.NewDTOBuilder(mockUserService, mockUserReactionService)

	handler := NewSongHandler(
		mockSongService,
		mockCollectionSongService,
		mockUserReactionService,
		mockCommentService,
		dtoBuilder,
	)
	httpHandler := handlers.MakeHandler(handler.GetByID)

	songID := uuid.New()
	mockSong := &models.Song{
		BaseModel: models.BaseModel{ID: songID},
		Title:     "Test Song",
	}

	testCases := []struct {
		name         string
		songID       string
		expectedCode int
		setupMocks   func()
	}{
		{
			name:         "success",
			songID:       songID.String(),
			expectedCode: http.StatusOK,
			setupMocks: func() {
				mockSongService.EXPECT().
					GetByID(gomock.Any(), songID, gomock.Any()).
					Return(mockSong, nil)
				mockUserReactionService.EXPECT().GetSongLikes(gomock.Any(), gomock.Any()).Return(int64(0))
				mockUserReactionService.EXPECT().GetSongDislikes(gomock.Any(), gomock.Any()).Return(int64(0))
				mockUserService.EXPECT().GetUserFollowersCount(gomock.Any(), gomock.Any()).Return(int64(0))
			},
		},
		{
			name:         "invalid song UUID",
			songID:       "invalid-uuid",
			expectedCode: http.StatusBadRequest,
			setupMocks:   func() {},
		},
		{
			name:         "song not found / service error",
			songID:       songID.String(),
			expectedCode: http.StatusInternalServerError,
			setupMocks: func() {
				mockSongService.EXPECT().
					GetByID(gomock.Any(), songID, gomock.Any()).
					Return(nil, errors.New("not found"))
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodGet, "/songs/"+tt.songID, nil)
			rr := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.songID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)

			if tt.expectedCode == http.StatusOK {
				var resp dto.SongDTO
				err := json.NewDecoder(rr.Body).Decode(&resp)
				assert.NoError(t, err)
			}
		})
	}
}

func TestSongHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSongService := mocks.NewMockSongService(ctrl)
	mockCollectionSongService := mocks.NewMockCollectionSongService(ctrl)
	mockUserReactionService := mocks.NewMockUserReactionService(ctrl)
	mockCommentService := mocks.NewMockCommentService(ctrl)
	mockUserService := mocks.NewMockUserService(ctrl)

	dtoBuilder := dto.NewDTOBuilder(mockUserService, mockUserReactionService)

	handler := NewSongHandler(
		mockSongService,
		mockCollectionSongService,
		mockUserReactionService,
		mockCommentService,
		dtoBuilder,
	)
	httpHandler := handlers.MakeHandler(handler.Create)

	userID := uuid.New()
	songID := uuid.New()

	testCases := []struct {
		name         string
		expectedCode int
		setupMocks   func()
		buildRequest func() *http.Request
	}{
		{
			name:         "success",
			expectedCode: http.StatusCreated,
			setupMocks: func() {
				mockSongService.EXPECT().
					SaveSong(gomock.Any(), gomock.Any()).
					Return(&models.Song{BaseModel: models.BaseModel{ID: songID}}, nil)

				mockSongService.EXPECT().
					GetByID(gomock.Any(), songID).
					Return(&models.Song{BaseModel: models.BaseModel{ID: songID}}, nil)

				mockUserReactionService.EXPECT().
					GetSongLikes(gomock.Any(), gomock.Any()).
					Return(int64(0))

				mockUserReactionService.EXPECT().
					GetSongDislikes(gomock.Any(), gomock.Any()).
					Return(int64(0))

				mockUserService.EXPECT().
					GetUserFollowersCount(gomock.Any(), gomock.Any()).
					Return(int64(0))
			},
			buildRequest: func() *http.Request {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)

				_ = writer.WriteField("userId", userID.String())
				_ = writer.WriteField("title", "Test Song")
				_ = writer.WriteField("genre", "Pop")
				_ = writer.WriteField("artists", "Artist1")
				_ = writer.WriteField("tags", "tag1")

				songWriter, _ := writer.CreateFormFile("song", "song.mp3")
				_, _ = songWriter.Write([]byte("dummy song data"))

				coverWriter, _ := writer.CreateFormFile("cover", "cover.jpg")
				_, _ = coverWriter.Write([]byte("dummy cover data"))

				_ = writer.Close()

				req := httptest.NewRequest(http.MethodPost, "/songs", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req
			},
		},
		{
			name:         "invalid user ID",
			expectedCode: http.StatusBadRequest,
			setupMocks:   func() {},
			buildRequest: func() *http.Request {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				_ = writer.WriteField("userId", "not-a-uuid")
				_ = writer.Close()
				req := httptest.NewRequest(http.MethodPost, "/songs", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req
			},
		},
		{
			name:         "missing song file",
			expectedCode: http.StatusBadRequest,
			setupMocks:   func() {},
			buildRequest: func() *http.Request {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				_ = writer.WriteField("userId", userID.String())
				coverWriter, _ := writer.CreateFormFile("cover", "cover.jpg")
				_, _ = coverWriter.Write([]byte("dummy cover data"))
				_ = writer.Close()
				req := httptest.NewRequest(http.MethodPost, "/songs", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req
			},
		},
		{
			name:         "save song failed",
			expectedCode: http.StatusInternalServerError,
			setupMocks: func() {
				mockSongService.EXPECT().
					SaveSong(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("save error"))
			},
			buildRequest: func() *http.Request {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				_ = writer.WriteField("userId", userID.String())
				_ = writer.WriteField("title", "Test Song")
				_ = writer.WriteField("genre", "Pop")

				songWriter, _ := writer.CreateFormFile("song", "song.mp3")
				_, _ = songWriter.Write([]byte("dummy song"))

				coverWriter, _ := writer.CreateFormFile("cover", "cover.jpg")
				_, _ = coverWriter.Write([]byte("dummy cover"))

				_ = writer.Close()
				req := httptest.NewRequest(http.MethodPost, "/songs", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			req := tt.buildRequest()
			rr := httptest.NewRecorder()

			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)

			if rr.Code == http.StatusCreated {
				var song dto.SongDTO
				err := json.NewDecoder(rr.Body).Decode(&song)
				assert.NoError(t, err)
				assert.Equal(t, songID, song.ID)
			}
		})
	}
}

func TestSongHandler_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSongService := mocks.NewMockSongService(ctrl)
	mockCollectionSongService := mocks.NewMockCollectionSongService(ctrl)
	mockUserReactionService := mocks.NewMockUserReactionService(ctrl)
	mockCommentService := mocks.NewMockCommentService(ctrl)
	mockUserService := mocks.NewMockUserService(ctrl)

	dtoBuilder := dto.NewDTOBuilder(mockUserService, mockUserReactionService)

	handler := NewSongHandler(
		mockSongService,
		mockCollectionSongService,
		mockUserReactionService,
		mockCommentService,
		dtoBuilder,
	)
	httpHandler := handlers.MakeHandler(handler.Update)

	songID := uuid.New()

	testCases := []struct {
		name         string
		songID       string
		expectedCode int
		setupMocks   func()
		buildRequest func() *http.Request
	}{
		{
			name:         "success",
			songID:       songID.String(),
			expectedCode: http.StatusOK,
			setupMocks: func() {
				mockSongService.EXPECT().
					UpdateSong(gomock.Any(), gomock.Any()).
					Return(nil)

				mockSongService.EXPECT().
					GetByID(gomock.Any(), songID).
					Return(&models.Song{BaseModel: models.BaseModel{ID: songID}}, nil)

				mockUserReactionService.EXPECT().
					GetSongLikes(gomock.Any(), gomock.Any()).
					Return(int64(0))

				mockUserReactionService.EXPECT().
					GetSongDislikes(gomock.Any(), gomock.Any()).
					Return(int64(0))

				mockUserService.EXPECT().
					GetUserFollowersCount(gomock.Any(), gomock.Any()).
					Return(int64(0))
			},
			buildRequest: func() *http.Request {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)

				_ = writer.WriteField("title", "Updated Title")
				_ = writer.WriteField("genre", "Rock")
				_ = writer.WriteField("artists", "Artist1")
				_ = writer.WriteField("tags", "tag1")

				songWriter, _ := writer.CreateFormFile("song", "updated.mp3")
				_, _ = songWriter.Write([]byte("updated song"))

				coverWriter, _ := writer.CreateFormFile("cover", "updated.jpg")
				_, _ = coverWriter.Write([]byte("updated cover"))

				writer.Close()

				req := httptest.NewRequest(http.MethodPut, "/songs/"+songID.String(), body)
				req = muxSetURLParam(req, "id", songID.String()) // mimic chi.URLParam
				req.Header.Set("Content-Type", writer.FormDataContentType())

				return req
			},
		},
		{
			name:         "invalid song ID",
			songID:       "invalid-uuid",
			expectedCode: http.StatusBadRequest,
			setupMocks:   func() {},
			buildRequest: func() *http.Request {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				_ = writer.Close()
				req := httptest.NewRequest(http.MethodPut, "/songs/invalid-uuid", body)
				req = muxSetURLParam(req, "id", "invalid-uuid")
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req
			},
		},
		{
			name:         "update failed",
			songID:       songID.String(),
			expectedCode: http.StatusInternalServerError,
			setupMocks: func() {
				mockSongService.EXPECT().
					UpdateSong(gomock.Any(), gomock.Any()).
					Return(errors.New("update error"))
			},
			buildRequest: func() *http.Request {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				_ = writer.WriteField("title", "Test")
				_ = writer.Close()

				req := httptest.NewRequest(http.MethodPut, "/songs/"+songID.String(), body)
				req = muxSetURLParam(req, "id", songID.String())
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req
			},
		},
		{
			name:         "get after update failed",
			songID:       songID.String(),
			expectedCode: http.StatusInternalServerError,
			setupMocks: func() {
				mockSongService.EXPECT().
					UpdateSong(gomock.Any(), gomock.Any()).
					Return(nil)

				mockSongService.EXPECT().
					GetByID(gomock.Any(), songID).
					Return(nil, errors.New("get failed"))
			},
			buildRequest: func() *http.Request {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				_ = writer.WriteField("title", "Test")
				_ = writer.Close()

				req := httptest.NewRequest(http.MethodPut, "/songs/"+songID.String(), body)
				req = muxSetURLParam(req, "id", songID.String())
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			req := tt.buildRequest()
			rr := httptest.NewRecorder()
			httpHandler.ServeHTTP(rr, req)
			assert.Equal(t, tt.expectedCode, rr.Code)
		})
	}
}

func TestSongHandler_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSongService := mocks.NewMockSongService(ctrl)
	mockCollectionSongService := mocks.NewMockCollectionSongService(ctrl)
	mockUserReactionService := mocks.NewMockUserReactionService(ctrl)
	mockCommentService := mocks.NewMockCommentService(ctrl)
	mockUserService := mocks.NewMockUserService(ctrl)

	dtoBuilder := dto.NewDTOBuilder(mockUserService, mockUserReactionService)

	handler := NewSongHandler(
		mockSongService,
		mockCollectionSongService,
		mockUserReactionService,
		mockCommentService,
		dtoBuilder,
	)

	httpHandler := handlers.MakeHandler(handler.Delete)

	songID := uuid.New()

	testCases := []struct {
		name         string
		songID       string
		expectedCode int
		setupMocks   func()
	}{
		{
			name:         "success",
			songID:       songID.String(),
			expectedCode: http.StatusNoContent,
			setupMocks: func() {
				mockSongService.
					EXPECT().
					Delete(gomock.Any(), songID).
					Return(nil)
			},
		},
		{
			name:         "invalid song ID",
			songID:       "invalid-uuid",
			expectedCode: http.StatusBadRequest,
			setupMocks:   func() {},
		},
		{
			name:         "delete failure",
			songID:       songID.String(),
			expectedCode: http.StatusInternalServerError,
			setupMocks: func() {
				mockSongService.
					EXPECT().
					Delete(gomock.Any(), songID).
					Return(errors.New("delete failed"))
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodDelete, "/songs/"+tt.songID, nil)
			req = muxSetURLParam(req, "id", tt.songID)
			rr := httptest.NewRecorder()

			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
		})
	}
}

func TestSongHandler_GetGenres(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSongService := mocks.NewMockSongService(ctrl)
	mockCollectionSongService := mocks.NewMockCollectionSongService(ctrl)
	mockUserReactionService := mocks.NewMockUserReactionService(ctrl)
	mockCommentService := mocks.NewMockCommentService(ctrl)
	mockUserService := mocks.NewMockUserService(ctrl)

	dtoBuilder := dto.NewDTOBuilder(mockUserService, mockUserReactionService)

	handler := NewSongHandler(
		mockSongService,
		mockCollectionSongService,
		mockUserReactionService,
		mockCommentService,
		dtoBuilder,
	)

	httpHandler := handlers.MakeHandler(handler.GetGenres)

	testCases := []struct {
		name             string
		setupMocks       func()
		expectedCode     int
		expectedResponse string
	}{
		{
			name: "genres with songs",
			setupMocks: func() {
				mockSongService.
					EXPECT().
					GetGenres(gomock.Any()).
					Return([]string{"rock", "jazz"})

				mockSongService.
					EXPECT().
					GetGenresMostPopularSong(gomock.Any(), "rock").
					Return(&models.Song{CoverURL: "rock.jpg"}, nil)

				mockSongService.
					EXPECT().
					GetGenresMostPopularSong(gomock.Any(), "jazz").
					Return(&models.Song{CoverURL: "jazz.png"}, nil)
			},
			expectedCode: http.StatusOK,
			expectedResponse: `[{"genreName":"rock","genreCover":"rock.jpg"},` +
				`{"genreName":"jazz","genreCover":"jazz.png"}]`,
		},
		{
			name: "genres with one error",
			setupMocks: func() {
				mockSongService.
					EXPECT().
					GetGenres(gomock.Any()).
					Return([]string{"rock", "metal"})

				mockSongService.
					EXPECT().
					GetGenresMostPopularSong(gomock.Any(), "rock").
					Return(nil, errors.New("error"))

				mockSongService.
					EXPECT().
					GetGenresMostPopularSong(gomock.Any(), "metal").
					Return(&models.Song{CoverURL: "metal.jpg"}, nil)
			},
			expectedCode: http.StatusOK,
			expectedResponse: `[{"genreName":"rock","genreCover":""},` +
				`{"genreName":"metal","genreCover":"metal.jpg"}]`,
		},
		{
			name: "empty genres list",
			setupMocks: func() {
				mockSongService.
					EXPECT().
					GetGenres(gomock.Any()).
					Return([]string{})
			},
			expectedCode:     http.StatusOK,
			expectedResponse: `[]`,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodGet, "/genres", nil)
			rr := httptest.NewRecorder()

			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			assert.JSONEq(t, tt.expectedResponse, rr.Body.String())
		})
	}
}

func TestSongHandler_GetSongComments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSongService := mocks.NewMockSongService(ctrl)
	mockCollectionSongService := mocks.NewMockCollectionSongService(ctrl)
	mockUserReactionService := mocks.NewMockUserReactionService(ctrl)
	mockCommentService := mocks.NewMockCommentService(ctrl)
	mockUserService := mocks.NewMockUserService(ctrl)

	dtoBuilder := dto.NewDTOBuilder(mockUserService, mockUserReactionService)

	handler := NewSongHandler(
		mockSongService,
		mockCollectionSongService,
		mockUserReactionService,
		mockCommentService,
		dtoBuilder,
	)

	httpHandler := handlers.MakeHandler(handler.GetSongComments)

	validSongID := uuid.New()
	validUserID := uuid.New()
	createdAt := time.Now()

	testCases := []struct {
		name         string
		songID       string
		queryParams  string
		setupMocks   func()
		expectedCode int
	}{
		{
			name:        "success",
			songID:      validSongID.String(),
			queryParams: "?page=1&limit=10",
			setupMocks: func() {
				mockCommentService.
					EXPECT().
					Where(gomock.Any(), &models.Comment{SongID: validSongID}, gomock.Any(), gomock.Any()).
					Return([]models.Comment{
						{
							BaseModel: models.BaseModel{
								ID:        uuid.New(),
								CreatedAt: createdAt,
							},
							Content: "Great song!",
							UserID:  validUserID,
							SongID:  validSongID,
							User: models.User{
								BaseModel: models.BaseModel{
									ID: validUserID,
								},
								Username:       "testuser",
								Role:           "listener",
								ProfilePicture: "pic.jpg",
								ProfileInfo:    "some info",
							},
						},
					}, nil)

				mockUserService.
					EXPECT().
					GetUserFollowersCount(gomock.Any(), validUserID).
					Return(int64(100))
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid song ID",
			songID:       "invalid-uuid",
			queryParams:  "",
			setupMocks:   func() {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "comment service error",
			songID:      validSongID.String(),
			queryParams: "",
			setupMocks: func() {
				mockCommentService.
					EXPECT().
					Where(gomock.Any(), &models.Comment{SongID: validSongID}, gomock.Any(), gomock.Any()).
					Return(nil, errors.New("db failure"))
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:        "invalid pagination values",
			songID:      validSongID.String(),
			queryParams: "?page=abc&limit=xyz",
			setupMocks: func() {
				mockCommentService.
					EXPECT().
					Where(gomock.Any(), &models.Comment{SongID: validSongID}, gomock.Any(), gomock.Any()).
					Return([]models.Comment{}, nil)
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			url := fmt.Sprintf("/songs/%s/comments%s", tt.songID, tt.queryParams)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req = muxSetURLParam(req, "id", tt.songID)

			rr := httptest.NewRecorder()
			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
		})
	}
}

func muxSetURLParam(req *http.Request, key, val string) *http.Request {
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
		URLParams: chi.RouteParams{
			Keys:   []string{key},
			Values: []string{val},
		},
	}))
}
