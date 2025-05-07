package domain

import (
	"net/url"
	"strconv"
)

type Pagination struct {
	Limit  int
	Offset int
}

func NewPagination(q url.Values) Pagination {
	limitDefault := 10
	offsetDefault := 0

	limitStr := q.Get("limit")
	if limitStr == "" {
		limitStr = "-1"
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		panic(err)
	}
	if limit < 0 {
		limit = limitDefault
	}

	offsetStr := q.Get("offset")
	if offsetStr == "" {
		offsetStr = "-1"
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		panic(err)
	}
	if offset < 0 {
		offset = offsetDefault
	}

	return Pagination{Offset: offset, Limit: limit}
}

type PaginationResponse[T any] struct {
	Limit  int
	Offset int
	Total  int
	Items  []T
}
