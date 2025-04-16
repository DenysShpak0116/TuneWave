package service

import (
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
)

type UserService struct {
	*GenericService[models.User]
}

func NewUserService(repo port.Repository[models.User]) *UserService {
	return &UserService{
		GenericService: NewGenericService(repo),
	}
}
