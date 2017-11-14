package database

import (
	"github.com/Masterminds/squirrel"
	h "../helpers"
	"fmt"
)

// Selection query builder helper.
// Used for URL search query params
type QuerySelector struct {
	OrderKey string
	OrderDirection string
	Limit int
	Offset int
	SearchQuery string
	SearchKey string
}

// Add query selector SQL to existing query
func (qs *QuerySelector) ApplyOnSelect(query squirrel.SelectBuilder) squirrel.SelectBuilder {
	// Apply search query
	if h.NotEmpty(qs.SearchKey) && h.NotEmpty(qs.SearchQuery) {
		key := fmt.Sprintf("%s like ", qs.SearchKey)
		value := "%" + qs.SearchQuery + "%"
		query = query.Where(key, value)
	}

	// Apply order
	if h.NotEmpty(qs.OrderKey) && h.NotEmpty(qs.OrderDirection) {
		if qs.OrderDirection == "asc" || qs.OrderDirection == "desc" {
			query = query.OrderBy(qs.OrderKey + " " + qs.OrderDirection)
		}
	}

	// Apply limit and offset
	if qs.Limit > 0 {
		query = query.Limit(uint64(qs.Limit)).Offset(uint64(qs.Offset))
	}

	return query
}
