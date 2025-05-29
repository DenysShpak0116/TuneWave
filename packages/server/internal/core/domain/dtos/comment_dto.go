package dtos

import (
	"time"

	"github.com/google/uuid"
)

type CommentDTO struct {
	ID uuid.UUID `json:"id"`
	Author UserDTO `json:"author"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}
