package services

import (
	"context"
)

type Service[T any] interface {
	Create(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id uint) (*T, error)
	Where(ctx context.Context, params *T) ([]T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id uint) error
	SkipTake(ctx context.Context, skip int, take int) (*[]T, error)
	CountWhere(ctx context.Context, params *T) int64
}
