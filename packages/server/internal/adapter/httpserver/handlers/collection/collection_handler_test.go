package collection

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/service/mocks"
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
		setupMocks     func() *http.Request
	}{
		{
			name:           "success",
			expectedStatus: http.StatusCreated,
			setupMocks: func() *http.Request {
				userID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
				collectionID := uuid.New()

				mockCollectionService.EXPECT().SaveCollection(gomock.Any(), gomock.Any()).Return(&models.Collection{
					BaseModel:   models.BaseModel{ID: collectionID},
					Title:       "Test Title",
					Description: "Test Description",
					UserID:      userID,
				}, nil)

				mockUserCollectionService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

				mockCollectionService.EXPECT().GetByID(gomock.Any(), collectionID, gomock.Any()).Return(&models.Collection{
					BaseModel:   models.BaseModel{ID: collectionID},
					Title:       "Test Title",
					Description: "Test Description",
					UserID:      userID,
					User: models.User{
						BaseModel: models.BaseModel{ID: userID},
						Username:  "Test User",
					},
				}, nil)

				mockUserService.EXPECT().GetByID(gomock.Any(), userID).Return(&models.User{
					BaseModel: models.BaseModel{ID: userID},
					Username:  "Test User",
				}, nil)

				var body bytes.Buffer
				writer := multipart.NewWriter(&body)

				_ = writer.WriteField("title", "Test Title")
				_ = writer.WriteField("description", "Test Description")

				part, _ := writer.CreateFormFile("cover", "cover.jpg")
				part.Write([]byte("fake image content"))

				writer.Close()

				req := httptest.NewRequest("POST", "/collections", &body)
				req.Header.Set("Content-Type", writer.FormDataContentType())

				ctx := context.WithValue(req.Context(), "userID", userID)
				return req.WithContext(ctx)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.setupMocks()
			rr := httptest.NewRecorder()
			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
