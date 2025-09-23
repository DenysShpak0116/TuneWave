package service

import (
	"log/slog"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
)

type VectorService struct {
	GenericService[models.Vector]
}

func NewVectorService(repo port.Repository[models.Vector], logger *slog.Logger) services.VectorService {
	return &VectorService{
		GenericService: GenericService[models.Vector]{
			repository: repo,
			logger:     logger,
		},
	}
}
