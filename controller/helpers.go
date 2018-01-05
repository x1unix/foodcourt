package controller

import (
	"foodcourt/database"
	"net/http"
	"strconv"
)

const qsKeySearchQuery = "q"
const qsKeyLimit = "limit"
const qsKeyOffset = "offset"
const qsKeyOrderBy = "order_by"
const qsKeyOrderDir = "order_dir"

// Parse data query selector from query params.
// URL ex: /items?q=some+query&limit=32&offset=2
func parseQuerySelector(qs *database.QuerySelector, r *http.Request) {
	queryParams := r.URL.Query()

	qs.SearchQuery = queryParams.Get(qsKeySearchQuery)
	qs.OrderKey = queryParams.Get(qsKeyOrderBy)
	qs.OrderDirection = queryParams.Get(qsKeyOrderDir)

	limit, limitErr := strconv.Atoi(queryParams.Get(qsKeyLimit))
	offset, offsetErr := strconv.Atoi(queryParams.Get(qsKeyOffset))

	if (limitErr == nil) && (offsetErr == nil) {
		qs.Limit = limit
		qs.Offset = offset
	}
}
