package services

import "github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"

type AuthService interface {
	HandleForgotPassword(req dto.ForgotPasswordRequest) error
	HandleResetPassword(req dto.ResetPasswordRequest) error
}
