package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GenericRepository[T any] struct {
	db       *gorm.DB
	preloads []string
}

func NewRepository[T any](db *gorm.DB) *GenericRepository[T] {
	return &GenericRepository[T]{db: db}
}

func (r *GenericRepository[T]) WithPreloads(preloads ...string) *GenericRepository[T] {
	newRepo := *r
	newRepo.preloads = preloads
	return &newRepo
}

func (r *GenericRepository[T]) applyPreloads(db *gorm.DB) *gorm.DB {
	for _, preload := range r.preloads {
		db = db.Preload(preload)
	}
	return db
}

func (r *GenericRepository[T]) Add(ctx context.Context, entity *T) error {
	err := r.db.WithContext(ctx).Create(entity).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *GenericRepository[T]) GetByID(ctx context.Context, id uuid.UUID) (*T, error) {
	var entity T
	query := r.applyPreloads(r.db.WithContext(ctx)).Model(&entity).Where("id = ?", id)

	err := query.First(&entity).Error
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *GenericRepository[T]) Where(ctx context.Context, params interface{}, args ...interface{}) ([]T, error) {
	var entities []T
	query := r.applyPreloads(r.db.WithContext(ctx))

	switch p := params.(type) {
	case *T:
		query = query.Where(p)
	case string:
		query = query.Where(p, args...)
	default:
		return nil, errors.New("unsupported parameter type for Where")
	}

	err := query.Find(&entities).Error
	if err != nil {
		return nil, err
	}

	// Load relations to entities

	return entities, nil
}

func (r *GenericRepository[T]) Update(ctx context.Context, entity *T) error {
	err := r.db.WithContext(ctx).Model(entity).Updates(entity).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *GenericRepository[T]) Delete(ctx context.Context, id uuid.UUID) error {
	var entity T
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *GenericRepository[T]) SkipTake(ctx context.Context, skip int, take int) (*[]T, error) {
	var entities []T
	query := r.applyPreloads(r.db.WithContext(ctx)).Offset(skip).Limit(take)

	err := query.Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return &entities, nil
}

func (r *GenericRepository[T]) CountWhere(ctx context.Context, params *T) int64 {
	var count int64
	query := r.db.WithContext(ctx).Model(new(T)).Where(params)

	query.Count(&count)
	return count
}

func (r *GenericRepository[T]) GetAll(ctx context.Context) ([]T, error) {
	var entities []T
	query := r.applyPreloads(r.db.WithContext(ctx))

	err := query.Find(&entities).Error
	if err != nil {
		return nil, err
	}

	return entities, nil
}
