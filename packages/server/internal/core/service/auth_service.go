package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
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
) *AuthService {
	return &AuthService{
		MailService:     mailService,
		TokenRepository: tokenRepository,
		UserService:     userService,
	}
}

func (as *AuthService) HandleForgotPassword(req dto.ForgotPasswordRequest) error {
	token := uuid.New().String()
	expiresAt := time.Now().Add(1 * time.Hour)

	err := as.TokenRepository.Add(context.Background(), &models.Token{
		Token:     token,
		Email:     req.Email,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return err
	}

	as.MailService.SendEmail(
		req.Email,
		"Password Reset",
		fmt.Sprintf("Click here to reset your password: http://localhost:8081/auth/reset-password?token=%s", token),
	)
	return nil
}

func (as *AuthService) HandleResetPassword(req dto.ResetPasswordRequest) error {
	tokens, err := as.TokenRepository.NewQuery(context.Background()).
		Where("token = ?", req.Token).
		Take(1).
		Find()
	if err != nil {
		return errors.New("invalid token")
	}
	token := tokens[0]

	if time.Now().After(token.ExpiresAt) {
		return errors.New("token expired")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = as.UserService.UpdateUserPassword(token.Email, string(hash))
	if err != nil {
		return err
	}

	_ = as.TokenRepository.Delete(context.TODO(), token.ID)
	return nil
}
