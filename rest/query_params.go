package rest

import (
	"net/url"
	"strconv"
	"net/http"
)

type QueryParamsJar struct {
	params url.Values
}

func (q *QueryParamsJar) GetSingle(key string) string {
	return q.params[key][0]
}

func (q *QueryParamsJar) Get(key string) []string {
	return q.params[key]
}

func (q *QueryParamsJar) Has(key string) bool {
	keys, ok := q.params[key]
	return ok && (len(keys) > 0)
}

func (q *QueryParamsJar) GetInt(key string) int {
	val, err := strconv.Atoi(q.GetSingle(key))

	if err != nil {
		return 0
	} else {
		return val
	}
}

// Read query params from the request
func QueryParams(request *http.Request) *QueryParamsJar {
	qp := QueryParamsJar{request.URL.Query()}

	return &qp
}

