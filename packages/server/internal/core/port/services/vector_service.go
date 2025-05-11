package services

import (
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
)

type VectorService interface {
	Service[models.Vector]
}
