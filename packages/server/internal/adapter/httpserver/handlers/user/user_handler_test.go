package user

import (
	"encoding/json"
	"errors"
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

func TestUserHandler_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	mockUserFollowerService := mocks.NewMockUserFollowerService(ctrl)
	mockUserReactionService := mocks.NewMockUserReactionService(ctrl)
	mockMessageService := mocks.NewMockMessageService(ctrl)

	dtoBuilder := dto.NewDTOBuilder(mockUserService, mockUserReactionService)

	handler := NewUserHandler(
		mockUserService,
		mockUserFollowerService,
		mockUserReactionService,
		mockMessageService,
		dtoBuilder,
	)

	httpHandler := handlers.MakeHandler(handler.GetAll)

	users := []models.User{
		{
			BaseModel: models.BaseModel{ID: uuid.New()},
			Username:  "testuser1",
		},
		{
			BaseModel: models.BaseModel{ID: uuid.New()},
			Username:  "testuser2",
		},
	}

	tests := []struct {
		name           string
		queryParams    string
		setupMocks     func()
		expectedStatus int
	}{
		{
			name:        "success with default pagination",
			queryParams: "",
			setupMocks: func() {
				mockUserService.EXPECT().
					Where(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(users, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "success with custom pagination",
			queryParams: "?page=2&limit=5",
			setupMocks: func() {
				mockUserService.EXPECT().
					Where(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(users, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "internal error from service",
			queryParams: "?page=abc&limit=xyz",
			setupMocks: func() {
				mockUserService.EXPECT().
					Where(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("db failure"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodGet, "/users"+tt.queryParams, nil)
			rr := httptest.NewRecorder()

			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if rr.Code == http.StatusOK {
				var actual []models.User
				err := json.Unmarshal(rr.Body.Bytes(), &actual)
				assert.NoError(t, err)
			} else {
				var actual map[string]string
				err := json.Unmarshal(rr.Body.Bytes(), &actual)
				assert.NoError(t, err)
			}
		})
	}
}
