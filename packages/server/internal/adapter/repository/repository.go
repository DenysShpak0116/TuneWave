package repository

import (
	"context"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GenericRepository[T any] struct {
	db       *gorm.DB
	preloads []string
}

func NewRepository[T any](db *gorm.DB) port.Repository[T] {
	return &GenericRepository[T]{db: db}
}

func (r *GenericRepository[T]) Add(ctx context.Context, entities ...*T) error {
	err := r.db.WithContext(ctx).Create(entities).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *GenericRepository[T]) Update(ctx context.Context, entity *T) error {
	err := r.db.WithContext(ctx).Model(entity).Updates(entity).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *GenericRepository[T]) Delete(ctx context.Context, id ...uuid.UUID) error {
	var entity T
	err := r.db.WithContext(ctx).Delete(&entity, id).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *GenericRepository[T]) Distinct(ctx context.Context, field string) []string {
	var fieldList []string

	var entities []T
	err := r.db.WithContext(ctx).Model(&entities).Distinct(field).Pluck(field, &fieldList).Error
	if err != nil {
		return []string{}
	}

	return fieldList
}
