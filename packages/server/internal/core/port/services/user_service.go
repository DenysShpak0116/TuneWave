package services

import (
	"context"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/dtos"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
)

type UserService interface {
	Service[models.User]
	GetUsers(
		ctx context.Context,
		page int,
		limit int,
	) ([]dtos.UserDTO, error)
	UpdateUserPassword(email string, hashedPassword string) error
}
