package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/config"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/service"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/service/mocks"
	"go.uber.org/mock/gomock"
)

func TestRegister_Success(t *testing.T) {
	// Arrange
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockAuthService := mocks.NewMockAuthService(controller)
	mockUserService := mocks.NewMockUserService(controller)
	dtoBuilder := dto.NewDTOBuilder(mockUserService, nil)
	cfg := &config.Config{
		JwtSecret: "random-secret",
		Google: config.GoogleConfig{
			ClientID:     "abc",
			ClientSecret: "abc",
		},
	}

	handler := NewAuthHandler(mockAuthService, mockUserService, dtoBuilder, cfg)

	email := "test@example.com"
	password := "testpassword"
	username := "testuser"

	reqBody := map[string]string{
		"email":    email,
		"password": password,
		"username": username,
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	mockUserService.EXPECT().
		First(gomock.Any(), &models.User{Email: email}).
		Return(nil, service.ErrNotFound)

	mockUserService.EXPECT().
		Create(gomock.Any(), gomock.AssignableToTypeOf(&models.User{})).
		Return(nil)

	mockUserService.EXPECT().
		GetUserFollowersCount(gomock.Any(), gomock.Any()).
		Return(0)

	// Act
	recorder := httptest.NewRecorder()
	err := handler.Register(recorder, req)

	// Assert
	if err != nil {
		t.Errorf("Register() error = %v, want nil", err)
	}

	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("Register() status code = %v, want %v", status, http.StatusCreated)
	}
}
