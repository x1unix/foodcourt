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
		token := r.URL.Query().Get(TokenQueryParam)

		if len(token) == 0 {
			rest.HttpErrorFromString(errAccessDenied, http.StatusForbidden).Write(&w)
			return
		}

		context.Set(r, TokenQueryParam, token)
		handler(w, r)
	}
}

func RequireAuth(handler rest.RequestHandler) rest.RequestHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get(TokenQueryParam)

		if len(token) == 0 {
			rest.HttpErrorFromString(errAccessDenied, http.StatusForbidden).Write(&w)
			return
		}

		session, valid := vault.IsSessionValid(r, token)
		if !valid {
			rest.HttpErrorFromString(errAccessDenied, http.StatusForbidden).Write(&w)
			return
		}

		sessionData, err := vault.BuildSessionData(session)

		if err != nil {
			rest.HttpErrorFromString(errAccessDenied, http.StatusForbidden).Write(&w)
			return
		}

		context.Set(r, ContextSessionKey, sessionData)
		context.Set(r, TokenQueryParam, token)

		handler(w, r)
	}
}


