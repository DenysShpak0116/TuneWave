package service

import (
	"log/slog"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
)

type EventService struct {
	GenericService[models.Event]
}

func NewEventService(repo port.Repository[models.Event], logger *slog.Logger) services.EventService {
	return &EventService{
		GenericService: GenericService[models.Event]{
			repository: repo,
			logger:     logger,
		},
	}
}
