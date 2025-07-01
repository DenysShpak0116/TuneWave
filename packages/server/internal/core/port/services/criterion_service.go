//go:generate mockgen -source=criterion_service.go -destination=../../service/mocks/criterion_service_mock.go -package=mocks -typed

package services

import (
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
)

type CriterionService interface {
	Service[models.Criterion]
}
