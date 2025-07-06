package collection

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
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

func TestCollectionHandler_CreateCollection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionService := mocks.NewMockCollectionService(ctrl)
	mockUserCollectionService := mocks.NewMockUserCollectionService(ctrl)
	mockUserReactionService := mocks.NewMockUserReactionService(ctrl)
	mockUserService := mocks.NewMockUserService(ctrl)

	dtoBuilder := dto.NewDTOBuilder(mockUserService, nil)
	handler := NewCollectionHandler(mockCollectionService, mockUserCollectionService, mockUserReactionService, mockUserService, dtoBuilder)

	httpHandler := handlers.MakeHandler(handler.CreateCollection)

	tests := []struct {
		name           string
		expectedStatus int
		setup          func() *http.Request
	}{
		{
			name:           "success",
			expectedStatus: http.StatusCreated,
			setup: func() *http.Request {
				userID := "123e4567-e89b-12d3-a456-426614174000"
				userUUID := uuid.MustParse(userID)
				collectionID := uuid.New()

				mockCollectionService.EXPECT().SaveCollection(gomock.Any(), gomock.Any()).Return(&models.Collection{
					BaseModel:   models.BaseModel{ID: collectionID},
					Title:       "Test Title",
					Description: "Test Description",
					UserID:      userUUID,
				}, nil)
				mockUserCollectionService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				mockUserService.EXPECT().GetUserFollowersCount(gomock.Any(), userUUID).Return(int64(0))
				mockCollectionService.EXPECT().GetByID(gomock.Any(), collectionID, gomock.Any()).Return(&models.Collection{
					BaseModel:   models.BaseModel{ID: collectionID},
					Title:       "Test Title",
					Description: "Test Description",
					UserID:      userUUID,
					User: models.User{
						BaseModel: models.BaseModel{ID: userUUID},
						Username:  "Test User",
					},
				}, nil)

				var body bytes.Buffer
				writer := multipart.NewWriter(&body)
				writer.WriteField("title", "Test Title")
				writer.WriteField("description", "Test Description")
				part, _ := writer.CreateFormFile("cover", "cover.jpg")
				part.Write([]byte("fake image content"))
				writer.Close()

				req := httptest.NewRequest("POST", "/collections", &body)
				req.Header.Set("Content-Type", writer.FormDataContentType())

				ctx := context.WithValue(req.Context(), "userID", userID)
				return req.WithContext(ctx)
			},
		},
		{
			name:           "missing user ID in context",
			expectedStatus: http.StatusBadRequest,
			setup: func() *http.Request {
				var body bytes.Buffer
				writer := multipart.NewWriter(&body)
				writer.WriteField("title", "Some Title")
				writer.WriteField("description", "Some Description")
				part, _ := writer.CreateFormFile("cover", "cover.jpg")
				part.Write([]byte("fake image content"))
				writer.Close()

				req := httptest.NewRequest("POST", "/collections", &body)
				req.Header.Set("Content-Type", writer.FormDataContentType())

				return req
			},
		},
		{
			name:           "missing cover file",
			expectedStatus: http.StatusBadRequest,
			setup: func() *http.Request {
				userID := "123e4567-e89b-12d3-a456-426614174000"

				var body bytes.Buffer
				writer := multipart.NewWriter(&body)
				writer.WriteField("title", "Test Title")
				writer.WriteField("description", "Test Description")
				writer.Close()

				req := httptest.NewRequest("POST", "/collections", &body)
				req.Header.Set("Content-Type", writer.FormDataContentType())

				ctx := context.WithValue(req.Context(), "userID", userID)
				return req.WithContext(ctx)
			},
		},
		{
			name:           "SaveCollection returns error",
			expectedStatus: http.StatusInternalServerError,
			setup: func() *http.Request {
				userID := "123e4567-e89b-12d3-a456-426614174000"

				mockCollectionService.EXPECT().SaveCollection(gomock.Any(), gomock.Any()).Return(nil, errors.New("save error"))

				var body bytes.Buffer
				writer := multipart.NewWriter(&body)
				writer.WriteField("title", "Test Title")
				writer.WriteField("description", "Test Description")
				part, _ := writer.CreateFormFile("cover", "cover.jpg")
				part.Write([]byte("fake image content"))
				writer.Close()

				req := httptest.NewRequest("POST", "/collections", &body)
				req.Header.Set("Content-Type", writer.FormDataContentType())

				ctx := context.WithValue(req.Context(), "userID", userID)
				return req.WithContext(ctx)
			},
		},
		{
			name:           "ParseMultipartForm returns error",
			expectedStatus: http.StatusBadRequest,
			setup: func() *http.Request {
				body := bytes.NewBufferString("invalid multipart body")
				req := httptest.NewRequest("POST", "/collections", body)
				req.Header.Set("Content-Type", "multipart/form-data; boundary=invalid")
				ctx := context.WithValue(req.Context(), "userID", "123e4567-e89b-12d3-a456-426614174000")
				return req.WithContext(ctx)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.setup()
			rr := httptest.NewRecorder()
			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestCollectionHandler_GetCollectionByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionService := mocks.NewMockCollectionService(ctrl)
	mockUserCollectionService := mocks.NewMockUserCollectionService(ctrl)
	mockUserReactionService := mocks.NewMockUserReactionService(ctrl)
	mockUserService := mocks.NewMockUserService(ctrl)

	dtoBuilder := dto.NewDTOBuilder(mockUserService, nil)
	handler := NewCollectionHandler(mockCollectionService, mockUserCollectionService, mockUserReactionService, mockUserService, dtoBuilder)

	httpHandler := handlers.MakeHandler(handler.GetCollectionByID)

	tests := []struct {
		name           string
		expectedStatus int
		setup          func() *http.Request
	}{
		{
			name:           "success",
			expectedStatus: http.StatusOK,
			setup: func() *http.Request {
				collectionID := uuid.New()

				mockCollectionService.EXPECT().GetByID(gomock.Any(), collectionID, gomock.Any()).Return(&models.Collection{
					BaseModel:   models.BaseModel{ID: collectionID},
					Title:       "Test Title",
					Description: "Test Description",
					User: models.User{
						BaseModel: models.BaseModel{ID: uuid.New()},
						Username:  "Test User",
					},
				}, nil)
				mockUserService.EXPECT().GetUserFollowersCount(gomock.Any(), gomock.Any()).Return(int64(0))

				req := httptest.NewRequest("GET", "/collections/"+collectionID.String(), nil)
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", collectionID.String())

				ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
				return req.WithContext(ctx)
			},
		},
		{
			name:           "invalid UUID format",
			expectedStatus: http.StatusBadRequest,
			setup: func() *http.Request {
				req := httptest.NewRequest("GET", "/collections/invalid-uuid", nil)
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", "invalid-uuid")

				ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
				return req.WithContext(ctx)
			},
		},
		{
			name:           "collection not found / GetByID returns error",
			expectedStatus: http.StatusInternalServerError,
			setup: func() *http.Request {
				collectionID := uuid.New()

				mockCollectionService.EXPECT().GetByID(gomock.Any(), collectionID, gomock.Any()).Return(nil, errors.New("not found"))

				req := httptest.NewRequest("GET", "/collections/"+collectionID.String(), nil)
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", collectionID.String())

				ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
				return req.WithContext(ctx)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.setup()
			rr := httptest.NewRecorder()
			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestCollectionHandler_DeleteCollection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionService := mocks.NewMockCollectionService(ctrl)
	mockUserCollectionService := mocks.NewMockUserCollectionService(ctrl)
	mockUserReactionService := mocks.NewMockUserReactionService(ctrl)
	mockUserService := mocks.NewMockUserService(ctrl)

	dtoBuilder := dto.NewDTOBuilder(mockUserService, nil)
	handler := NewCollectionHandler(mockCollectionService, mockUserCollectionService, mockUserReactionService, mockUserService, dtoBuilder)

	httpHandler := handlers.MakeHandler(handler.DeleteCollection)

	tests := []struct {
		name           string
		expectedStatus int
		setup          func() *http.Request
	}{
		{
			name:           "success",
			expectedStatus: http.StatusNoContent,
			setup: func() *http.Request {
				collectionID := uuid.New()

				mockCollectionService.EXPECT().Delete(gomock.Any(), collectionID).Return(nil)

				req := httptest.NewRequest("DELETE", "/collections/"+collectionID.String(), nil)
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", collectionID.String())

				ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
				return req.WithContext(ctx)
			},
		},
		{
			name:           "invalid UUID format",
			expectedStatus: http.StatusBadRequest,
			setup: func() *http.Request {
				req := httptest.NewRequest("DELETE", "/collections/invalid-uuid", nil)
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", "invalid-uuid")

				ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
				return req.WithContext(ctx)
			},
		},
		{
			name:           "Delete returns error",
			expectedStatus: http.StatusInternalServerError,
			setup: func() *http.Request {
				collectionID := uuid.New()

				mockCollectionService.EXPECT().Delete(gomock.Any(), collectionID).Return(errors.New("some error"))

				req := httptest.NewRequest("DELETE", "/collections/"+collectionID.String(), nil)
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", collectionID.String())

				ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
				return req.WithContext(ctx)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.setup()
			rr := httptest.NewRecorder()
			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestCollectionHandler_UpdateCollection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionService := mocks.NewMockCollectionService(ctrl)
	mockUserCollectionService := mocks.NewMockUserCollectionService(ctrl)
	mockUserReactionService := mocks.NewMockUserReactionService(ctrl)
	mockUserService := mocks.NewMockUserService(ctrl)

	dtoBuilder := dto.NewDTOBuilder(mockUserService, nil)
	handler := NewCollectionHandler(mockCollectionService, mockUserCollectionService, mockUserReactionService, mockUserService, dtoBuilder)

	httpHandler := handlers.MakeHandler(handler.UpdateCollection)

	tests := []struct {
		name           string
		expectedStatus int
		setup          func() *http.Request
	}{
		{
			name:           "success",
			expectedStatus: http.StatusOK,
			setup: func() *http.Request {
				collectionID := uuid.New()
				userID := uuid.New()

				mockCollectionService.EXPECT().GetByID(gomock.Any(), collectionID, gomock.Any()).Return(&models.Collection{
					BaseModel: models.BaseModel{ID: collectionID},
					Title:     "Old Title",
					User: models.User{
						BaseModel: models.BaseModel{ID: userID},
					},
				}, nil)

				mockCollectionService.EXPECT().UpdateCollection(gomock.Any(), collectionID, gomock.Any()).Return(&models.Collection{
					BaseModel:   models.BaseModel{ID: collectionID},
					Title:       "New Title",
					Description: "New Description",
				}, nil)

				mockCollectionService.EXPECT().GetByID(gomock.Any(), collectionID, gomock.Any()).Return(&models.Collection{
					BaseModel:   models.BaseModel{ID: collectionID},
					Title:       "New Title",
					Description: "New Description",
					User: models.User{
						BaseModel: models.BaseModel{ID: userID},
						Username:  "Test User",
					},
				}, nil)

				mockUserService.EXPECT().GetUserFollowersCount(gomock.Any(), gomock.Any()).Return(int64(0))

				var body bytes.Buffer
				writer := multipart.NewWriter(&body)
				writer.WriteField("title", "New Title")
				writer.WriteField("description", "New Description")
				part, _ := writer.CreateFormFile("cover", "cover.jpg")
				part.Write([]byte("fake image content"))
				writer.Close()

				req := httptest.NewRequest("PUT", "/collections/"+collectionID.String(), &body)
				req.Header.Set("Content-Type", writer.FormDataContentType())

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", collectionID.String())
				ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
				return req.WithContext(ctx)
			},
		},
		{
			name:           "invalid UUID",
			expectedStatus: http.StatusBadRequest,
			setup: func() *http.Request {
				req := httptest.NewRequest("PUT", "/collections/invalid-id", nil)
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", "invalid-id")
				ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
				return req.WithContext(ctx)
			},
		},
		{
			name:           "error parsing form",
			expectedStatus: http.StatusBadRequest,
			setup: func() *http.Request {
				collectionID := uuid.New()
				req := httptest.NewRequest("PUT", "/collections/"+collectionID.String(), strings.NewReader("invalid form data"))
				req.Header.Set("Content-Type", "multipart/form-data; boundary=--invalid")

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", collectionID.String())
				ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
				return req.WithContext(ctx)
			},
		},
		{
			name:           "GetByID returns error",
			expectedStatus: http.StatusInternalServerError,
			setup: func() *http.Request {
				collectionID := uuid.New()

				mockCollectionService.EXPECT().GetByID(gomock.Any(), collectionID, gomock.Any()).Return(nil, errors.New("not found"))

				var body bytes.Buffer
				writer := multipart.NewWriter(&body)
				writer.WriteField("title", "Some Title")
				writer.WriteField("description", "Some Description")
				writer.Close()

				req := httptest.NewRequest("PUT", "/collections/"+collectionID.String(), &body)
				req.Header.Set("Content-Type", writer.FormDataContentType())

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", collectionID.String())
				ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
				return req.WithContext(ctx)
			},
		},
		{
			name:           "UpdateCollection returns error",
			expectedStatus: http.StatusInternalServerError,
			setup: func() *http.Request {
				collectionID := uuid.New()
				userID := uuid.New()

				mockCollectionService.EXPECT().GetByID(gomock.Any(), collectionID, gomock.Any()).Return(&models.Collection{
					BaseModel: models.BaseModel{ID: collectionID},
					User: models.User{
						BaseModel: models.BaseModel{ID: userID},
					},
				}, nil)

				mockCollectionService.EXPECT().UpdateCollection(gomock.Any(), collectionID, gomock.Any()).Return(nil, errors.New("update failed"))

				var body bytes.Buffer
				writer := multipart.NewWriter(&body)
				writer.WriteField("title", "New Title")
				writer.WriteField("description", "New Description")
				writer.Close()

				req := httptest.NewRequest("PUT", "/collections/"+collectionID.String(), &body)
				req.Header.Set("Content-Type", writer.FormDataContentType())

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", collectionID.String())
				ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
				return req.WithContext(ctx)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.setup()
			rr := httptest.NewRecorder()
			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
