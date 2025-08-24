package service

import (
	"log/slog"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
)

type CommentService struct {
	*GenericService[models.Comment]
}

func NewCommentService(repo port.Repository[models.Comment], logger *slog.Logger) services.CommentService {
	return &CommentService{
		GenericService: NewGenericService(repo, logger),
	}
}
