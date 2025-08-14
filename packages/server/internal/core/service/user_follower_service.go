package service

import (
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
)

type UserFollowerService struct {
	GenericService[models.UserFollower]
}

func NewUserFollowerService(repo port.Repository[models.UserFollower]) services.UserFollowerService {
	return &UserFollowerService{
		GenericService: GenericService[models.UserFollower]{
			Repository: repo,
		},
	}
}
