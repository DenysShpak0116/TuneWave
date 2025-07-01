//go:generate mockgen -source=user_follower_service.go -destination=../../service/mocks/user_follower_service_mock.go -package=mocks -typed

package services

import "github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"

type UserFollowerService interface {
	Service[models.UserFollower]
}
