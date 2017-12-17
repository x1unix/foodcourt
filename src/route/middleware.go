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
const errUnauthorized = "unauthorized"
const errAccessDenied = "access denied"
const paramUserId = "userId"

// Middleware guard that requires API token to be passed
func RequireToken(handler rest.RequestHandler) rest.RequestHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		token := rest.GetToken(r)

		if len(token) == 0 {
			rest.HttpErrorFromString(errUnauthorized, http.StatusUnauthorized).Write(&w)
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
			rest.HttpErrorFromString(errUnauthorized, http.StatusUnauthorized).Write(&w)
			return
		}

		context.Set(r, TokenQueryParam, token)

		exists, err := vault.SessionExists(token)

		if err != nil {
			rest.Error(err).Write(&w)
			return
		}

		if !exists {
			rest.HttpErrorFromString(errUnauthorized, http.StatusUnauthorized).Write(&w)
			return
		}

		session, _ := vault.GetSession(token)

		// Check token ttl
		now := time.Now().Unix()

		if now >= session.TTL {
			vault.RevealToken(token)
			rest.HttpErrorFromString("session timed out", http.StatusUnauthorized).Write(&w)
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


// Middleware checks if the passed date in request is correct
func RequireValidDate(handler rest.RequestHandler) rest.RequestHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		date := rest.Params(r).GetString("date")

		// Date must be in format YYYYMMDD
		if rest.DateFormatValid(date) {
			handler(w, r)
		} else {
			rest.ErrorFromString("Invalid date format (expected: YYYYMMDD)", 400).Write(&w)
		}
	}
}

// Allow access only if client's group is equals or higher that specified
// of if the section belongs to the client's user (for ex. route like /foo/bar/{userId})
func OnlySelfOrGroup(handler rest.RequestHandler, minAccessGroup int) rest.RequestHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		token := rest.GetToken(r)
		session, _ := vault.GetSession(token)
		targetUserId := rest.Params(r).GetInt(paramUserId)
		currentUser := session.User

		// Check access
		if (currentUser.Level <= minAccessGroup) || (currentUser.ID == targetUserId) {
			handler(w, r)
		} else {
			rest.ErrorFromString("access denied", http.StatusForbidden)
		}
	}
}

// Allow access only if client is manager (or higher)
// of if the section belongs to the client's user (for ex. route like /foo/bar/{userId})
func OnlySelfOrManager(handler rest.RequestHandler) rest.RequestHandler {
	return OnlySelfOrGroup(handler, auth.LEVEL_MANAGER)
}