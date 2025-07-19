package result

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

func TestResultHandler_SendResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockResultService := mocks.NewMockResultService(ctrl)
	mockCollectionSongService := mocks.NewMockCollectionSongService(ctrl)

	handler := NewResultHandler(mockResultService, mockCollectionSongService)
	httpHandler := handlers.MakeHandler(handler.SendResult)

	userID := uuid.New()
	collectionID := uuid.New()

	validRequest := dto.SendResultRequest{
		Results: []dto.ResultRequest{
			{
				Song1ID: uuid.NewString(),
				ComparedTo: []struct {
					Song2ID string `json:"song2Id"`
					Result  int    `json:"result"`
				}{
					{Song2ID: uuid.NewString(), Result: 1},
					{Song2ID: uuid.NewString(), Result: -1},
				},
			},
		},
	}

	expectedResults := []models.Result{
		{
			BaseModel:        models.BaseModel{ID: uuid.New()},
			SongRank:         1,
			UserID:           userID,
			CollectionSongID: uuid.New(),
		},
		{
			BaseModel:        models.BaseModel{ID: uuid.New()},
			SongRank:         2,
			UserID:           userID,
			CollectionSongID: uuid.New(),
		},
	}

	tests := []struct {
		name           string
		userCtx        bool
		collectionID   string
		requestBody    interface{}
		mockSetup      func()
		expectedStatus int
	}{
		{
			name:         "success",
			userCtx:      true,
			collectionID: collectionID.String(),
			requestBody:  validRequest,
			mockSetup: func() {
				mockResultService.EXPECT().
					ProcessUserResults(gomock.Any(), userID, collectionID, validRequest).
					Return(expectedResults, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid user id",
			userCtx:        false,
			collectionID:   collectionID.String(),
			requestBody:    validRequest,
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid collection id",
			userCtx:        true,
			collectionID:   "not-a-uuid",
			requestBody:    validRequest,
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid request body",
			userCtx:        true,
			collectionID:   collectionID.String(),
			requestBody:    "invalid-json",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:         "service error",
			userCtx:      true,
			collectionID: collectionID.String(),
			requestBody:  validRequest,
			mockSetup: func() {
				mockResultService.EXPECT().
					ProcessUserResults(gomock.Any(), userID, collectionID, validRequest).
					Return(nil, errors.New("service failed"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			var body []byte
			switch v := tt.requestBody.(type) {
			case string:
				body = []byte(v)
			default:
				b, _ := json.Marshal(v)
				body = b
			}

			req := httptest.NewRequest(http.MethodPost, "/collections/"+tt.collectionID+"/results", bytes.NewReader(body))
			if tt.userCtx {
				req = req.WithContext(context.WithValue(req.Context(), "userID", userID.String()))
			}

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.collectionID)
			req = req.WithContext(contextWithChi(req.Context(), rctx))

			rr := httptest.NewRecorder()
			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func contextWithChi(ctx context.Context, rctx *chi.Context) context.Context {
	return context.WithValue(ctx, chi.RouteCtxKey, rctx)
}
