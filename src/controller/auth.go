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

	if (idErr != nil) {
		logger.GetLogger().Error(idErr)
		rest.Error(idErr).Write(&w)
		user.Dispose()
		return
	}

	// Create session
	sessionData, sessErr := vault.NewSessionTicket(r, &w, user.ID)

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