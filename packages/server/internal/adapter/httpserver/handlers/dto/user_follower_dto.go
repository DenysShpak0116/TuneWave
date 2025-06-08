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

func NewUserFollowerDTO(userFollower *models.UserFollower) *UserFollowerDTO {
	return &UserFollowerDTO{
		ID:       userFollower.ID,
		User:     *NewUserDTO(&userFollower.User),
		Follower: *NewUserDTO(&userFollower.Follower),
	}

}
