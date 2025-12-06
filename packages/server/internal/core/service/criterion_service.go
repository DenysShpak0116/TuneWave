package service

import (
	"log/slog"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
)

type CriterionService struct {
	GenericService[models.Criterion]
}

func NewCriterionService(repo port.Repository[models.Criterion], logger *slog.Logger) services.CriterionService {
	return &CriterionService{
		GenericService: GenericService[models.Criterion]{
			repository: repo,
			logger:     logger,
		},
	}
}
