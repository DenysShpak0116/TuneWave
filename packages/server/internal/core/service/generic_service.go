package service

import (
	"context"
	"fmt"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/helpers/query"
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

func (s *GenericService[T]) GetByID(ctx context.Context, id uuid.UUID, preloads ...string) (*T, error) {
	query := s.Repository.NewQuery(ctx).Where("id = ?", id).Take(1)
	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	entities, err := query.Find()
	if err != nil {
		return nil, err
	}

	if len(entities) == 0 {
		return nil, fmt.Errorf("entity with id %s not found", id)
	}
	return &entities[0], nil
}

func (s *GenericService[T]) Where(ctx context.Context, params *T, opts ...query.Option) ([]T, error) {
	cfg := query.Build(opts...)

	result, err := s.Repository.NewQuery(ctx).
		Where(params).
		Order(cfg.Sort).
		Skip(cfg.Offset).
		Take(cfg.Limit).
		Preload(cfg.Preloads...).Find()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *GenericService[T]) Update(ctx context.Context, entity *T) (*T, error) {
	return s.Repository.Update(ctx, entity)
}

func (s *GenericService[T]) Delete(ctx context.Context, id uuid.UUID) error {
	return s.Repository.Delete(ctx, id)
}

func (s *GenericService[T]) CountWhere(ctx context.Context, params *T) (int64, error) {
	entities, err := s.Repository.NewQuery(ctx).
		Where(params).Count()
	if err != nil {
		return 0, err
	}

	return entities, nil
}
