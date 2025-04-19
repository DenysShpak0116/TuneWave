package dtos

import (
	"time"

	"github.com/google/uuid"
)

type CommentDTO struct {
	ID         uuid.UUID `json:"id"`
	AuthorID   uuid.UUID `json:"authorId"`
	AuthorName string    `json:"authorName"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"createdAt"`
}
