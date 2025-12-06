package criterion

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/service/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCriterionHandler_CreateCriterion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCriterionService := mocks.NewMockCriterionService(ctrl)
	handler := NewCriterionHandler(mockCriterionService)
	httpHandler := handlers.MakeHandler(handler.CreateCriterion)

	tests := []struct {
		name           string
		body           string
		setupMocks     func()
		expectedStatus int
	}{
		{
			name: "success",
			body: `{"name": "Quality"}`,
			setupMocks: func() {
				mockCriterionService.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid request body",
			body:           `invalid-json`,
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "internal error",
			body: `{"name": "Speed"}`,
			setupMocks: func() {
				mockCriterionService.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			req := httptest.NewRequest(http.MethodPost, "/criterions", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestCriterionHandler_GetCriterions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCriterionService := mocks.NewMockCriterionService(ctrl)
	handler := NewCriterionHandler(mockCriterionService)
	httpHandler := handlers.MakeHandler(handler.GetCriterions)

	expected := []models.Criterion{{Name: "Accuracy"}, {Name: "Relevance"}}

	tests := []struct {
		name           string
		setupMocks     func()
		expectedStatus int
	}{
		{
			name: "success",
			setupMocks: func() {
				mockCriterionService.EXPECT().
					Where(gomock.Any(), gomock.Any()).
					Return(expected, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "internal error",
			setupMocks: func() {
				mockCriterionService.EXPECT().
					Where(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("db failure"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodGet, "/criterions", nil)
			rr := httptest.NewRecorder()

			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestCriterionHandler_UpdateCriterion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCriterionService := mocks.NewMockCriterionService(ctrl)
	handler := NewCriterionHandler(mockCriterionService)
	httpHandler := handlers.MakeHandler(handler.UpdateCriterion)

	id := uuid.New()

	tests := []struct {
		name           string
		idParam        string
		body           string
		setupMocks     func()
		expectedStatus int
	}{
		{
			name:    "success",
			idParam: id.String(),
			body:    `{"name": "Updated Criterion"}`,
			setupMocks: func() {
				mockCriterionService.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid UUID",
			idParam:        "invalid-uuid",
			body:           `{}`,
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid JSON",
			idParam:        id.String(),
			body:           `invalid`,
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "update failed",
			idParam: id.String(),
			body:    `{"name": "Try"}`,
			setupMocks: func() {
				mockCriterionService.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodPut, "/criterions/"+tt.idParam, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("id", tt.idParam)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			rr := httptest.NewRecorder()
			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestCriterionHandler_DeleteCriterion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCriterionService := mocks.NewMockCriterionService(ctrl)
	handler := NewCriterionHandler(mockCriterionService)
	httpHandler := handlers.MakeHandler(handler.DeleteCriterion)

	id := uuid.New()

	tests := []struct {
		name           string
		idParam        string
		setupMocks     func()
		expectedStatus int
	}{
		{
			name:    "success",
			idParam: id.String(),
			setupMocks: func() {
				mockCriterionService.EXPECT().
					Delete(gomock.Any(), id).
					Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "invalid UUID",
			idParam:        "invalid-id",
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "delete failed",
			idParam: id.String(),
			setupMocks: func() {
				mockCriterionService.EXPECT().
					Delete(gomock.Any(), id).
					Return(errors.New("fail"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodDelete, "/criterions/"+tt.idParam, nil)
			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("id", tt.idParam)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			rr := httptest.NewRecorder()
			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
