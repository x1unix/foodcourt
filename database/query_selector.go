package database

import (
	"fmt"
	h "foodcourt/helpers"

	"github.com/Masterminds/squirrel"
)

// Selection query builder helper.
// Used for URL search query params
type QuerySelector struct {
	OrderKey       string
	OrderDirection string
	Limit          int
	Offset         int
	SearchQuery    string
	SearchKey      string
}

// Add query selector SQL to existing query
func (qs *QuerySelector) ApplyOnSelect(query squirrel.SelectBuilder) squirrel.SelectBuilder {
	// Apply search query
	if h.NotEmpty(qs.SearchKey) && h.NotEmpty(qs.SearchQuery) {
		// Doesn't work - sql: converting Exec argument $1 type: unsupported type []interface {}, a slice of interface
		//
		// key := fmt.Sprintf("%s LIKE ?", qs.SearchKey)
		// val := fmt.Sprint("%", qs.SearchQuery, "%")

		// query = query.Where(key, val)

		// So I've used a workaround:
		key := fmt.Sprintf("%s LIKE ", qs.SearchKey)
		val := fmt.Sprint("%", qs.SearchQuery, "%")
		val = h.EscapeString(val)

		query = query.Where(key + val)
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
