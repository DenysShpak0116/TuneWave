package auth

import (
	"fmt"
	"time"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/config"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	authService        services.AuthService
	userService        services.UserService
	dtoBuilder         *dto.DTOBuilder
	googleClientID     string
	googleClientSecret string
	jwtSecret          string
}

func NewAuthHandler(
	authService services.AuthService,
	userService services.UserService,
	dtoBuilder *dto.DTOBuilder,
	cfg *config.Config,
) *AuthHandler {
	goth.UseProviders(
		google.New(
			cfg.Google.ClientID,
			cfg.Google.ClientSecret,
			"http://localhost:8081/auth/google/callback",
			"email",
			"profile",
		),
	)

	return &AuthHandler{
		authService:        authService,
		userService:        userService,
		dtoBuilder:         dtoBuilder,
		googleClientID:     cfg.Google.ClientID,
		googleClientSecret: cfg.Google.ClientSecret,
		jwtSecret:          cfg.JwtSecret,
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	fmt.Println(err == nil)
	return err == nil
}

func (ah *AuthHandler) GenerateTokens(userID string) (string, string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID,
		"exp":    time.Now().Add(5 * time.Hour).Unix(),
	})

	accessTokenStr, err := accessToken.SignedString([]byte(ah.jwtSecret))
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID,
		"exp":    time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	refreshTokenStr, err := refreshToken.SignedString([]byte(ah.jwtSecret))
	if err != nil {
		return "", "", err
	}

	return accessTokenStr, refreshTokenStr, nil
}
