package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/google/uuid"
)

type ChatService struct {
	*GenericService[models.Chat]
	userRepository port.Repository[models.User]
}

func NewChatService(
	repo port.Repository[models.Chat],
	logger *slog.Logger,
	userRepository port.Repository[models.User],
) services.ChatService {
	return &ChatService{
		GenericService: NewGenericService(repo, logger),
		userRepository: userRepository,
	}
}

func (cs *ChatService) GetOrCreateGroupChat(ctx context.Context, userIDs []uuid.UUID, name string) (*models.Chat, error) {
	const op = "core.service.ChatService.GetOrCreateGroupChat"
	logger := cs.logger.With(slog.String("op", op))

	preloads := []string{"ChatUsers"}
	chats, err := cs.repository.NewQuery(ctx).Preload(preloads...).Find()
	if err != nil {
		logger.Error("Error while retrieving chats", "err", err.Error())
		return nil, err
	}

	for _, chat := range chats {
		if sameParticipants(chat.ChatUsers, userIDs) {
			logger.Info("Returned existing chat")
			return &chat, nil
		}
	}

	users, err := cs.userRepository.NewQuery(ctx).
		Where("id IN ?", userIDs).
		Find()
	if err != nil {
		logger.Error("Failed to load users", "err", err)
		return nil, err
	}

	if name == "" {
		for i, user := range users {
			if i > 2 {
				name = fmt.Sprintf("%s and others", name)
				break
			}
			if i > 0 {
				name += ", "
			}
			name += user.Username
		}
	}

	chatUsers := make([]models.ChatUser, 0, len(users))
	for _, user := range users {
		chatUsers = append(chatUsers, models.ChatUser{UserID: user.ID})
	}

	chat := &models.Chat{
		Name:      name,
		ChatUsers: chatUsers,
	}

	if err := cs.repository.Add(ctx, chat); err != nil {
		logger.Error("Failed to create new chat", "err", err)
		return nil, err
	}

	logger.Info("New chat created successfully", "users", userIDs)
	return chat, nil
}

func sameParticipants(chatUsers []models.ChatUser, ids []uuid.UUID) bool {
	if len(chatUsers) != len(ids) {
		return false
	}

	idMap := make(map[uuid.UUID]bool)
	for _, cu := range chatUsers {
		idMap[cu.UserID] = true
	}
	for _, id := range ids {
		if !idMap[id] {
			return false
		}
	}
	return true
}
