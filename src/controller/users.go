package controller

import (
	"net/http"
	"../shared/rest"
	"../shared/database"
	"../shared/logger"
	"../model"
	"encoding/json"
	"gopkg.in/go-playground/validator.v9"
)

const USERS_VAR_UID = "id"

// Get user by id
// (GET /api/users/{id:[0-9]+})
func GetUserById(w http.ResponseWriter, r *http.Request) {
	con := database.GetInstance()
	userId := rest.Params(r).GetString(USERS_VAR_UID)

	mod := model.User{DB: con}
	err, data := mod.FindById(userId);
	mod.Dispose()

	if (err != nil) {
		logger.GetLogger().Error(err)
		rest.Error(err).Write(&w)
	} else {
		rest.Success(data).Write(&w)
	}


}

// Select all users
// (GET /api/users/)
func GetUsers(w http.ResponseWriter, r *http.Request) {
	db := database.GetInstance()

	err, data := model.Users(db).GetAll()

	db.Close()

	if (err != nil) {
		logger.GetLogger().Error(err)
		rest.Error(err).Write(&w)
	} else {
		rest.Success(data).Write(&w)
	}
}

// Add new user
// (POST /api/users/)
func AddUser(w http.ResponseWriter, r *http.Request) {
	// Source model
	var user model.User

	// Extract request data
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	defer r.Body.Close()

	if (err != nil) {
		rest.HttpError(err, http.StatusBadRequest).Write(&w)
		return
	}

	validate := validator.New()

	// Validate
	validErrors := validate.Struct(&user)

	if (validErrors != nil) {
		rest.HttpError(validErrors, http.StatusBadRequest).Write(&w)
		return
	}

	// Assign new DB connection instance to the model
	user.DB = database.GetInstance()

	// Check if user exists
	err, exists := user.Exists()


	if (err != nil) {
		// If query error occurred - close DB and return error
		defer user.DB.Close()
		rest.Error(err).Write(&w)
		return;
	}

	if (exists) {
		// If user already exists - return error
		defer user.DB.Close()
		rest.ErrorFromString("User already exists", http.StatusConflict).Write(&w)
		return;
	}


	// Create new user
	error := user.Create()

	user.Dispose()

	if (error != nil) {
		rest.Error(error).Write(&w)
	} else {
		rest.Echo("Success").Write(&w)
	}
}
