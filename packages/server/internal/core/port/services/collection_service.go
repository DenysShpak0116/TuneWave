//go:generate mockgen -source=collection_service.go -destination=../../service/mocks/collection_service_mock.go -package=mocks -typed

package services

import (
	"context"
	"mime/multipart"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/google/uuid"
)

type SaveCollectionParams struct {
	Title       string
	Description string
	CoverHeader *multipart.FileHeader
	Cover       multipart.File
	UserID      uuid.UUID
}

type UpdateCollectionParams struct {
	Title       string
	Description string
	CoverHeader *multipart.FileHeader
	Cover       multipart.File
	UserID      uuid.UUID
}

type CollectionService interface {
	Service[models.Collection]
	GetMany(ctx context.Context, limit, page int, sort, order string, preloads ...string) ([]models.Collection, error)
	SaveCollection(ctx context.Context, saveCollectionParams SaveCollectionParams) (*models.Collection, error)
	UpdateCollection(ctx context.Context, id uuid.UUID, updateCollectionParams UpdateCollectionParams) (*models.Collection, error)
	GetCollectionSongs(ctx context.Context, collectionID uuid.UUID, search, sortBy, order string, page, limit int) ([]models.Song, error)
}
