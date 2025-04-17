package user

import (
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
)

type UserHandler struct {
	UserService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}
