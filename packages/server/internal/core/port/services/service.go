package services

import (
	"context"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/helpers/query"
	"github.com/google/uuid"
)

type Service[T any] interface {
	Create(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id uuid.UUID, preloads ...string) (*T, error)
	Where(ctx context.Context, params *T, opts ...query.Option) ([]T, error)
	Update(ctx context.Context, entity *T) (*T, error)
	Delete(ctx context.Context, id uuid.UUID) error
	CountWhere(ctx context.Context, params *T) (int64, error)
}
