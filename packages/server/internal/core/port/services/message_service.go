//go:generate mockgen -source=message_service.go -destination=../../service/mocks/message_service_mock.go -package=mocks -typed

package services

import "github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"

type MessageService interface {
	Service[models.Message]
}
