package port

import (
	"context"

	"github.com/google/uuid"
)

type Repository[T any] interface {
	Add(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id uuid.UUID) (*T, error)
	GetAll(ctx context.Context) ([]T, error)
	Where(ctx context.Context, params interface{}, args ...interface{}) ([]T, error)
	Update(ctx context.Context, entity *T) (*T, error)
	Delete(ctx context.Context, id uuid.UUID) error
	SkipTake(ctx context.Context, skip int, take int) (*[]T, error)
	CountWhere(ctx context.Context, params *T) int64
}
