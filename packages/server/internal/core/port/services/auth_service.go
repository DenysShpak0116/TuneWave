//go:generate mockgen -source=auth_service.go -destination=../../service/mocks/auth_service_mock.go -package=mocks -typed

package services

import "context"

type AuthService interface {
	HandleForgotPassword(email string) (string, error)
	HandleResetPassword(ctx context.Context, token, newPassword string) error
}
