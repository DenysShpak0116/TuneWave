package repository

import (
	"context"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
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

func (qb *QueryBuilder[T]) Where(params interface{}, args ...interface{}) port.Query[T] {
	switch p := params.(type) {
	case *T:
		qb.query = qb.query.Where(p)
	case string:
		qb.query = qb.query.Where(p, args...)
	default:
		panic("unsupported parameter type for Where")
	}
	return qb
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
