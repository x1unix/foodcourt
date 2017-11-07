package controller

import (
	"net/http"
	"../shared/database"
	"../shared/rest"
	"../shared/vault"
	"../shared/logger"
	"../model"
	"encoding/json"
)

func Login(w http.ResponseWriter, r *http.Request) {
	// Form request body
	user := model.User{}

	// Extract request data
	decoder := json.NewDecoder(r.Body)
	decodeErr := decoder.Decode(&user)
	defer r.Body.Close()

	if (decodeErr != nil) {
		rest.HttpError(decodeErr, http.StatusBadRequest).Write(&w)
		return
	}

	// Create DB connection instance
	user.DB = database.GetInstance()

	// Check credentials
	matchErr, matches := user.MatchCredentials()

	// Check query error
	if (matchErr != nil) {
		logger.GetLogger().Error(matchErr)
		rest.HttpError(matchErr, http.StatusInternalServerError).Write(&w)
		user.Dispose()
		return
	}

	// return 401 authorized in case of bad credentials
	if (!matches) {
		rest.HttpErrorFromString("Wrong email or password", http.StatusUnauthorized).Write(&w)
		user.Dispose()
		return
	}

	// If everything is ok - get user by id
	idErr := user.GetId()

	// Generate a token
	token, err := vault.NewToken()

	if (err != nil) {
		rest.Error(err).Write(&w)
		return
	}

	if (idErr != nil) {
		logger.GetLogger().Error(idErr)
		rest.Error(idErr).Write(&w)
		user.Dispose()
		return
	}

	// Create session
	sessionData, sessErr := vault.NewSession(token, user.ID)

	if (sessErr != nil) {
		logger.GetLogger().Error(sessErr)
		rest.Error(sessErr).Write(&w)
		user.Dispose()
		return
	}

	// return session data on success
	rest.Success(sessionData).Write(&w)
	user.Dispose()
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