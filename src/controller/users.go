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

// Update user information
// (PUT /api/users/{id:[0-9]+})
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	db := database.GetInstance()
	userId := rest.Params(r).GetInt(USERS_VAR_UID)

	mod := model.User{DB: db, ID: userId}

	err, ifExists := mod.IdExists()

	// Check if error occurred on user id check
	if (err != nil) {
		mod.Dispose()
		logger.GetLogger().Error(err)
		rest.Error(err).Write(&w)
		return
	}

	// Return 404 if user isn't exists
	if !ifExists {
		mod.Dispose()
		rest.HttpErrorFromString("User doesn't exists", 404).Write(&w)
		return
	}

	// Extract request data
	decoder := json.NewDecoder(r.Body)
	decodeErr := decoder.Decode(&mod)
	defer r.Body.Close()

	if (err != nil) {
		rest.HttpError(decodeErr, http.StatusBadRequest).Write(&w)
		return
	}

	// Modify data in DB
	createErr := mod.Update()

	// Write error if occurred
	if (createErr != nil) {
		logger.GetLogger().Error(err)
		rest.Error(createErr).Write(&w)
		return
	}

	// Write success message if everything is OK
	rest.Echo("Success").Write(&w)
}

// Delete a user
// (DELETE /api/users/{id:[0-9]+})
func DropUser(w http.ResponseWriter, r *http.Request) {
	db := database.GetInstance()
	userId := rest.Params(r).GetInt(USERS_VAR_UID)

	mod := model.User{DB: db, ID: userId}

	err, ifExists := mod.IdExists()

	// Check if error occurred on user id check
	if (err != nil) {
		mod.Dispose()
		logger.GetLogger().Error(err)
		rest.Error(err).Write(&w)
		return
	}

	// Return 404 if user isn't exists
	if !ifExists {
		mod.Dispose()
		rest.HttpErrorFromString("User doesn't exists", 404).Write(&w)
		return
	}

	delErr := mod.Delete()

	if (delErr != nil) {
		mod.Dispose()
		rest.Error(delErr).Write(&w)
		return
	}

	mod.Dispose()
	rest.Echo("Success")
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
	err = user.Create()

	user.Dispose()

	if (err != nil) {
		rest.Error(err).Write(&w)
	} else {
		rest.Echo("Success").Write(&w)
	}
}
