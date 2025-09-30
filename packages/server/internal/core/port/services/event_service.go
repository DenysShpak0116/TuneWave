//go:generate mockgen -source=event_service.go -destination=../../service/mocks/event_service_mock.go -package=mocks -typed

package services

import (
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
)

type EventService interface {
	Service[models.Event]
}
