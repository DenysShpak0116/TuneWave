package services

import (
	"context"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/google/uuid"
)

type ChatService interface {
	Service[models.Chat]
	GetOrCreatePrivateChat(ctx context.Context, user1, user2 uuid.UUID) (*models.Chat, error)
}
