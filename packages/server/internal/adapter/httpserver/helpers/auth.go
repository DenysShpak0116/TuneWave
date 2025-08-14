package helpers

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func ParseToken(jwtSecret, tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
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

func GetUserID(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return uuid.UUID{}, errors.New("userId not found in context")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return uuid.UUID{}, errors.New("userId could not be parsed to uuid")
	}

	return userUUID, nil
}
