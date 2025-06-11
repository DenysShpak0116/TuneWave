package dto

import (
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

func NewUserDTO(user *models.User) UserDTO {
	return UserDTO{
		ID:             user.ID,
		Username:       user.Username,
		Role:           user.Role,
		ProfilePicture: user.ProfilePicture,
		ProfileInfo:    user.ProfileInfo,
		Followers:      int64(len(user.Followers)),
	}
}

// type UserExtendedDTO struct {
// 	CreatedAt   time.Time       `json:"createdAt"`
// 	Songs       []SongDTO       `json:"songs"`
// 	Collections []CollectionDTO `json:"collections"`
// 	Chats       []ChatDTO       `json:"chats"`
// 	Follows     []UserDTO       `json:"follows"`
// 	Followers   []UserDTO       `json:"followers"`

// 	ID             uuid.UUID `json:"id"`
// 	Username       string    `json:"username"`
// 	Role           string    `json:"role"`
// 	ProfileInfo    string    `json:"profileInfo"`
// 	Email          string    `json:"email"`
// 	ProfilePicture string    `json:"profilePictureUrl"`
// }
