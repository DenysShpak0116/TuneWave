package vector

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/service/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestVectorHandler_GetSongVectors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVectorService := mocks.NewMockVectorService(ctrl)
	mockCollectionSongService := mocks.NewMockCollectionSongService(ctrl)
	mockCriterionService := mocks.NewMockCriterionService(ctrl)
	dtoBuilder := dto.NewDTOBuilder(nil, nil)
	handler := NewVectorHandler(mockVectorService, mockCollectionSongService, mockCriterionService, dtoBuilder)
	httpHandler := handlers.MakeHandler(handler.GetSongVectors)

	colID := uuid.New()
	songID := uuid.New()
	vectorID := uuid.New()
	criterionID := uuid.New()

	tests := []struct {
		name           string
		colParam       string
		songParam      string
		setupMocks     func()
		expectedStatus int
	}{
		{
			name:     "success with vectors",
			colParam: colID.String(), songParam: songID.String(),
			setupMocks: func() {
				mockCollectionSongService.EXPECT().
					First(gomock.Any(), &models.CollectionSong{CollectionID: colID, SongID: songID}, gomock.Any()).
					Return(&models.CollectionSong{
						Vectors: []models.Vector{
							{
								BaseModel:   models.BaseModel{ID: vectorID},
								Mark:        "5",
								CriterionID: criterionID,
								Criterion:   models.Criterion{BaseModel: models.BaseModel{ID: criterionID}, Name: "Quality"},
							},
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "invalid collection id",
			colParam: "bad", songParam: songID.String(),
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "invalid song id",
			colParam: colID.String(), songParam: "bad",
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "service error",
			colParam: colID.String(), songParam: songID.String(),
			setupMocks: func() {
				mockCollectionSongService.EXPECT().
					First(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("fail"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodGet,
				fmt.Sprintf("/collections/%s/%s/vectors", tt.colParam, tt.songParam), nil)
			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("id", tt.colParam)
			ctx.URLParams.Add("song-id", tt.songParam)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			rr := httptest.NewRecorder()
			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			var actual interface{}
			err := json.Unmarshal(rr.Body.Bytes(), &actual)
			assert.NoError(t, err)
		})
	}
}

func TestVectorHandler_CreateSongVectors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVectorService := mocks.NewMockVectorService(ctrl)
	mockCollectionSongService := mocks.NewMockCollectionSongService(ctrl)
	mockCriterionService := mocks.NewMockCriterionService(ctrl)
	dtoBuilder := dto.NewDTOBuilder(nil, nil)
	handler := NewVectorHandler(mockVectorService, mockCollectionSongService, mockCriterionService, dtoBuilder)
	httpHandler := handlers.MakeHandler(handler.CreateSongVectors)

	colID := uuid.New()
	songID := uuid.New()
	csID := uuid.New()
	vecID := uuid.New()
	critID := uuid.New()

	requestBody := fmt.Sprintf(`{"vectors":[{"criterionId":"%s","mark":"10"}]}`, critID.String())

	tests := []struct {
		name           string
		colParam       string
		songParam      string
		body           string
		setupMocks     func()
		expectedStatus int
	}{
		{
			name:     "success create",
			colParam: colID.String(), songParam: songID.String(),
			body: requestBody,
			setupMocks: func() {
				mockCollectionSongService.EXPECT().
					First(gomock.Any(), &models.CollectionSong{CollectionID: colID, SongID: songID}).
					Return(&models.CollectionSong{BaseModel: models.BaseModel{ID: csID}}, nil)
				mockVectorService.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil)
				mockCollectionSongService.EXPECT().
					GetByID(gomock.Any(), csID, gomock.Any()).
					Return(&models.CollectionSong{Vectors: []models.Vector{
						{
							BaseModel:   models.BaseModel{ID: vecID},
							Mark:        "10",
							CriterionID: critID,
							Criterion: models.Criterion{
								BaseModel: models.BaseModel{ID: critID},
								Name:      "X",
							},
						},
					}}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:     "bad collection id",
			colParam: "bad", songParam: songID.String(),
			body:           requestBody,
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "bad song id",
			colParam: colID.String(), songParam: "bad",
			body:           requestBody,
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "invalid payload",
			colParam: colID.String(), songParam: songID.String(),
			body:           `not-json`,
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "service First failure",
			colParam: colID.String(), songParam: songID.String(),
			body: requestBody,
			setupMocks: func() {
				mockCollectionSongService.EXPECT().
					First(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("fail"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:     "vector create failure",
			colParam: colID.String(), songParam: songID.String(),
			body: requestBody,
			setupMocks: func() {
				mockCollectionSongService.EXPECT().
					First(gomock.Any(), gomock.Any()).
					Return(&models.CollectionSong{BaseModel: models.BaseModel{ID: csID}}, nil)
				mockVectorService.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(errors.New("fail"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodPost,
				fmt.Sprintf("/collections/%s/%s/vectors", tt.colParam, tt.songParam),
				strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("id", tt.colParam)
			ctx.URLParams.Add("song-id", tt.songParam)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			rr := httptest.NewRecorder()
			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			var actual interface{}
			json.Unmarshal(rr.Body.Bytes(), &actual)
		})
	}
}

func TestVectorHandler_UpdateSongVectors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVectorService := mocks.NewMockVectorService(ctrl)
	mockCollectionSongService := mocks.NewMockCollectionSongService(ctrl)
	mockCriterionService := mocks.NewMockCriterionService(ctrl)
	dtoBuilder := dto.NewDTOBuilder(nil, nil)

	handler := NewVectorHandler(mockVectorService, mockCollectionSongService, mockCriterionService, dtoBuilder)
	httpHandler := handlers.MakeHandler(handler.UpdateSongVectors)

	colID := uuid.New()
	songID := uuid.New()
	vecID := uuid.New()
	critID := uuid.New()

	body := fmt.Sprintf(`{"vectors":[{"id":"%s","criterionId":"%s","mark":"7"}]}`, vecID, critID)

	tests := []struct {
		name           string
		colParam       string
		songParam      string
		body           string
		setupMocks     func()
		expectedStatus int
	}{
		{
			name:     "success update",
			colParam: colID.String(), songParam: songID.String(),
			body: body,
			setupMocks: func() {
				mockVectorService.EXPECT().
					Update(gomock.Any(), &models.Vector{
						BaseModel:   models.BaseModel{ID: vecID},
						Mark:        "7",
						CriterionID: critID,
					}).
					Return(nil)

				mockCollectionSongService.EXPECT().
					First(gomock.Any(),
						&models.CollectionSong{CollectionID: colID, SongID: songID},
						"Vectors", "Vectors.Criterion").
					Return(&models.CollectionSong{
						Vectors: []models.Vector{
							{
								BaseModel:   models.BaseModel{ID: vecID},
								Mark:        "7",
								CriterionID: critID,
								Criterion: models.Criterion{
									BaseModel: models.BaseModel{ID: critID},
									Name:      "Lyrical",
								},
							},
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "invalid collection id",
			colParam: "bad", songParam: songID.String(),
			body:           body,
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "invalid song id",
			colParam: colID.String(), songParam: "bad",
			body:           body,
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "invalid JSON",
			colParam: colID.String(), songParam: songID.String(),
			body:           `{"vectors":bad}`,
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "update failure",
			colParam: colID.String(), songParam: songID.String(),
			body: body,
			setupMocks: func() {
				mockVectorService.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(errors.New("fail"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:     "get collection song fails",
			colParam: colID.String(), songParam: songID.String(),
			body: body,
			setupMocks: func() {
				mockVectorService.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil)
				mockCollectionSongService.EXPECT().
					First(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("fail"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodPut,
				fmt.Sprintf("/collections/%s/%s/vectors", tt.colParam, tt.songParam),
				strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("id", tt.colParam)
			ctx.URLParams.Add("song-id", tt.songParam)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			rr := httptest.NewRecorder()
			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			var actual interface{}
			json.Unmarshal(rr.Body.Bytes(), &actual)
		})
	}
}

func TestVectorHandler_DeleteSongVectors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVectorService := mocks.NewMockVectorService(ctrl)
	mockCollectionSongService := mocks.NewMockCollectionSongService(ctrl)
	mockCriterionService := mocks.NewMockCriterionService(ctrl)
	dtoBuilder := dto.NewDTOBuilder(nil, nil)

	handler := NewVectorHandler(mockVectorService, mockCollectionSongService, mockCriterionService, dtoBuilder)
	httpHandler := handlers.MakeHandler(handler.DeleteSongVectors)

	colID := uuid.New()
	songID := uuid.New()
	vecID := uuid.New()

	tests := []struct {
		name           string
		colParam       string
		songParam      string
		setupMocks     func()
		expectedStatus int
	}{
		{
			name:     "success delete",
			colParam: colID.String(), songParam: songID.String(),
			setupMocks: func() {
				mockCollectionSongService.EXPECT().
					First(gomock.Any(),
						&models.CollectionSong{CollectionID: colID, SongID: songID},
						"Vectors").
					Return(&models.CollectionSong{
						Vectors: []models.Vector{{BaseModel: models.BaseModel{ID: vecID}}},
					}, nil)
				mockVectorService.EXPECT().
					Delete(gomock.Any(), vecID).
					Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:     "invalid collection id",
			colParam: "bad", songParam: songID.String(),
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "invalid song id",
			colParam: colID.String(), songParam: "bad",
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "First error",
			colParam: colID.String(), songParam: songID.String(),
			setupMocks: func() {
				mockCollectionSongService.EXPECT().
					First(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("fail"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:     "Delete error",
			colParam: colID.String(), songParam: songID.String(),
			setupMocks: func() {
				mockCollectionSongService.EXPECT().
					First(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&models.CollectionSong{
						Vectors: []models.Vector{{BaseModel: models.BaseModel{ID: vecID}}},
					}, nil)
				mockVectorService.EXPECT().
					Delete(gomock.Any(), vecID).
					Return(errors.New("fail"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodDelete,
				fmt.Sprintf("/collections/%s/%s/vectors", tt.colParam, tt.songParam),
				nil)
			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("id", tt.colParam)
			ctx.URLParams.Add("song-id", tt.songParam)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			rr := httptest.NewRecorder()
			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestVectorHandler_HasAllVectors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVectorService := mocks.NewMockVectorService(ctrl)
	mockCollectionSongService := mocks.NewMockCollectionSongService(ctrl)
	mockCriterionService := mocks.NewMockCriterionService(ctrl)
	dtoBuilder := dto.NewDTOBuilder(nil, nil)

	handler := NewVectorHandler(mockVectorService, mockCollectionSongService, mockCriterionService, dtoBuilder)
	httpHandler := handlers.MakeHandler(handler.HasAllVectors)

	colID := uuid.New()

	tests := []struct {
		name           string
		colParam       string
		setupMocks     func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:     "all vectors present",
			colParam: colID.String(),
			setupMocks: func() {
				mockCriterionService.EXPECT().
					CountWhere(gomock.Any(), gomock.Any()).
					Return(int64(2), nil)

				mockCollectionSongService.EXPECT().
					Where(gomock.Any(),
						&models.CollectionSong{CollectionID: colID},
						gomock.Any()).
					Return([]models.CollectionSong{
						{Vectors: []models.Vector{{}, {}}},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]bool{"hasAllVectors": true},
		},
		{
			name:     "partial vectors",
			colParam: colID.String(),
			setupMocks: func() {
				mockCriterionService.EXPECT().CountWhere(gomock.Any(), gomock.Any()).Return(int64(3), nil)
				mockCollectionSongService.EXPECT().
					Where(gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]models.CollectionSong{{Vectors: []models.Vector{{}}}}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]bool{"hasAllVectors": false},
		},
		{
			name:           "invalid id",
			colParam:       "bad",
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"error": "invalid collection id"},
		},
		{
			name:     "no criteria",
			colParam: colID.String(),
			setupMocks: func() {
				mockCriterionService.EXPECT().
					CountWhere(gomock.Any(), gomock.Any()).
					Return(int64(0), nil)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   map[string]interface{}{"error": "there is no criteria"},
		},
		{
			name:     "no collection songs",
			colParam: colID.String(),
			setupMocks: func() {
				mockCriterionService.EXPECT().CountWhere(gomock.Any(), gomock.Any()).Return(int64(2), nil)
				mockCollectionSongService.EXPECT().
					Where(gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]models.CollectionSong{}, nil)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   map[string]interface{}{"error": "there is no collection song"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/collections/%s/has-all-vectors", tt.colParam), nil)
			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("id", tt.colParam)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			rr := httptest.NewRecorder()
			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			var actual map[string]bool
			if tt.expectedStatus == http.StatusOK {
				err := json.Unmarshal(rr.Body.Bytes(), &actual)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedBody, actual)
			}
		})
	}
}
