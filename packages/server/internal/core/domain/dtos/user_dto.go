package dtos

import (
	"time"

	"github.com/google/uuid"
)

type UserDTO struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	ProfilePicture string    `json:"profilePictureUrl"`
}

type UserExtendedDTO struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	ProfileInfo    string    `json:"profileInfo"`
	Email          string    `json:"email"`
	ProfilePicture string    `json:"profilePictureUrl"`

	CreatedAt   time.Time       `json:"createdAt"`
	Songs       []SongDTO       `json:"songs"`
	Collections []CollectionDTO `json:"collections"`
	Chats       []ChatDTO       `json:"chats"`
}
