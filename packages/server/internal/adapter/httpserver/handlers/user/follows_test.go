package user

import (
	"context"
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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserHandler_FollowUser(t *testing.T) {
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

	httpHandler := handlers.MakeHandler(handler.FollowUser)

	targetUserID := uuid.New()
	currentUserID := uuid.New()

	tests := []struct {
		name         string
		targetID     string
		setupMocks   func()
		setupContext func(*http.Request) *http.Request
		expectedCode int
	}{
		{
			name:     "success",
			targetID: targetUserID.String(),
			setupMocks: func() {
				mockUserFollowerService.EXPECT().
					First(gomock.Any(), &models.UserFollower{
						UserID:     targetUserID,
						FollowerID: currentUserID,
					}).
					Return(&models.UserFollower{}, nil)

				mockUserFollowerService.EXPECT().
					Create(gomock.Any(), &models.UserFollower{
						UserID:     targetUserID,
						FollowerID: currentUserID,
					}).
					Return(nil)

				mockUserFollowerService.EXPECT().
					First(gomock.Any(), &models.UserFollower{
						UserID:     targetUserID,
						FollowerID: currentUserID,
					}, gomock.Any()).
					Return(&models.UserFollower{
						UserID:     targetUserID,
						FollowerID: currentUserID,
					}, nil)

				mockUserService.EXPECT().
					GetUserFollowersCount(gomock.Any(), gomock.Any()).
					Return(int64(0)).
					AnyTimes()
			},
			setupContext: func(r *http.Request) *http.Request {
				ctx := context.WithValue(r.Context(), "userID", currentUserID.String())
				return r.WithContext(ctx)
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid target user ID",
			targetID:     "invalid-uuid",
			setupMocks:   func() {},
			setupContext: func(r *http.Request) *http.Request { return r },
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "unauthorized user",
			targetID:     targetUserID.String(),
			setupMocks:   func() {},
			setupContext: func(r *http.Request) *http.Request { return r },
			expectedCode: http.StatusBadRequest,
		},
		{
			name:     "already followed",
			targetID: targetUserID.String(),
			setupMocks: func() {
				mockUserFollowerService.
					EXPECT().
					First(gomock.Any(), &models.UserFollower{
						UserID:     targetUserID,
						FollowerID: currentUserID,
					}).
					Return(nil, service.ErrNotFound)
			},
			setupContext: func(r *http.Request) *http.Request {
				ctx := context.WithValue(r.Context(), "userID", currentUserID.String())
				return r.WithContext(ctx)
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:     "internal error on first",
			targetID: targetUserID.String(),
			setupMocks: func() {
				mockUserFollowerService.
					EXPECT().
					First(gomock.Any(), &models.UserFollower{
						UserID:     targetUserID,
						FollowerID: currentUserID,
					}).
					Return(nil, errors.New("db error"))
			},
			setupContext: func(r *http.Request) *http.Request {
				ctx := context.WithValue(r.Context(), "userID", currentUserID.String())
				return r.WithContext(ctx)
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			url := fmt.Sprintf("/users/%s/follow", tt.targetID)
			req := httptest.NewRequest(http.MethodPost, url, nil)
			req = helpers.MuxSetURLParam(req, "id", tt.targetID)

			if tt.setupContext != nil {
				req = tt.setupContext(req)
			}

			rr := httptest.NewRecorder()
			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
		})
	}
}

func TestUserHandler_IsFollowed(t *testing.T) {
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

	httpHandler := handlers.MakeHandler(handler.IsFollowed)

	currentUserID := uuid.New()
	targetUserID := uuid.New()

	tests := []struct {
		name             string
		targetID         string
		setupMocks       func()
		setupContext     func(*http.Request) *http.Request
		expectedCode     int
		expectedResponse string
	}{
		{
			name:     "user is followed",
			targetID: targetUserID.String(),
			setupMocks: func() {
				mockUserFollowerService.EXPECT().
					First(gomock.Any(), &models.UserFollower{
						UserID:     targetUserID,
						FollowerID: currentUserID,
					}).
					Return(&models.UserFollower{}, nil)
			},
			setupContext: func(r *http.Request) *http.Request {
				ctx := context.WithValue(r.Context(), "userID", currentUserID.String())
				return r.WithContext(ctx)
			},
			expectedCode:     http.StatusOK,
			expectedResponse: `{"isFollowed": true}`,
		},
		{
			name:     "user is not followed",
			targetID: targetUserID.String(),
			setupMocks: func() {
				mockUserFollowerService.EXPECT().
					First(gomock.Any(), &models.UserFollower{
						UserID:     targetUserID,
						FollowerID: currentUserID,
					}).
					Return(nil, service.ErrNotFound)
			},
			setupContext: func(r *http.Request) *http.Request {
				ctx := context.WithValue(r.Context(), "userID", currentUserID.String())
				return r.WithContext(ctx)
			},
			expectedCode:     http.StatusOK,
			expectedResponse: `{"isFollowed": false}`,
		},
		{
			name:       "invalid target user ID",
			targetID:   "invalid-uuid",
			setupMocks: func() {},
			setupContext: func(r *http.Request) *http.Request {
				ctx := context.WithValue(r.Context(), "userID", currentUserID.String())
				return r.WithContext(ctx)
			},
			expectedCode:     http.StatusBadRequest,
			expectedResponse: `{"error":"wrong target user id"}`,
		},
		{
			name:             "unauthorized user",
			targetID:         targetUserID.String(),
			setupMocks:       func() {},
			setupContext:     func(r *http.Request) *http.Request { return r },
			expectedCode:     http.StatusBadRequest,
			expectedResponse: `{"error":"wrong user id"}`,
		},
		{
			name:     "internal error on follower check",
			targetID: targetUserID.String(),
			setupMocks: func() {
				mockUserFollowerService.EXPECT().
					First(gomock.Any(), &models.UserFollower{
						UserID:     targetUserID,
						FollowerID: currentUserID,
					}).
					Return(nil, errors.New("db error"))
			},
			setupContext: func(r *http.Request) *http.Request {
				ctx := context.WithValue(r.Context(), "userID", currentUserID.String())
				return r.WithContext(ctx)
			},
			expectedCode:     http.StatusInternalServerError,
			expectedResponse: `{"error":"troubles to check if you followed"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			url := fmt.Sprintf("/users/%s/is-followed", tt.targetID)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req = helpers.MuxSetURLParam(req, "id", tt.targetID)

			if tt.setupContext != nil {
				req = tt.setupContext(req)
			}

			rr := httptest.NewRecorder()
			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)

			if tt.expectedResponse != "" {
				assert.JSONEq(t, tt.expectedResponse, rr.Body.String())
			}
		})
	}
}
