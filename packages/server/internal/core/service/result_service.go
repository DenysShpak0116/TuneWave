package service

import (
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
)

type ResultService struct {
	GenericService[models.Result]
}

func NewResultService(repo port.Repository[models.Result]) services.ResultService {
	return &ResultService{
		GenericService: GenericService[models.Result]{
			Repository: repo,
		},
	}
}
