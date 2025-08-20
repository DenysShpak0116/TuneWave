package repository

import (
	"context"
	"log/slog"
	"reflect"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type GenericRepository[T any] struct {
	db     *gorm.DB
	redis  *redis.Client
	logger *slog.Logger
}

func NewRepository[T any](db *gorm.DB, redis *redis.Client, logger *slog.Logger) port.Repository[T] {
	return &GenericRepository[T]{
		db:     db,
		redis:  redis,
		logger: logger,
	}
}

func (r *GenericRepository[T]) Add(ctx context.Context, entities ...*T) error {
	const op = "adapter.repository.Add"
	logger := r.logger.With(slog.String("op", op))

	if err := r.db.WithContext(ctx).Create(entities).Error; err != nil {
		logger.Error("error while creating %s, err: %s", reflect.TypeOf(entities).Name(), err.Error())
		return err
	}

	logger.Info("%s successfully created", reflect.TypeOf(entities).Name())
	return nil
}

func (r *GenericRepository[T]) Update(ctx context.Context, entity *T) error {
	const op = "adapter.repository.Update"
	logger := r.logger.With(slog.String("op", op))

	if err := r.db.WithContext(ctx).Model(entity).Updates(entity).Error; err != nil {
		logger.Error("error while updating %s, err: %s", reflect.TypeOf(entity).Name(), err.Error())
		return err
	}

	logger.Info("%s successfully updated", reflect.TypeOf(entity).Name())
	return nil
}

func (r *GenericRepository[T]) Delete(ctx context.Context, id ...uuid.UUID) error {
	const op = "adapter.repository.Delete"
	logger := r.logger.With(slog.String("op", op))

	var entity T
	if err := r.db.WithContext(ctx).Delete(&entity, id).Error; err != nil {
		logger.Error("error while deleting %s, err: %s", reflect.TypeOf(entity).Name(), err.Error())
		return err
	}

	logger.Info("%s successfully deleted", reflect.TypeOf(entity).Name())
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
