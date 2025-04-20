package services

import (
	"context"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/dtos"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/google/uuid"
)

type UserService interface {
	Service[models.User]
	GetFullDTOByID(ctx context.Context, id uuid.UUID) (*dtos.UserExtendedDTO, error)
	GetUsers(
		ctx context.Context,
		page int,
		limit int,
	) ([]dtos.UserDTO, error)
	UpdateUserPassword(email string, hashedPassword string) error
}
