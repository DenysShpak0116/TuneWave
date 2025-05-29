package dtos

import (
	"time"

	"github.com/google/uuid"
)

type MessageDTO struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`

	Content string `json:"content"`

	SenderID uuid.UUID `json:"senderId"`
}
