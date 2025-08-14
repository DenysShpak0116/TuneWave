package service

import (
	"context"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/google/uuid"
)

type UserReactionService struct {
	*GenericService[models.UserReaction]
}

func NewUserReactionService(repo port.Repository[models.UserReaction]) services.UserReactionService {
	return &UserReactionService{
		GenericService: NewGenericService(repo),
	}
}

func (svc *UserReactionService) GetSongLikes(ctx context.Context, songID uuid.UUID) int64 {
	likes, err := svc.CountWhere(ctx, &models.UserReaction{
		Type:   "like",
		SongID: songID,
	})
	if err != nil {
		likes = 0
	}
	return likes
}

func (svc *UserReactionService) GetSongDislikes(ctx context.Context, songID uuid.UUID) int64 {
	dislikes, err := svc.CountWhere(ctx, &models.UserReaction{
		Type:   "dislike",
		SongID: songID,
	})
	if err != nil {
		dislikes = 0
	}
	return dislikes
}
