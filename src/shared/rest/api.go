package rest

import (
	"../vault"
	"github.com/gorilla/context"
	"net/http"
)

const QueryParamToken = "token"
const ContextSessionKey = "session"

func GetToken(r *http.Request) string {
	return r.URL.Query().Get(QueryParamToken)
}

// Gets session from the request's context
func GetSession(r *http.Request) *vault.Session {
	return context.Get(r, ContextSessionKey).(*vault.Session)
}
