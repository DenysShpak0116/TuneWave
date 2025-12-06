package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/service"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/service/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUserHandler_GetByID(t *testing.T) {
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

	httpHandler := handlers.MakeHandler(handler.GetByID)

	userID := uuid.New()

	tests := []struct {
		name         string
		inputID      string
		setupMocks   func()
		expectedCode int
	}{
		{
			name:    "success",
			inputID: userID.String(),
			setupMocks: func() {
				mockUserService.EXPECT().
					GetByID(gomock.Any(), userID, "Follows", "Follows.User", "Followers", "Followers.Follower").
					Return(&models.User{
						BaseModel: models.BaseModel{ID: userID},
						Username:  "John",
					}, nil)

				mockUserService.EXPECT().
					GetUserFollowersCount(gomock.Any(), userID).
					Return(int64(0))
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid user ID",
			inputID:      "invalid-uuid",
			setupMocks:   func() {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:    "user not found",
			inputID: userID.String(),
			setupMocks: func() {
				mockUserService.EXPECT().
					GetByID(gomock.Any(), userID, "Follows", "Follows.User", "Followers", "Followers.Follower").
					Return(nil, service.ErrNotFound)
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			url := fmt.Sprintf("/users/%s", tt.inputID)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req = helpers.MuxSetURLParam(req, "id", tt.inputID)

			rr := httptest.NewRecorder()
			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
		})
	}
}

func TestUserHandler_GetChats(t *testing.T) {
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

	httpHandler := handlers.MakeHandler(handler.GetChats)

	userID := uuid.New()
	chatID := uuid.New()
	otherUserID := uuid.New()

	tests := []struct {
		name         string
		userID       uuid.UUID
		setupMocks   func()
		expectedCode int
	}{
		{
			name:   "success",
			userID: userID,
			setupMocks: func() {
				mockUserService.
					EXPECT().
					GetByID(gomock.Any(), userID, gomock.Any()).
					Return(&models.User{
						BaseModel: models.BaseModel{ID: userID},
						ChatUsers: []models.ChatUser{
							{
								BaseModel: models.BaseModel{},
								UserID:    userID,
								User:      models.User{BaseModel: models.BaseModel{ID: userID}, Username: "Me"},
								ChatID:    chatID,
								Chat: models.Chat{
									BaseModel: models.BaseModel{ID: chatID},
									ChatUsers: []models.ChatUser{
										{
											UserID: userID,
											User:   models.User{BaseModel: models.BaseModel{ID: userID}, Username: "Me"},
										},
										{
											UserID: otherUserID,
											User:   models.User{BaseModel: models.BaseModel{ID: otherUserID}, Username: "OtherUser", ProfilePicture: "avatar.png"},
										},
									},
								},
							},
						},
					}, nil)

				mockMessageService.
					EXPECT().
					Last(gomock.Any(), &models.Message{ChatID: chatID}).
					Return(&models.Message{Content: "Hello!"}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:   "user not found",
			userID: userID,
			setupMocks: func() {
				mockUserService.
					EXPECT().
					GetByID(gomock.Any(), userID, gomock.Any()).
					Return(nil, service.ErrNotFound)
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodGet, "/chats", nil)
			ctx := context.WithValue(req.Context(), "userID", tt.userID.String())
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)

			if rr.Code == http.StatusOK {
				var actual []map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &actual)
				assert.NoError(t, err)

				require.NotEmpty(t, actual)
				assert.Equal(t, "OtherUser", actual[0]["username"])
				assert.Equal(t, "Hello!", actual[0]["lastMessage"])
			} else {
				var actual map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &actual)
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserHandler_GetUserCollections(t *testing.T) {
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

	httpHandler := handlers.MakeHandler(handler.GetByID)
	userID := uuid.New()

	tests := []struct {
		name           string
		idParam        string
		setupMocks     func()
		expectedStatus int
	}{
		{
			name:    "success",
			idParam: userID.String(),
			setupMocks: func() {
				mockUserService.EXPECT().GetByID(gomock.Any(), userID, gomock.Any()).Return(&models.User{
					BaseModel: models.BaseModel{ID: userID},
					UserCollections: []models.UserCollection{
						{
							Collection: models.Collection{
								BaseModel: models.BaseModel{ID: uuid.New()},
								Title:     "Chill Vibes",
								User: models.User{
									BaseModel: models.BaseModel{ID: userID},
									Username:  "testuser",
								},
							},
						},
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid UUID",
			idParam:        "invalid-uuid",
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "user not found",
			idParam: userID.String(),
			setupMocks: func() {
				mockUserService.EXPECT().GetByID(gomock.Any(), userID, gomock.Any()).Return(nil, errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodGet, "/users/"+tt.idParam+"/collections", nil)
			ctx := context.WithValue(req.Context(), "userID", userID.String())
			req = req.WithContext(ctx)

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
				URLParams: chi.RouteParams{
					Keys:   []string{"id"},
					Values: []string{tt.idParam},
				},
			}))
			rr := httptest.NewRecorder()

			httpHandler.ServeHTTP(rr, req)
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
