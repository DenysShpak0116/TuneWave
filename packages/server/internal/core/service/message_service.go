package service

import (
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
)

type MessageService struct {
	GenericService[models.Message]
}

func NewMessageService(repo port.Repository[models.Message]) services.MessageService {
	return &MessageService{
		GenericService: GenericService[models.Message]{
			Repository: repo,
		},
	}
}
