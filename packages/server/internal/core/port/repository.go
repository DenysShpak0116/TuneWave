package port

import (
	"context"

	"github.com/google/uuid"
)

type Query[T any] interface {
	Where(params any, args ...any) Query[T]
	First(params any, args ...any) (T, error)
	Last(params any, args ...any) (T, error)
	Order(order string) Query[T]
	Skip(offset int) Query[T]
	Take(limit int) Query[T]
	Preload(preloads ...string) Query[T]
	Delete() error
	Find() ([]T, error)
	Count() (int64, error)
	Join(query string, args ...any) Query[T]
	Group(group string) Query[T]
}

type Repository[T any] interface {
	Add(ctx context.Context, entities ...*T) error
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id ...uuid.UUID) error
	Distinct(ctx context.Context, field string) []string
	NewQuery(ctx context.Context) Query[T]
}
