package service

import (
	"context"
	"log/slog"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/google/uuid"
)

type UserReactionService struct {
	*GenericService[models.UserReaction]
}

func NewUserReactionService(repo port.Repository[models.UserReaction], logger *slog.Logger) services.UserReactionService {
	return &UserReactionService{
		GenericService: NewGenericService(repo, logger),
	}
}

func (svc *UserReactionService) GetSongLikes(ctx context.Context, songID uuid.UUID) int64 {
	const op = "core.service.UserReactionService.GetSongLikes"
	logger := svc.logger.With(
		slog.String("op", op),
	)

	likes, err := svc.CountWhere(ctx, &models.UserReaction{
		Type:   "like",
		SongID: songID,
	})
	if err != nil {
		logger.Error("Failed to get song likes", "id", songID.String(), "err", err.Error())
		likes = 0
	}

	logger.Info("Song likes successfully retrieved", "id", songID.String(), "err", err.Error())
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
