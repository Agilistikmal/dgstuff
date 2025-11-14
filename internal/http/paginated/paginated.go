package paginated

import "math"

type Paginated[T any] struct {
	Page        int  `json:"page"`
	Limit       int  `json:"limit"`
	Data        []T  `json:"data"`
	TotalPages  int  `json:"total_pages"`
	TotalItems  int  `json:"total_items"`
	HasNext     bool `json:"has_next"`
	HasPrevious bool `json:"has_previous"`
}

func NewPaginated[T any](page int, limit int) *Paginated[T] {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	return &Paginated[T]{Page: page, Limit: limit}
}

func (p *Paginated[T]) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

func (p *Paginated[T]) CalculateMetadata(totalItems int64) {
	totalPages := int(math.Ceil(float64(totalItems) / float64(p.Limit)))
	p.TotalPages = int(totalPages)
	p.TotalItems = int(totalItems)
	p.HasNext = p.Page < int(totalPages)
	p.HasPrevious = p.Page > 1
}
