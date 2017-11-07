package auth

import (
	"net/http"
	"../rest"
	 "../vault"
	"github.com/gorilla/context"
)

const TokenQueryParam = "token"
const ContextSessionKey = "session"
const errAccessDenied = "Access Denied"

// Middleware guard that required api token to be passed
func RequireToken(handler rest.RequestHandler) rest.RequestHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		token := rest.GetToken(r)

		if len(token) == 0 {
			rest.HttpErrorFromString(errAccessDenied, http.StatusForbidden).Write(&w)
			return
		}

		context.Set(r, TokenQueryParam, token)
		handler(w, r)
	}
}

// Middleware guard that requires user to be authorized
func RequireAuth(handler rest.RequestHandler) rest.RequestHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		token := rest.GetToken(r)

		if len(token) == 0 {
			rest.HttpErrorFromString(errAccessDenied, http.StatusForbidden).Write(&w)
			return
		}

		context.Set(r, TokenQueryParam, token)

		exists, err := vault.SessionExists(token)

		if err != nil {
			rest.Error(err).Write(&w)
			return
		}

		if !exists {
			rest.HttpErrorFromString(errAccessDenied, http.StatusForbidden).Write(&w)
			return
		}

		handler(w, r)
	}
}

