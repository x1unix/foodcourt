package rest

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// Query params reader wrapper
type ParamsJar struct {
	Params map[string]string
}

func (q *ParamsJar) GetString(key string) string {
	return q.Params[key]
}

func (q *ParamsJar) Has(key string) bool {
	_, ok := q.Params[key]
	return ok
}

func (q *ParamsJar) GetInt(key string) int {
	val, err := strconv.Atoi(q.GetString(key))

	if err != nil {
		return 0
	} else {
		return val
	}
}

// Read query params from the request
func Params(request *http.Request) *ParamsJar {
	qp := ParamsJar{
		Params: mux.Vars(request),
	}

	return &qp
}
