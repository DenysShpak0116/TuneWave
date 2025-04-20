package dtos

import "github.com/google/uuid"

type UserDTO struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	ProfilePicture string    `json:"profilePictureUrl"`
}

type ExtendedUserDTO struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	ProfileInfo    string    `json:"profileInfo"`
	Email          string    `json:"email"`
	ProfilePicture string    `json:"profilePictureUrl"`

	CreatedAt   string          `json:"createdAt"`
	Songs       []SongDTO       `json:"songs"`
	Collections []CollectionDTO `json:"collections"`
	Chats       []ChatDTO       `json:"chats"`
}
