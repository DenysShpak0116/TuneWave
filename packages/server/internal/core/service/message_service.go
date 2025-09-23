package service

import (
	"log/slog"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
)

type MessageService struct {
	GenericService[models.Message]
}

func NewMessageService(repo port.Repository[models.Message], logger *slog.Logger) services.MessageService {
	return &MessageService{
		GenericService: GenericService[models.Message]{
			repository: repo,
			logger: logger,
		},
	}
}
