package rest

import (
	"net/http"
	"../vault"
	"github.com/gorilla/context"
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