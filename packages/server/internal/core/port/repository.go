package port

import (
	"context"

	"github.com/google/uuid"
)

type Query[T any] interface {
	Where(params interface{}, args ...interface{}) Query[T]
	Order(order string) Query[T]
	Skip(offset int) Query[T]
	Take(limit int) Query[T]
	Preload(preloads ...string) Query[T]
	Find() ([]T, error)
	Count() (int64, error)
}

type Repository[T any] interface {
	Add(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) (*T, error)
	Delete(ctx context.Context, id uuid.UUID) error
	NewQuery(ctx context.Context) Query[T]
}
