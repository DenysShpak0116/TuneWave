//go:generate mockgen -source=vector_service.go -destination=../../service/mocks/vector_service_mock.go -package=mocks -typed

package services

import (
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
)

type VectorService interface {
	Service[models.Vector]
}
