package query

import "fmt"

type Option interface {
	apply(*listParams)
}

type listParams struct {
	Preloads []string
	Sort     string
	Offset   int
	Limit    int
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

func WithSort(orderBy, sort string) Option {
	return optionFunc(func(p *listParams) {
		p.Sort = fmt.Sprintf("%s %s", orderBy, sort)
	})
}

func WithPreloads(preloads ...string) Option {
	return optionFunc(func(p *listParams) {
		p.Preloads = append(p.Preloads, preloads...)
	})
}

func Build(opts ...Option) *listParams {
	params := &listParams{
		Limit:    20,
		Offset:   0,
		Sort:     "created_at desc",
		Preloads: []string{},
	}

	for _, opt := range opts {
		opt.apply(params)
	}

	if params.Limit <= 0 {
		params.Limit = 20
	}
	if params.Offset < 0 {
		params.Offset = 0
	}
	if params.Sort == "" {
		params.Sort = "created_at desc"
	}

	return params
}
