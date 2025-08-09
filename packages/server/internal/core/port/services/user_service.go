//go:generate mockgen -source=user_service.go -destination=../../service/mocks/user_service_mock.go -package=mocks -typed

package services

import (
	"context"
	"mime/multipart"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/google/uuid"
)

type UpdatePfpParams struct {
	UserID    uuid.UUID
	Pfp       multipart.File
	PfpHeader *multipart.FileHeader
}

type UserService interface {
	Service[models.User]
	GetUsers(
		ctx context.Context,
		page int,
		limit int,
	) ([]models.User, error)
	UpdateUserPassword(email string, hashedPassword string) error
	UpdateUserPfp(ctx context.Context, pfpParams UpdatePfpParams) error
	GetUserFollowersCount(ctx context.Context, userID uuid.UUID) int64
}
