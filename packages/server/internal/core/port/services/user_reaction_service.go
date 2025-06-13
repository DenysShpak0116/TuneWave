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
