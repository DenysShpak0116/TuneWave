package service

import (
	"context"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
	"github.com/google/uuid"
)

type GenericService[T any] struct {
	Repository port.Repository[T]
}

func NewGenericService[T any](repo port.Repository[T]) *GenericService[T] {
	return &GenericService[T]{Repository: repo}
}

func (s *GenericService[T]) Create(ctx context.Context, entity *T) error {
	return s.Repository.Add(ctx, entity)
}

func (s *GenericService[T]) GetByID(ctx context.Context, id uuid.UUID) (*T, error) {
	return s.Repository.GetByID(ctx, id)
}

func (s *GenericService[T]) Where(ctx context.Context, params *T) ([]T, error) {
	return s.Repository.Where(ctx, params)
}

func (s *GenericService[T]) Update(ctx context.Context, entity *T) error {
	return s.Repository.Update(ctx, entity)
}

func (s *GenericService[T]) Delete(ctx context.Context, id uuid.UUID) error {
	return s.Repository.Delete(ctx, id)
}

func (s *GenericService[T]) SkipTake(ctx context.Context, skip int, take int) (*[]T, error) {
	return s.Repository.SkipTake(ctx, skip, take)
}

func (s *GenericService[T]) CountWhere(ctx context.Context, params *T) int64 {
	return s.Repository.CountWhere(ctx, params)
}
