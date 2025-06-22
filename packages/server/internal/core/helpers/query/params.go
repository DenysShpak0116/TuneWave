package query

import (
	"fmt"
	"slices"
)

type Option interface {
	apply(*listParams)
}

type listParams struct {
	allowedFields []string
	Preloads      []string
	SortBy        string
	orderBy       string
	sort          string
	defaultField  string
	Offset        int
	Limit         int
}
type optionFunc func(*listParams)

func (f optionFunc) apply(p *listParams) {
	f(p)
}

func WithPagination(page, limit int) Option {
	return optionFunc(func(p *listParams) {
		if limit <= 0 || limit > 100 {
			limit = 20
		}
		if page <= 0 {
			page = 1
		}
		p.Limit = limit
		p.Offset = (page - 1) * limit
	})
}

func WithSort(
	orderBy, sort string,
	allowedFields []string,
	defaultField string,
) Option {
	return optionFunc(func(p *listParams) {
		p.sort = sort
		p.orderBy = orderBy
		p.allowedFields = allowedFields
		p.defaultField = defaultField
	})
}

func WithPreloads(preloads ...string) Option {
	return optionFunc(func(p *listParams) {
		p.Preloads = append(p.Preloads, preloads...)
	})
}

func Build(opts ...Option) *listParams {
	params := &listParams{
		Limit:         -1,
		Offset:        0,
		SortBy:        "created_at desc",
		Preloads:      []string{},
		allowedFields: []string{"created_at"},
		defaultField:  "created_at",
	}

	for _, opt := range opts {
		opt.apply(params)
	}

	if params.Limit >= 0 && params.Limit != -1 {
		params.Limit = 20
	}
	if params.Offset < 0 {
		params.Offset = 0
	}

	if params.sort != "asc" && params.sort != "desc" {
		params.sort = "desc"
	}

	if len(params.allowedFields) == 0 {
		params.allowedFields = []string{"created_at"}
	}
	if !slices.Contains(params.allowedFields, params.defaultField) {
		params.defaultField = "created_at"
	}
	if !slices.Contains(params.allowedFields, params.orderBy) {
		params.orderBy = params.defaultField
	}

	params.SortBy = fmt.Sprintf("%s %s", params.orderBy, params.sort)

	return params
}
