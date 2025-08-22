package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
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
	logger          *slog.Logger
}

func NewAuthService(
	mailService services.MailService,
	tokenRepository port.Repository[models.Token],
	userService services.UserService,
	logger *slog.Logger,
) services.AuthService {
	return &AuthService{
		MailService:     mailService,
		TokenRepository: tokenRepository,
		UserService:     userService,
		logger:          logger,
	}
}

func (as *AuthService) HandleForgotPassword(email string) (string, error) {
	const op = "core.service.AuthSerivce.HandleForgotPassword"
	logger := as.logger.With(
		slog.String("op", op),
	)

	token := uuid.New().String()
	expiresAt := time.Now().Add(1 * time.Hour)

	newToken := &models.Token{
		Token:     token,
		Email:     email,
		ExpiresAt: expiresAt,
	}
	if err := as.TokenRepository.Add(context.Background(), newToken); err != nil {
		logger.Error("Failed to add token", "err", err.Error())
		return "", err
	}

	as.MailService.SendEmail(email, "Password Reset", fmt.Sprintf("Token for password: %s", token))
	logger.Info("Token sent succesfully", "email", email)
	return token, nil
}

func (as *AuthService) HandleResetPassword(ctx context.Context, token, newPassword string) error {
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

	_ = as.TokenRepository.Delete(ctx, foundToken.ID)
	return nil
}
