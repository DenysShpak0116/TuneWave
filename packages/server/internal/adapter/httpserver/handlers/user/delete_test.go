package user

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/service/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserHandler_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	mockUserFollowerService := mocks.NewMockUserFollowerService(ctrl)
	mockUserReactionService := mocks.NewMockUserReactionService(ctrl)
	mockMessageService := mocks.NewMockMessageService(ctrl)

	dtoBuilder := dto.NewDTOBuilder(mockUserService, nil)

	handler := NewUserHandler(
		mockUserService,
		mockUserFollowerService,
		mockUserReactionService,
		mockMessageService,
		dtoBuilder,
	)

	httpHandler := handlers.MakeHandler(handler.Delete)

	validUserID := uuid.New()

	tests := []struct {
		name             string
		userID           string
		setupMocks       func()
		expectedCode     int
		expectedResponse string
	}{
		{
			name:   "success",
			userID: validUserID.String(),
			setupMocks: func() {
				mockUserService.
					EXPECT().
					Delete(gomock.Any(), validUserID).
					Return(nil)
			},
			expectedCode:     http.StatusNoContent,
			expectedResponse: "",
		},
		{
			name:   "invalid user ID",
			userID: "invalid-uuid",
			setupMocks: func() {
			},
			expectedCode:     http.StatusBadRequest,
			expectedResponse: `{"error":"invalid user id"}`,
		},
		{
			name:   "delete failed",
			userID: validUserID.String(),
			setupMocks: func() {
				mockUserService.
					EXPECT().
					Delete(gomock.Any(), validUserID).
					Return(errors.New("delete failed"))
			},
			expectedCode:     http.StatusInternalServerError,
			expectedResponse: `{"error":"failed to delete user"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodDelete, "/users/"+tt.userID, nil)
			req = helpers.MuxSetURLParam(req, "id", tt.userID)

			rr := httptest.NewRecorder()
			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)

			if tt.expectedResponse != "" {
				assert.JSONEq(t, tt.expectedResponse, rr.Body.String())
			} else {
				assert.Empty(t, rr.Body.String())
			}
		})
	}
}
