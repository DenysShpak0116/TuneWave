package services

import "github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"

type UserService interface {
	Service[models.User]
	UpdateUserPassword(email string, hashedPassword string) error
}
