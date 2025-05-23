package dtos

import (
	"time"

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

type UserExtendedDTO struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	Role           string    `json:"role"`
	ProfileInfo    string    `json:"profileInfo"`
	Email          string    `json:"email"`
	ProfilePicture string    `json:"profilePictureUrl"`

	CreatedAt   time.Time       `json:"createdAt"`
	Songs       []SongDTO       `json:"songs"`
	Collections []CollectionDTO `json:"collections"`
	Chats       []ChatDTO       `json:"chats"`
	Follows     []UserDTO       `json:"follows"`
	Followers   []UserDTO       `json:"followers"`
}
