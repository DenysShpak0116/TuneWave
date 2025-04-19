package dtos

import "github.com/google/uuid"

type AuthorDTO struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Role string    `json:"role"`
}
