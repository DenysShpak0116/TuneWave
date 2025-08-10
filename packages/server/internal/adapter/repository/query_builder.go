package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
	"github.com/redis/go-redis/v9"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QueryBuilder[T any] struct {
	ctx      context.Context
	repo     *GenericRepository[T]
	query    *gorm.DB
	preloads []string
}

func NewGenericRepository[T any](db *gorm.DB, redis *redis.Client) *GenericRepository[T] {
	return &GenericRepository[T]{
		db:    db,
		redis: redis,
	}
}

func (r *GenericRepository[T]) NewQuery(ctx context.Context) port.Query[T] {
	return &QueryBuilder[T]{
		ctx:   ctx,
		repo:  r,
		query: r.db.WithContext(ctx).Model(new(T)),
	}
}

func (qb *QueryBuilder[T]) generateCacheKey(params any, args ...any) string {
	key := fmt.Sprintf("%T:", new(T))
	switch p := params.(type) {
	case nil:
		key += "all"
	case uint, int, int64, uuid.UUID:
		key += fmt.Sprintf("%v", p)
	case string:
		key += fmt.Sprintf("%s:%v", p, args)
	case *T:
		key += fmt.Sprintf("%v", p)
	}
	for _, preload := range qb.preloads {
		key += fmt.Sprintf(":preload:%s", preload)
	}
	return key
}

func (qb *QueryBuilder[T]) First(params any, args ...any) (T, error) {
	var result T
	cacheKey := qb.generateCacheKey(params, args...)

	if qb.repo.redis != nil {
		cached, err := qb.repo.redis.Get(qb.ctx, cacheKey).Result()
		if err == nil {
			if err := json.Unmarshal([]byte(cached), &result); err == nil {
				return result, nil
			}
		}
	}

	db := qb.query
	for _, preload := range qb.preloads {
		db = db.Preload(preload)
	}

	var err error
	switch p := params.(type) {
	case nil:
		err = db.First(&result).Error
	case uint, int, int64, uuid.UUID:
		err = db.First(&result, p).Error
	case string:
		if len(args) > 0 {
			err = db.Where(p, args...).First(&result).Error
		} else {
			err = db.First(&result, p).Error
		}
	case *T:
		err = db.Where(p).First(&result).Error
	default:
		err = fmt.Errorf("unsupported parameter type for First: %T", p)
	}

	if err == nil && qb.repo.redis != nil {
		data, err := json.Marshal(result)
		if err == nil {
			qb.repo.redis.Set(qb.ctx, cacheKey, data, 10*time.Minute) // Set TTL to 10 minutes
		}
	}

	return result, err
}

func (qb *QueryBuilder[T]) Find() ([]T, error) {
	var entities []T
	cacheKey := qb.generateCacheKey(nil)

	if qb.repo.redis != nil {
		cached, err := qb.repo.redis.Get(qb.ctx, cacheKey).Result()
		if err == nil {
			if err := json.Unmarshal([]byte(cached), &entities); err == nil {
				return entities, nil
			}
		}
	}

	db := qb.query
	for _, preload := range qb.preloads {
		db = db.Preload(preload)
	}

	err := db.Find(&entities).Error
	if err != nil {
		return nil, err
	}

	if qb.repo.redis != nil {
		// Cache the result
		data, err := json.Marshal(entities)
		if err == nil {
			qb.repo.redis.Set(qb.ctx, cacheKey, data, 10*time.Minute) // Set TTL to 10 minutes
		}
	}

	return entities, nil
}

func (qb *QueryBuilder[T]) Delete() error {
	cacheKey := qb.generateCacheKey(nil)
	err := qb.query.Delete(new(T)).Error
	if err == nil && qb.repo.redis != nil {
		qb.repo.redis.Del(qb.ctx, cacheKey)
	}
	return err
}

func (qb *QueryBuilder[T]) Where(params any, args ...any) port.Query[T] {
	switch p := params.(type) {
	case *T:
		qb.query = qb.query.Where(p)
	case string:
		qb.query = qb.query.Where(p, args...)
	}
	return qb
}

func (qb *QueryBuilder[T]) Last(params any, args ...any) (T, error) {
	var result T
	var err error

	db := qb.query
	for _, preload := range qb.preloads {
		db = db.Preload(preload)
	}

	switch p := params.(type) {
	case nil:
		err = db.Last(&result).Error
	case uint, int, int64, uuid.UUID:
		err = db.Last(&result, p).Error
	case string:
		if len(args) > 0 {
			err = db.Where(p, args...).Last(&result).Error
		} else {
			err = db.Last(&result, p).Error
		}
	case *T:
		err = db.Where(p).Last(&result).Error
	default:
		err = fmt.Errorf("unsupported parameter type for Last: %T", p)
	}

	return result, err
}

func (qb *QueryBuilder[T]) Order(order string) port.Query[T] {
	qb.query = qb.query.Order(order)
	return qb
}

func (qb *QueryBuilder[T]) Skip(offset int) port.Query[T] {
	qb.query = qb.query.Offset(offset)
	return qb
}

func (qb *QueryBuilder[T]) Take(limit int) port.Query[T] {
	qb.query = qb.query.Limit(limit)
	return qb
}

func (qb *QueryBuilder[T]) Preload(preloads ...string) port.Query[T] {
	qb.preloads = append(qb.preloads, preloads...)
	return qb
}

func (qb *QueryBuilder[T]) Count() (int64, error) {
	var count int64
	err := qb.query.Count(&count).Error
	return count, err
}

func (qb *QueryBuilder[T]) Join(query string, args ...any) port.Query[T] {
	qb.query = qb.query.Joins(query, args...)
	return qb
}

func (qb *QueryBuilder[T]) Group(group string) port.Query[T] {
	qb.query = qb.query.Group(group)
	return qb
}
