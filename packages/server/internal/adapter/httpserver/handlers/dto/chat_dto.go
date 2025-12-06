package dto

import (
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/google/uuid"
)

type ChatDTO struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Users []UserDTO `json:"users"`
}

func (b *DTOBuilder) BuildChatDTO(chat *models.Chat) *ChatDTO {
	var users []UserDTO
	for _, cu := range chat.ChatUsers {
		users = append(users, *b.BuildUserDTO(&cu.User))
	}

	return &ChatDTO{
		ID:    chat.ID,
		Name:  chat.Name,
		Users: users,
	}
}

type ChatExtendedDTO struct {
	ID        uuid.UUID    `json:"id"`
	Name      string       `json:"name"`
	CreatedAt string       `json:"createdAt"`
	Users     []UserDTO    `json:"users"`
	Messages  []MessageDTO `json:"messages"`
}

func (b *DTOBuilder) BuildChatExtendedDTO(chat *models.Chat, messages []models.Message) *ChatExtendedDTO {
	var users []UserDTO
	for _, cu := range chat.ChatUsers {
		users = append(users, *b.BuildUserDTO(&cu.User))
	}

	var msgs []MessageDTO
	for _, m := range messages {
		msgs = append(msgs, *b.BuildMessageDTO(&m))
	}

	return &ChatExtendedDTO{
		ID:        chat.ID,
		Name:      chat.Name,
		CreatedAt: chat.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Users:     users,
		Messages:  msgs,
	}
}
