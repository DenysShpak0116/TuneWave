package services

import (
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
)

type ResultService interface {
	Service[models.Result]
}
