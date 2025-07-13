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
