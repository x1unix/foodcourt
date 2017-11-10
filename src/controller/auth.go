package controller

import (
	"net/http"
	"../shared/database"
	"../shared/rest"
	"../shared/vault"
	"../shared/logger"
	"../shared/auth"
	"encoding/json"
	sq "github.com/Masterminds/squirrel"
)

func Login(w http.ResponseWriter, r *http.Request) {
	// Form request body
	form := auth.Credentials{}

	// Extract request data
	decoder := json.NewDecoder(r.Body)
	decodeErr := decoder.Decode(&form)
	defer r.Body.Close()

	if decodeErr != nil {
		rest.HttpError(decodeErr, http.StatusBadRequest).Write(&w)
		return
	}

	// Create DB connection instance
	db := database.GetInstance()

	// Close db at the end
	defer db.Close()

	// Search query
	query := sq.Eq{"email": form.Email, "password": auth.Password(form.Password)}

	// Try to find user with specified credentials
	searchErr, ifExists := auth.UserExists(db, query);

	// Search error
	if searchErr != nil {
		logger.GetLogger().Error(searchErr)
		rest.HttpError(searchErr, http.StatusInternalServerError).Write(&w)
		return
	}

	// return 401 authorized in case of bad credentials
	if !ifExists {
		rest.HttpErrorFromString("Wrong email or password", http.StatusUnauthorized).Write(&w)
		return
	}

	matchErr, cuser := auth.Find(db, query);

	// Check query error
	if matchErr != nil {
		logger.GetLogger().Error(matchErr)
		rest.HttpError(matchErr, http.StatusInternalServerError).Write(&w)
		return
	}

	// Generate a token
	token, err := vault.NewToken()

	if err != nil {
		rest.Error(err).Write(&w)
		return
	}

	// Create session
	sessionData, sessErr := vault.NewSession(token, cuser)

	if sessErr != nil {
		logger.GetLogger().Error(sessErr)
		rest.Error(sessErr).Write(&w)
		return
	}

	// return session data on success
	rest.Success(sessionData).Write(&w)
}

// Get session details
// (GET - /api/session)
func GetSessionInfo(w http.ResponseWriter, r *http.Request) {
	token := rest.GetToken(r)

	sess, err := vault.GetSession(token)

	if err != nil {
		rest.Error(err).Write(&w)
		return
	}

	rest.Success(sess).Write(&w)
}

// User logout
// (POST - /api/logout)
func Logout(w http.ResponseWriter, r *http.Request) {
	token := rest.GetToken(r)

	vault.RevealToken(token)

	rest.Echo("Success").Write(&w)
}