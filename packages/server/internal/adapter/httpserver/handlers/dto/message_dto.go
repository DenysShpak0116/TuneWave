package dto

import (
	"time"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/google/uuid"
)

type MessageDTO struct {
	CreatedAt time.Time `json:"createdAt"`
	ID        uuid.UUID `json:"id"`
	SenderID  uuid.UUID `json:"senderId"`
	Content   string    `json:"content"`
}

func NewMessageDTO(message *models.Message) *MessageDTO {
	return &MessageDTO{
		ID:        message.ID,
		Content:   message.Content,
		CreatedAt: message.CreatedAt,
		SenderID:  message.SenderID,
	}
}
