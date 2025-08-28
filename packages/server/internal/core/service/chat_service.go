package service

import (
	"context"
	"log/slog"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/google/uuid"
)

type ChatService struct {
	*GenericService[models.Chat]
}

func NewChatService(repo port.Repository[models.Chat], logger *slog.Logger) services.ChatService {
	return &ChatService{
		GenericService: NewGenericService(repo, logger),
	}
}

func (cs *ChatService) GetOrCreatePrivateChat(ctx context.Context, user1, user2 uuid.UUID) (*models.Chat, error) {
	const op = "core.service.ChatService.GetOrCreatePrivateChat"
	logger := cs.logger.With(
		slog.String("op", op),
	)

	chats, err := cs.repository.NewQuery(ctx).
		Where("(user_id1 = ? AND user_id2 = ?) OR (user_id1 = ? AND user_id2 = ?)", user1, user2, user2, user1).
		Find()
	if err != nil {
		logger.Error("Error while trying to retrieve chats", "err", err.Error())
		return nil, err
	}
	if len(chats) >= 1 {
		logger.Info("Returned first found chat")
		return &chats[0], nil
	}

	chat := &models.Chat{UserID1: user1, UserID2: user2}
	if err = cs.repository.Add(ctx, chat); err != nil {
		logger.Error("Failed to create new chat", "err", err)
		return nil, err
	}

	logger.Info("New chat created successfully", "user1", user1.String(), "user2", user2.String())
	return chat, nil
}
