package backend

import (
	"babel/openapi/gen/babelapi"
)

const MaxPageSize = 50

func normalizePagination(p *babelapi.Pagination) *babelapi.Pagination {
	page := p.Page
	if page < 0 {
		page = 0
	}

	pageSize := p.PageSize
	if pageSize <= 0 {
		pageSize = 1
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	return &babelapi.Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}
