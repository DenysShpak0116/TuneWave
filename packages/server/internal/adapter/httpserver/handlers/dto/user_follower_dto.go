package dto

import (
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/google/uuid"
)

type UserFollowerDTO struct {
	ID       uuid.UUID `json:"id"`
	User     UserDTO   `json:"user"`
	Follower UserDTO   `json:"follower"`
}

func (b *DTOBuilder) BuildUserFollowerDTO(userFollower *models.UserFollower) UserFollowerDTO {
	return UserFollowerDTO{
		ID:       userFollower.ID,
		User:     b.BuildUserDTO(&userFollower.User),
		Follower: b.BuildUserDTO(&userFollower.Follower),
	}
}
