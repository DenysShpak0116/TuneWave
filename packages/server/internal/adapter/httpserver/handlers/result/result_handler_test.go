package result_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/handlers/helpers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/result"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/service/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResultHandler_SendResult(t *testing.T) {
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

	tests := []struct {
		name           string
		userCtx        bool
		collectionID   string
		requestBody    interface{}
		mockSetup      func(rs *mocks.ResultService)
		expectedStatus int
	}{
		{
			name:         "success",
			userCtx:      true,
			collectionID: collectionID.String(),
			requestBody:  validRequest,
			mockSetup: func(rs *mocks.ResultService) {
				rs.On(
					"ProcessUserResults",
					helpers.CtxMatcher,
					userID,
					collectionID,
					validRequest,
				).Return(nil, nil).Once()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid user id",
			userCtx:        false,
			collectionID:   collectionID.String(),
			requestBody:    validRequest,
			mockSetup:      func(rs *mocks.ResultService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid collection id",
			userCtx:        true,
			collectionID:   "not-a-uuid",
			requestBody:    validRequest,
			mockSetup:      func(rs *mocks.ResultService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid request body",
			userCtx:        true,
			collectionID:   collectionID.String(),
			requestBody:    "invalid-json",
			mockSetup:      func(rs *mocks.ResultService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:         "service error",
			userCtx:      true,
			collectionID: collectionID.String(),
			requestBody:  validRequest,
			mockSetup: func(rs *mocks.ResultService) {
				rs.On(
					"ProcessUserResults",
					helpers.CtxMatcher,
					userID,
					collectionID,
					validRequest,
				).Return(nil, errors.New("service failed")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := mocks.NewResultService(t)
			tt.mockSetup(rs)
			handler := result.NewResultHandler(rs)

			var bodyBytes []byte
			switch v := tt.requestBody.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				var err error
				bodyBytes, err = json.Marshal(v)
				require.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/collections/"+tt.collectionID+"/results", bytes.NewReader(bodyBytes))
			if tt.userCtx {
				req = req.WithContext(helpers.CtxWithUserID(req.Context(), userID))
			}

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.collectionID)
			req = req.WithContext(contextWithChi(req.Context(), rctx))

			rr := httptest.NewRecorder()

			err := handler.SendResult(rr, req)
			if err != nil {
				helpers.AssertAPIError(t, rr, err, tt.expectedStatus)
				return
			}

			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.expectedStatus == http.StatusOK {
				var actual []models.Result
				err := json.Unmarshal(rr.Body.Bytes(), &actual)
				require.NoError(t, err)
			}
		})
	}
}

func contextWithChi(ctx context.Context, rctx *chi.Context) context.Context {
	return context.WithValue(ctx, chi.RouteCtxKey, rctx)
}
