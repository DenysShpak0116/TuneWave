package dto

import (
	"context"

	"github.com/google/uuid"
)

type UserService interface {
	GetUserFollowersCount(ctx context.Context, userID uuid.UUID) int64
}

type SongReactionService interface {
	GetSongLikes(ctx context.Context, songID uuid.UUID) int64
	GetSongDislikes(ctx context.Context, songID uuid.UUID) int64
}

type DTOBuilder struct {
	userService         UserService
	songReactionService SongReactionService
}

func NewDTOBuilder(userSvc UserService, reactionSvc SongReactionService) *DTOBuilder {
	return &DTOBuilder{
		userService:         userSvc,
		songReactionService: reactionSvc,
	}
}
