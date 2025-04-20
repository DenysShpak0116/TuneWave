package dtos

import "github.com/google/uuid"

type MessageDTO struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt string    `json:"createdAt"`

	Content string `json:"content"`

	Sender UserDTO `json:"sender"`
}
