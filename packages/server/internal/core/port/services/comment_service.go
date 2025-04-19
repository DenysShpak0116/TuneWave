package services

import (
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
)

type CommentService interface {
	Service[models.Comment]
}
