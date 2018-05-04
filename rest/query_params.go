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

// Has checks if query parameter was provided
func (q *QueryParamsJar) Has(key string) bool {
	keys, ok := q.params[key]
	return ok && (len(keys) > 0)
}

// GetBulkInt returns array of int params provided in GET request
func (q *QueryParamsJar) GetBulkInt(key string) (*[]int, int) {
	params := q.Get(key)
	output := make([]int, 0)
	count := 0

	for _, strval := range params {
		val, err := strconv.Atoi(strval)

		if err == nil {
			output = append(output, val)
			count++
		}
	}

	return &output, count
}

// GetInt returns a first value of query parameter
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

