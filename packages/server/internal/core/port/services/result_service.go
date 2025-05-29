package services

import (
	"context"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/dtos"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/google/uuid"
)

type ResultService interface {
	Service[models.Result]

	ProcessUserResults(
		ctx context.Context,
		userID, collectionID uuid.UUID,
		request dto.SendResultRequest,
	) ([]dtos.UserResultsDTO, error)
	GetUserResults(ctx context.Context, userID, collectionID uuid.UUID) ([]dtos.UserResultsDTO, error)
	GetCollectiveResults(ctx context.Context, collectionID uuid.UUID) (map[string]interface{}, error)
}
