package dto

import (
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/google/uuid"
)

type ChatDTO struct {
	ID uuid.UUID `json:"id"`

	User1 UserDTO `json:"user1"`
	User2 UserDTO `json:"user2"`
}

func (b *DTOBuilder) BuildChatDTO(chat *models.Chat) *ChatDTO {
	return &ChatDTO{
		ID:    chat.ID,
		User1: *b.BuildUserDTO(&chat.User1),
		User2: UserDTO{},
	}
}

type ChatExtendedDTO struct {
	ID uuid.UUID `json:"id"`

	CreatedAt string `json:"createdAt"`

	User1 UserDTO `json:"user1"`
	User2 UserDTO `json:"user2"`

	Messages []MessageDTO `json:"messages"`
}
