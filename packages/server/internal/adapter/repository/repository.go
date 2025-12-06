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
	logger := r.logger.With(
		slog.String("op", op),
		slog.String("object", reflect.TypeOf(entities).Name()),
	)

	if err := r.db.WithContext(ctx).Create(entities).Error; err != nil {
		logger.Error("error while creating", "err", err.Error())
		return err
	}

	logger.Info("successfully created")
	return nil
}

func (r *GenericRepository[T]) Update(ctx context.Context, entity *T) error {
	const op = "adapter.repository.Update"
	logger := r.logger.With(
		slog.String("op", op),
		slog.String("object", reflect.TypeOf(entity).Name()),
	)

	if err := r.db.WithContext(ctx).Model(entity).Updates(entity).Error; err != nil {
		logger.Error("error while updating", "err", err.Error())
		return err
	}

	logger.Info("successfully updated")
	return nil
}

func (r *GenericRepository[T]) Delete(ctx context.Context, id ...uuid.UUID) error {
	const op = "adapter.repository.Delete"

	var entity T
	logger := r.logger.With(
		slog.String("op", op),
		slog.String("object", reflect.TypeOf(entity).Name()),
	)

	if err := r.db.WithContext(ctx).Delete(&entity, id).Error; err != nil {
		logger.Error("error while deleting", "err", err.Error())
		return err
	}

	logger.Info("successfully deleted")
	return nil
}

func (r *GenericRepository[T]) Distinct(ctx context.Context, field string) []string {
	var fieldList []string

	var entities []T
	if err := r.db.WithContext(ctx).
		Model(&entities).
		Distinct(field).
		Pluck(field, &fieldList).Error; err != nil {
		return []string{}
	}

	return fieldList
}
