package dtos

import "github.com/google/uuid"

type ChatDTO struct {
	ID uuid.UUID `json:"id"`

	User1 UserDTO `json:"user1"`
	User2 UserDTO `json:"user2"`
}

type ChatExtendedDTO struct {
	ID uuid.UUID `json:"id"`

	CreatedAt string `json:"createdAt"`

	User1 UserDTO `json:"user1"`
	User2 UserDTO `json:"user2"`

	Messages []MessageDTO `json:"messages"`
}
