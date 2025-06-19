package dto

import (
	"time"

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
		Followers:      b.CountUserFollowers(user.ID),
	}
}

type FullUserDTO struct {
	CreatedAt   time.Time       `json:"createdAt"`
	Songs       []SongDTO       `json:"songs"`
	Collections []CollectionDTO `json:"collections"`
	Chats       []ChatDTO       `json:"chats"`
	Follows     []UserDTO       `json:"follows"`
	Followers   []UserDTO       `json:"followers"`

	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	Role           string    `json:"role"`
	ProfileInfo    string    `json:"profileInfo"`
	Email          string    `json:"email"`
	ProfilePicture string    `json:"profilePictureUrl"`
}

func (b *DTOBuilder) BuildFullUserDTO(user *models.User) FullUserDTO {
	songs := make([]SongDTO, 0)
	for _, song := range user.Songs {
		songs = append(songs, *b.BuildSongDTO(&song))
	}

	collections := make([]CollectionDTO, 0)
	for _, collection := range user.Collections {
		collections = append(collections, *b.BuildCollectionDTO(&collection))
	}

	chats := make([]ChatDTO, 0)
	for _, chat := range user.Chats1 {
		chats = append(chats, *b.BuildChatDTO(&chat))
	}
	for _, chat := range user.Chats2 {
		chats = append(chats, *b.BuildChatDTO(&chat))
	}

	follows := make([]UserDTO, 0)
	for _, follow := range user.Follows {
		follows = append(follows, *b.BuildUserDTO(&follow.User))
	}

	followers := make([]UserDTO, 0)
	for _, follower := range user.Followers {
		followers = append(followers, *b.BuildUserDTO(&follower.Follower))
	}

	return FullUserDTO{
		CreatedAt:      user.CreatedAt,
		Songs:          songs,
		Collections:    collections,
		Chats:          chats,
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
