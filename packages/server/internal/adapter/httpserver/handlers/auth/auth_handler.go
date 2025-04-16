package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	UserService        services.UserService
	GoogleClientID     string
	GoogleClientSecret string
	JWTSecret          string
}

func NewAuthHandler(
	userService services.UserService,
	googleClientID,
	googleClientSecret string,
	jwtSecret string,
) *AuthHandler {
	goth.UseProviders(
		google.New(
			googleClientID,
			googleClientSecret,
			"http://localhost:8081/auth/google/callback",
			"email",
			"profile",
		),
	)

	return &AuthHandler{
		UserService:        userService,
		GoogleClientID:     googleClientID,
		GoogleClientSecret: googleClientSecret,
		JWTSecret:          jwtSecret,
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
		"exp":    time.Now().Add(15 * time.Minute).Unix(),
	})

	accessTokenStr, err := accessToken.SignedString([]byte(ah.JWTSecret))
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID,
		"exp":    time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	refreshTokenStr, err := refreshToken.SignedString([]byte(ah.JWTSecret))
	if err != nil {
		return "", "", err
	}

	return accessTokenStr, refreshTokenStr, nil
}

func (ah *AuthHandler) ParseToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(ah.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["userId"] == nil {
		return "", errors.New("invalid token claims")
	}

	return claims["userId"].(string), nil
}
