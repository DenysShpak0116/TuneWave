package collection

import (
	"bytes"
	"context"
	"encoding/json"
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

func TestCollectionHandler_GetUsersCollections(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionService := mocks.NewMockCollectionService(ctrl)
	mockUserCollectionService := mocks.NewMockUserCollectionService(ctrl)
	mockUserReactionService := mocks.NewMockUserReactionService(ctrl)
	mockUserService := mocks.NewMockUserService(ctrl)

	dtoBuilder := dto.NewDTOBuilder(mockUserService, nil)
	handler := NewCollectionHandler(
		mockCollectionService,
		mockUserCollectionService,
		mockUserReactionService,
		mockUserService,
		dtoBuilder,
	)

	httpHandler := handlers.MakeHandler(handler.GetUsersCollections)

	fixedUserID := uuid.New()
	fixedCollection := models.Collection{
		BaseModel: models.BaseModel{ID: uuid.New()},
		Title:     "Test Collection",
	}
	fixedUserCollection := models.UserCollection{
		UserID:     fixedUserID,
		Collection: fixedCollection,
	}

	tests := []struct {
		name           string
		expectedStatus int
		setup          func() *http.Request
	}{
		{
			name:           "success",
			expectedStatus: http.StatusOK,
			setup: func() *http.Request {
				mockUserCollectionService.EXPECT().
					Where(gomock.Any(), &models.UserCollection{UserID: fixedUserID}, gomock.Any()).
					Return([]models.UserCollection{fixedUserCollection}, nil)

				req := httptest.NewRequest("GET", "/collections", nil)
				ctx := context.WithValue(req.Context(), "userID", fixedUserID.String())
				return req.WithContext(ctx)
			},
		},
		{
			name:           "invalid user ID",
			expectedStatus: http.StatusBadRequest,
			setup: func() *http.Request {
				req := httptest.NewRequest("GET", "/collections", nil)
				return req
			},
		},
		{
			name:           "internal service error",
			expectedStatus: http.StatusInternalServerError,
			setup: func() *http.Request {
				mockUserCollectionService.EXPECT().
					Where(gomock.Any(), &models.UserCollection{UserID: fixedUserID}, gomock.Any()).
					Return(nil, errors.New("db error"))

				req := httptest.NewRequest("GET", "/collections", nil)
				ctx := context.WithValue(req.Context(), "userID", fixedUserID.String())
				return req.WithContext(ctx)
			},
		},
		{
			name:           "empty result",
			expectedStatus: http.StatusOK,
			setup: func() *http.Request {
				mockUserCollectionService.EXPECT().
					Where(gomock.Any(), &models.UserCollection{UserID: fixedUserID}, gomock.Any()).
					Return([]models.UserCollection{}, nil)

				req := httptest.NewRequest("GET", "/collections", nil)
				ctx := context.WithValue(req.Context(), "userID", fixedUserID.String())
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

func TestCollectionHandler_GetCollections(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionService := mocks.NewMockCollectionService(ctrl)
	mockUserCollectionService := mocks.NewMockUserCollectionService(ctrl)
	mockUserReactionService := mocks.NewMockUserReactionService(ctrl)
	mockUserService := mocks.NewMockUserService(ctrl)

	dtoBuilder := dto.NewDTOBuilder(mockUserService, nil)
	handler := NewCollectionHandler(
		mockCollectionService,
		mockUserCollectionService,
		mockUserReactionService,
		mockUserService,
		dtoBuilder,
	)

	httpHandler := handlers.MakeHandler(handler.GetCollections)

	tests := []struct {
		name           string
		queryParams    string
		mockSetup      func()
		expectedStatus int
	}{
		{
			name:        "success with defaults",
			queryParams: "",
			mockSetup: func() {
				mockCollectionService.EXPECT().Where(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return([]models.Collection{
					{
						BaseModel: models.BaseModel{ID: uuid.New()},
						Title:     "Test Collection",
						User: models.User{
							BaseModel: models.BaseModel{ID: uuid.New()},
							Username:  "TestUser",
						},
					},
				}, nil)

				mockUserService.EXPECT().GetUserFollowersCount(gomock.Any(), gomock.Any()).Return(int64(0))
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "invalid limit and page should fallback to defaults",
			queryParams: "?limit=invalid&page=invalid",
			mockSetup: func() {
				mockCollectionService.EXPECT().Where(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return([]models.Collection{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "collection service returns error",
			queryParams: "?limit=10&page=1",
			mockSetup: func() {
				mockCollectionService.EXPECT().Where(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(nil, errors.New("some db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := httptest.NewRequest("GET", "/collections"+tt.queryParams, nil)
			rr := httptest.NewRecorder()
			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestCollectionHandler_AddCollectionToUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionService := mocks.NewMockCollectionService(ctrl)
	mockUserCollectionService := mocks.NewMockUserCollectionService(ctrl)
	mockUserReactionService := mocks.NewMockUserReactionService(ctrl)
	mockUserService := mocks.NewMockUserService(ctrl)

	dtoBuilder := dto.NewDTOBuilder(mockUserService, nil)
	handler := NewCollectionHandler(
		mockCollectionService,
		mockUserCollectionService,
		mockUserReactionService,
		mockUserService,
		dtoBuilder,
	)

	httpHandler := handlers.MakeHandler(handler.AddCollectionToUser)

	fixedUserID := uuid.New()
	fixedCollectionID := uuid.New()
	fixedCollection := models.Collection{
		BaseModel: models.BaseModel{ID: fixedCollectionID},
		Title:     "Test Collection",
	}
	userCollection := models.UserCollection{
		UserID:     fixedUserID,
		Collection: fixedCollection,
	}

	tests := []struct {
		name           string
		expectedStatus int
		setup          func() *http.Request
		mockSetup      func()
	}{
		{
			name:           "success",
			expectedStatus: http.StatusCreated,
			setup: func() *http.Request {
				req := httptest.NewRequest("POST", "/collections/"+fixedCollectionID.String(), nil)
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", fixedCollectionID.String())
				ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
				ctx = context.WithValue(ctx, "userID", fixedUserID.String())
				return req.WithContext(ctx)
			},
			mockSetup: func() {
				mockUserCollectionService.EXPECT().
					First(gomock.Any(), &models.UserCollection{
						UserID:       fixedUserID,
						CollectionID: fixedCollectionID,
					}).
					Return(nil, nil)

				mockUserCollectionService.EXPECT().
					Create(gomock.Any(), &models.UserCollection{
						UserID:       fixedUserID,
						CollectionID: fixedCollectionID,
					}).
					Return(nil)

				mockUserCollectionService.EXPECT().
					Where(gomock.Any(), &models.UserCollection{UserID: fixedUserID}, gomock.Any()).
					Return([]models.UserCollection{userCollection}, nil)
			},
		},
		{
			name:           "invalid collection UUID",
			expectedStatus: http.StatusBadRequest,
			setup: func() *http.Request {
				req := httptest.NewRequest("POST", "/collections/invalid-uuid", nil)
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", "invalid-uuid")
				ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
				ctx = context.WithValue(ctx, "userID", fixedUserID.String())
				return req.WithContext(ctx)
			},
			mockSetup: func() {},
		},
		{
			name:           "missing user ID in context",
			expectedStatus: http.StatusBadRequest,
			setup: func() *http.Request {
				req := httptest.NewRequest("POST", "/collections/"+fixedCollectionID.String(), nil)
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", fixedCollectionID.String())
				ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
				return req.WithContext(ctx)
			},
			mockSetup: func() {},
		},
		{
			name:           "error on First",
			expectedStatus: http.StatusBadRequest,
			setup: func() *http.Request {
				req := httptest.NewRequest("POST", "/collections/"+fixedCollectionID.String(), nil)
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", fixedCollectionID.String())
				ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
				ctx = context.WithValue(ctx, "userID", fixedUserID.String())
				return req.WithContext(ctx)
			},
			mockSetup: func() {
				mockUserCollectionService.EXPECT().
					First(gomock.Any(), &models.UserCollection{
						UserID:       fixedUserID,
						CollectionID: fixedCollectionID,
					}).
					Return(nil, errors.New("not found"))
			},
		},
		{
			name:           "error on Create",
			expectedStatus: http.StatusBadRequest,
			setup: func() *http.Request {
				req := httptest.NewRequest("POST", "/collections/"+fixedCollectionID.String(), nil)
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", fixedCollectionID.String())
				ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
				ctx = context.WithValue(ctx, "userID", fixedUserID.String())
				return req.WithContext(ctx)
			},
			mockSetup: func() {
				mockUserCollectionService.EXPECT().
					First(gomock.Any(), &models.UserCollection{
						UserID:       fixedUserID,
						CollectionID: fixedCollectionID,
					}).
					Return(nil, nil)

				mockUserCollectionService.EXPECT().
					Create(gomock.Any(), &models.UserCollection{
						UserID:       fixedUserID,
						CollectionID: fixedCollectionID,
					}).
					Return(errors.New("create error"))
			},
		},
		{
			name:           "error on Where",
			expectedStatus: http.StatusInternalServerError,
			setup: func() *http.Request {
				req := httptest.NewRequest("POST", "/collections/"+fixedCollectionID.String(), nil)
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", fixedCollectionID.String())
				ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
				ctx = context.WithValue(ctx, "userID", fixedUserID.String())
				return req.WithContext(ctx)
			},
			mockSetup: func() {
				mockUserCollectionService.EXPECT().
					First(gomock.Any(), &models.UserCollection{
						UserID:       fixedUserID,
						CollectionID: fixedCollectionID,
					}).
					Return(nil, nil)

				mockUserCollectionService.EXPECT().
					Create(gomock.Any(), &models.UserCollection{
						UserID:       fixedUserID,
						CollectionID: fixedCollectionID,
					}).
					Return(nil)

				mockUserCollectionService.EXPECT().
					Where(gomock.Any(), &models.UserCollection{UserID: fixedUserID}, gomock.Any()).
					Return(nil, errors.New("db error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			req := tt.setup()
			rr := httptest.NewRecorder()
			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestCollectionHandler_RemoveCollectionFromUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionService := mocks.NewMockCollectionService(ctrl)
	mockUserCollectionService := mocks.NewMockUserCollectionService(ctrl)
	mockUserReactionService := mocks.NewMockUserReactionService(ctrl)
	mockUserService := mocks.NewMockUserService(ctrl)

	dtoBuilder := dto.NewDTOBuilder(mockUserService, nil)
	handler := NewCollectionHandler(
		mockCollectionService,
		mockUserCollectionService,
		mockUserReactionService,
		mockUserService,
		dtoBuilder,
	)

	httpHandler := handlers.MakeHandler(handler.RemoveCollectionFromUser)

	fixedUserID := uuid.New()
	fixedCollectionID := uuid.New()
	fixedUserCollection := &models.UserCollection{
		BaseModel:    models.BaseModel{ID: uuid.New()},
		UserID:       fixedUserID,
		CollectionID: fixedCollectionID,
	}

	tests := []struct {
		name           string
		expectedStatus int
		setup          func() *http.Request
		mockSetup      func()
	}{
		{
			name:           "success",
			expectedStatus: http.StatusNoContent,
			setup: func() *http.Request {
				req := httptest.NewRequest("DELETE", "/collections/"+fixedCollectionID.String(), nil)
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", fixedCollectionID.String())
				ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
				ctx = context.WithValue(ctx, "userID", fixedUserID.String())
				return req.WithContext(ctx)
			},
			mockSetup: func() {
				mockUserCollectionService.EXPECT().
					First(gomock.Any(), &models.UserCollection{
						UserID:       fixedUserID,
						CollectionID: fixedCollectionID,
					}).
					Return(fixedUserCollection, nil)

				mockUserCollectionService.EXPECT().
					Delete(gomock.Any(), fixedUserCollection.ID).
					Return(nil)
			},
		},
		{
			name:           "invalid collection ID",
			expectedStatus: http.StatusBadRequest,
			setup: func() *http.Request {
				req := httptest.NewRequest("DELETE", "/collections/invalid-id", nil)
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", "invalid-id")
				ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
				ctx = context.WithValue(ctx, "userID", fixedUserID.String())
				return req.WithContext(ctx)
			},
			mockSetup: func() {},
		},
		{
			name:           "invalid user ID",
			expectedStatus: http.StatusBadRequest,
			setup: func() *http.Request {
				req := httptest.NewRequest("DELETE", "/collections/"+fixedCollectionID.String(), nil)
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", fixedCollectionID.String())
				ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
				return req.WithContext(ctx)
			},
			mockSetup: func() {},
		},
		{
			name:           "First returns error",
			expectedStatus: http.StatusInternalServerError,
			setup: func() *http.Request {
				req := httptest.NewRequest("DELETE", "/collections/"+fixedCollectionID.String(), nil)
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", fixedCollectionID.String())
				ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
				ctx = context.WithValue(ctx, "userID", fixedUserID.String())
				return req.WithContext(ctx)
			},
			mockSetup: func() {
				mockUserCollectionService.EXPECT().
					First(gomock.Any(), &models.UserCollection{
						UserID:       fixedUserID,
						CollectionID: fixedCollectionID,
					}).
					Return(nil, errors.New("not found"))
			},
		},
		{
			name:           "Delete returns error",
			expectedStatus: http.StatusInternalServerError,
			setup: func() *http.Request {
				req := httptest.NewRequest("DELETE", "/collections/"+fixedCollectionID.String(), nil)
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", fixedCollectionID.String())
				ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
				ctx = context.WithValue(ctx, "userID", fixedUserID.String())
				return req.WithContext(ctx)
			},
			mockSetup: func() {
				mockUserCollectionService.EXPECT().
					First(gomock.Any(), &models.UserCollection{
						UserID:       fixedUserID,
						CollectionID: fixedCollectionID,
					}).
					Return(fixedUserCollection, nil)

				mockUserCollectionService.EXPECT().
					Delete(gomock.Any(), fixedUserCollection.ID).
					Return(errors.New("delete failed"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := tt.setup()
			rr := httptest.NewRecorder()
			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestCollectionHandler_GetCollectionSongs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionService := mocks.NewMockCollectionService(ctrl)
	mockUserCollectionService := mocks.NewMockUserCollectionService(ctrl)
	mockUserReactionService := mocks.NewMockUserReactionService(ctrl)
	mockUserService := mocks.NewMockUserService(ctrl)
	mockSongReactionService := mocks.NewMockUserReactionService(ctrl)

	dtoBuilder := dto.NewDTOBuilder(mockUserService, mockSongReactionService)
	handler := NewCollectionHandler(
		mockCollectionService,
		mockUserCollectionService,
		mockUserReactionService,
		mockUserService,
		dtoBuilder,
	)

	httpHandler := handlers.MakeHandler(handler.GetCollectionSongs)

	collectionID := uuid.New()
	song := models.Song{BaseModel: models.BaseModel{ID: uuid.New()}, Title: "Test Song"}
	songList := []models.Song{song}

	tests := []struct {
		name           string
		url            string
		expectedStatus int
		mockSetup      func()
	}{
		{
			name:           "success",
			url:            "/collections/" + collectionID.String() + "/songs?page=1&limit=5&sortBy=title&order=asc",
			expectedStatus: http.StatusOK,
			mockSetup: func() {
				mockCollectionService.EXPECT().
					GetCollectionSongs(gomock.Any(), collectionID, "", "title", "asc", 1, 5).
					Return(songList, nil)

				mockSongReactionService.EXPECT().GetSongLikes(gomock.Any(), gomock.Any()).Return(int64(0))
				mockSongReactionService.EXPECT().GetSongDislikes(gomock.Any(), gomock.Any()).Return(int64(0))
			},
		},
		{
			name:           "invalid collection ID",
			url:            "/collections/invalid-uuid/songs",
			expectedStatus: http.StatusBadRequest,
			mockSetup:      func() {},
		},
		{
			name:           "invalid page & limit fallback",
			url:            "/collections/" + collectionID.String() + "/songs?page=abc&limit=-2",
			expectedStatus: http.StatusOK,
			mockSetup: func() {
				mockCollectionService.EXPECT().
					GetCollectionSongs(gomock.Any(), collectionID, "", "created_at", "desc", 1, 10).
					Return(songList, nil)

				mockSongReactionService.EXPECT().GetSongLikes(gomock.Any(), gomock.Any()).Return(int64(0))
				mockSongReactionService.EXPECT().GetSongDislikes(gomock.Any(), gomock.Any()).Return(int64(0))
			},
		},
		{
			name:           "internal service error",
			url:            "/collections/" + collectionID.String() + "/songs",
			expectedStatus: http.StatusInternalServerError,
			mockSetup: func() {
				mockCollectionService.EXPECT().
					GetCollectionSongs(gomock.Any(), collectionID, "", "created_at", "desc", 1, 10).
					Return(nil, errors.New("db error"))

			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := httptest.NewRequest("GET", tt.url, nil)
			rctx := chi.NewRouteContext()

			idPart := strings.Split(tt.url, "/")[2]
			rctx.URLParams.Add("id", idPart)

			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedStatus == http.StatusOK {
				var songs []models.Song
				err := json.NewDecoder(rr.Body).Decode(&songs)
				assert.NoError(t, err)
				assert.Len(t, songs, len(songList))
				assert.Equal(t, "Test Song", songs[0].Title)
			}
		})
	}
}
