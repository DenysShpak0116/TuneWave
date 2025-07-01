//go:generate mockgen -source=user_collection_service.go -destination=../../service/mocks/user_collection_service_mock.go -package=mocks -typed

package services

import "github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"

type UserCollectionService interface {
	Service[models.UserCollection]
}
