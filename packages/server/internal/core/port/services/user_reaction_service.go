//go:generate mockgen -source=user_reaction_service.go -destination=../../service/mocks/user_reaction_service_mock.go -package=mocks -typed

package services

import (
	"context"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/google/uuid"
)

type UserReactionService interface {
	Service[models.UserReaction]
	GetSongLikes(ctx context.Context, songID uuid.UUID) int64
	GetSongDislikes(ctx context.Context, songID uuid.UUID) int64
}
