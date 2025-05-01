package service

import (
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
)

type UserReactionService struct {
	*GenericService[models.UserReaction]
}

func NewUserReactionService(repo port.Repository[models.UserReaction]) services.UserReactionService {
	return &UserReactionService{
		GenericService: NewGenericService(repo),
	}
}
