package repository

import (
	"context"
	"fmt"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QueryBuilder[T any] struct {
	ctx      context.Context
	repo     *GenericRepository[T]
	query    *gorm.DB
	preloads []string
}

func (r *GenericRepository[T]) NewQuery(ctx context.Context) port.Query[T] {
	return &QueryBuilder[T]{
		ctx:   ctx,
		repo:  r,
		query: r.db.WithContext(ctx).Model(new(T)),
	}
}

func (qb *QueryBuilder[T]) Where(params any, args ...any) port.Query[T] {
	switch p := params.(type) {
	case *T:
		qb.query = qb.query.Where(p)
	case string:
		qb.query = qb.query.Where(p, args...)
	default:
		fmt.Println("unsupported parameter type for Where")
	}
	return qb
}

func (qb *QueryBuilder[T]) First(params any, args ...any) (T, error) {
	var result T
	var err error

	db := qb.query
	for _, preload := range qb.preloads {
		db = db.Preload(preload)
	}

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

	return result, err
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

func (qb *QueryBuilder[T]) Find() ([]T, error) {
	var entities []T
	db := qb.query

	for _, preload := range qb.preloads {
		db = db.Preload(preload)
	}

	err := db.Find(&entities).Error
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (qb *QueryBuilder[T]) Count() (int64, error) {
	var count int64
	err := qb.query.Count(&count).Error
	return count, err
}

func (qb *QueryBuilder[T]) Delete() error {
	return qb.query.Delete(new(T)).Error
}

func (qb *QueryBuilder[T]) Join(query string, args ...interface{}) port.Query[T] {
	qb.query = qb.query.Joins(query, args...)
	return qb
}

func (qb *QueryBuilder[T]) Group(group string) port.Query[T] {
	qb.query = qb.query.Group(group)
	return qb
}
