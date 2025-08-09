package service

import (
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
)

type CriterionService struct {
	GenericService[models.Criterion]
}

func NewCriterionService(repo port.Repository[models.Criterion]) services.CriterionService {
	return &CriterionService{
		GenericService: GenericService[models.Criterion]{
			Repository: repo,
		},
	}
}
