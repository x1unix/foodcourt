package rest

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// Query params reader wrapper
type QueryParams struct {
	Params map[string]string
}

func (q *QueryParams) GetString(key string) string {
	return q.Params[key]
}

func (q *QueryParams) GetInt(key string) int {
	val, err := strconv.Atoi(q.GetString(key))

	if err != nil {
		return 0
	} else {
		return val
	}
}

// Read query params from the request
func Params(request *http.Request) *QueryParams {
	qp := QueryParams{
		Params: mux.Vars(request),
	}

	return &qp
}
