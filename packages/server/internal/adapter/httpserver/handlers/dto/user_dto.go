package dto

import (
	"context"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/google/uuid"
)

type UserDTO struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	Role           string    `json:"role"`
	ProfilePicture string    `json:"profilePictureUrl"`
	ProfileInfo    string    `json:"profileInfo"`
	Followers      int64     `json:"followers"`
}

func (b *DTOBuilder) BuildUserDTO(user *models.User) *UserDTO {
	return &UserDTO{
		ID:             user.ID,
		Username:       user.Username,
		Role:           user.Role,
		ProfilePicture: user.ProfilePicture,
		ProfileInfo:    user.ProfileInfo,
		Followers:      b.userService.GetUserFollowersCount(context.Background(), user.ID),
	}
}

type FullUserDTO struct {
	Follows   []UserDTO `json:"follows"`
	Followers []UserDTO `json:"followers"`

	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	Role           string    `json:"role"`
	ProfileInfo    string    `json:"profileInfo"`
	Email          string    `json:"email"`
	ProfilePicture string    `json:"profilePictureUrl"`
}

func (b *DTOBuilder) BuildFullUserDTO(user *models.User) FullUserDTO {
	follows := make([]UserDTO, 0)
	for _, follow := range user.Follows {
		follows = append(follows, *b.BuildUserDTO(&follow.User))
	}

	followers := make([]UserDTO, 0)
	for _, follower := range user.Followers {
		followers = append(followers, *b.BuildUserDTO(&follower.Follower))
	}

	return FullUserDTO{
		Follows:        follows,
		Followers:      followers,
		ID:             user.ID,
		Username:       user.Username,
		Role:           user.Role,
		ProfileInfo:    user.ProfileInfo,
		Email:          user.Email,
		ProfilePicture: user.ProfileInfo,
	}
}
