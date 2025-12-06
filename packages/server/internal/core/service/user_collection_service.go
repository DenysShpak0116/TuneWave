package service

import (
	"log/slog"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
)

type UserCollectionService struct {
	GenericService[models.UserCollection]
}

func NewUserCollectionService(repo port.Repository[models.UserCollection], logger *slog.Logger) services.UserCollectionService {
	return &UserCollectionService{
		GenericService: GenericService[models.UserCollection]{
			repository: repo,
			logger:     logger,
		},
	}
}
