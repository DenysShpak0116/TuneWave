//go:generate mockgen -source=comment_service.go -destination=../../service/mocks/comment_service_mock.go -package=mocks -typed

package services

import (
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
)

type CommentService interface {
	Service[models.Comment]
}
