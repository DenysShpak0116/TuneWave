package service

import (
	"context"
	"errors"
	"log"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/helpers/query"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GenericService[T any] struct {
	Repository port.Repository[T]
}

func NewGenericService[T any](repo port.Repository[T]) *GenericService[T] {
	return &GenericService[T]{Repository: repo}
}

func (s *GenericService[T]) Create(ctx context.Context, entities ...*T) error {
	return s.Repository.Add(ctx, entities...)
}

func (s *GenericService[T]) GetByID(ctx context.Context, id uuid.UUID, preloads ...string) (*T, error) {
	entity, err := s.Repository.NewQuery(ctx).Preload(preloads...).First(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, ErrInternal
	}

	return &entity, nil
}

func (s *GenericService[T]) Where(ctx context.Context, params *T, opts ...query.Option) ([]T, error) {
	cfg := query.Build(opts...)

	query := s.Repository.NewQuery(ctx).Where(params).Order(cfg.SortBy)
	if cfg.Limit != -1 {
		query = query.Skip(cfg.Offset).Take(cfg.Limit)
	}
	result, err := query.Preload(cfg.Preloads...).Find()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, ErrInternal
	}

	return result, nil
}

func (s *GenericService[T]) First(ctx context.Context, params *T, preloads ...string) (*T, error) {
	result, err := s.Repository.NewQuery(ctx).Preload(preloads...).First(params)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		log.Printf("service.First: %v\n", err)
		return nil, ErrInternal
	}

	return &result, nil
}

func (s *GenericService[T]) Update(ctx context.Context, entity *T) error {
	return s.Repository.Update(ctx, entity)
}

func (s *GenericService[T]) Delete(ctx context.Context, id ...uuid.UUID) error {
	return s.Repository.Delete(ctx, id...)
}

func (s *GenericService[T]) CountWhere(ctx context.Context, params *T) (int64, error) {
	entities, err := s.Repository.NewQuery(ctx).
		Where(params).Count()
	if err != nil {
		return 0, err
	}

	return entities, nil
}
