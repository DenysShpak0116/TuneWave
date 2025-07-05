//go:generate mockgen -source=auth_service.go -destination=../../service/mocks/auth_service_mock.go -package=mocks -typed

package services

type AuthService interface {
	HandleForgotPassword(email string) (string, error)
	HandleResetPassword(token, newPassword string) error
}
