package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log/slog"
	"net/http"
	"time"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/config"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/go-chi/render"
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
	logger             *slog.Logger

	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func NewAuthHandler(
	authService services.AuthService,
	userService services.UserService,
	dtoBuilder *dto.DTOBuilder,
	cfg *config.Config,
	logger *slog.Logger,
) *AuthHandler {
	privateKey, publicKey, err := helpers.GenerateKeys()
	if err != nil {
		logger.Error("Failed to generate RSA keys", "err", err.Error())
		panic("cannot start AuthHandler without RSA keys")
	}

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
		logger:             logger,
		privateKey:         privateKey,
		publicKey:          publicKey,
	}
}

func (ah *AuthHandler) GetPublicKey(w http.ResponseWriter, r *http.Request) error {
	pubASN1, err := x509.MarshalPKIXPublicKey(ah.publicKey)
	if err != nil {
		ah.logger.Error("Failed to marshal public key", "err", err.Error())
		http.Error(w, "failed to marshal public key", http.StatusInternalServerError)
		return err
	}

	pemBlock := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	}

	stringKey := string(pem.EncodeToMemory(&pemBlock))

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{
		"publicKey": stringKey,
	})
	return nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (ah *AuthHandler) GenerateTokens(userID string) (string, string, error) {
	const op = "adapter.httpserver.handlers.auth.AuthHandler.GenerateTokens"
	logger := ah.logger.With(
		"op", op,
		"userID", userID,
	)

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID,
		"exp":    time.Now().Add(5 * time.Hour).Unix(),
	})

	accessTokenStr, err := accessToken.SignedString([]byte(ah.jwtSecret))
	if err != nil {
		logger.Error("Failed to create access token", "err", err.Error())
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID,
		"exp":    time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	refreshTokenStr, err := refreshToken.SignedString([]byte(ah.jwtSecret))
	if err != nil {
		logger.Error("Failed to create refresh token", "err", err.Error())
		return "", "", err
	}

	logger.Info("Access and refresh tokens successfully generated")
	return accessTokenStr, refreshTokenStr, nil
}
