package rest

import "net/http"

const QueryParamToken = "token"

func GetToken(r *http.Request) string {
	return r.URL.Query().Get(QueryParamToken)
}
