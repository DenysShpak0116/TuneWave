package services

import (
	"context"

	"github.com/google/uuid"
)

type Service[T any] interface {
	Create(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id uuid.UUID) (*T, error)
	Where(ctx context.Context, params *T) ([]T, error)
	Update(ctx context.Context, entity *T) (*T, error)
	Delete(ctx context.Context, id uuid.UUID) error
	CountWhere(ctx context.Context, params *T) (int64, error)
}
