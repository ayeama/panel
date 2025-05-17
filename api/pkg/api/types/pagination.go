package types

// total_pages: Math.ceil(total / limit)
// current_page: (offset / limit) + 1
// next: (offset + limit) < total
// previous: offset > 0
type PaginationResponse[T any] struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
	Items  []T `json:"items"`
}
