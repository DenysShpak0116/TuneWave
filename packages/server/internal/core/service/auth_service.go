package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	MailService     services.MailService
	TokenRepository port.Repository[models.Token]
	UserService     services.UserService
}

func NewAuthService(
	mailService services.MailService,
	tokenRepository port.Repository[models.Token],
	userService services.UserService,
) services.AuthService {
	return &AuthService{
		MailService:     mailService,
		TokenRepository: tokenRepository,
		UserService:     userService,
	}
}

func (as *AuthService) HandleForgotPassword(email string) (string, error) {
	token := uuid.New().String()
	expiresAt := time.Now().Add(1 * time.Hour)

	newToken := &models.Token{
		Token:     token,
		Email:     email,
		ExpiresAt: expiresAt,
	}
	if err := as.TokenRepository.Add(context.Background(), newToken); err != nil {
		return "", err
	}

	as.MailService.SendEmail(email, "Password Reset", fmt.Sprintf("Token for password: %s", token))
	return token, nil
}

func (as *AuthService) HandleResetPassword(token, newPassword string) error {
	foundToken, err := as.TokenRepository.NewQuery(context.Background()).
		First("token = ?", token)
	if err != nil {
		return errors.New("invalid token")
	}

	if time.Now().After(foundToken.ExpiresAt) {
		return errors.New("token expired")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err = as.UserService.UpdateUserPassword(foundToken.Email, string(hash)); err != nil {
		return err
	}

	_ = as.TokenRepository.Delete(context.TODO(), foundToken.ID)
	return nil
}
