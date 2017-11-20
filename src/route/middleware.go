package route

import (
	"../shared/auth"
	"../shared/database"
	"../shared/rest"
	"../shared/vault"
	"github.com/gorilla/context"
	"net/http"
	"strconv"
	"time"
)

const TokenQueryParam = "token"
const errAccessDenied = "Access Denied"

// Middleware guard that requires API token to be passed
func RequireToken(handler rest.RequestHandler) rest.RequestHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		token := rest.GetToken(r)

		if len(token) == 0 {
			rest.HttpErrorFromString(errAccessDenied, http.StatusUnauthorized).Write(&w)
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
			rest.HttpErrorFromString(errAccessDenied, http.StatusUnauthorized).Write(&w)
			return
		}

		context.Set(r, TokenQueryParam, token)

		exists, err := vault.SessionExists(token)

		if err != nil {
			rest.Error(err).Write(&w)
			return
		}

		if !exists {
			rest.HttpErrorFromString(errAccessDenied, http.StatusUnauthorized).Write(&w)
			return
		}

		session, _ := vault.GetSession(token)

		// Check token ttl
		now := time.Now().Unix()

		if now >= session.TTL {
			vault.RevealToken(token)
			rest.HttpErrorFromString("Session timed out", http.StatusUnauthorized).Write(&w)
			return
		}

		// Bind session to the context if everything is ok
		context.Set(r, rest.ContextSessionKey, session)

		handler(w, r)
	}
}

// Middleware guard that requires admin privileges to access
func RequireLevel(userLevel int, strict bool, handler rest.RequestHandler) rest.RequestHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		session := rest.GetSession(r)
		strUserId := strconv.Itoa(session.UserId)

		db := database.GetInstance()
		err, u := auth.FindById(db, strUserId)
		defer db.Close()

		if err != nil {
			rest.Error(err).Write(&w)
			return
		}

		match := false

		if strict {
			match = userLevel == u.Level
		} else {
			match = userLevel >= u.Level
		}

		if !match {
			rest.HttpErrorFromString(errAccessDenied, http.StatusForbidden).Write(&w)
			return
		}

		handler(w, r)
	}
}

func RequireAdmin(handler rest.RequestHandler) rest.RequestHandler {
	return RequireAuth(RequireLevel(0, true, handler))
}
