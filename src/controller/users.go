package controller

import (
	"../shared/auth"
	"../shared/database"
	"../shared/logger"
	"../shared/rest"
	"encoding/json"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

// Why not VAR_ID? It's official GO code convention
// and I don't like it too...
const VarID = "id"

// Get user by id
// (GET /api/users/{id:[0-9]+})
func GetUserById(w http.ResponseWriter, r *http.Request) {
	con := database.GetInstance()
	defer con.Close()
	userId := rest.Params(r).GetString(VarID)

	err, data := auth.FindById(con, userId)

	if err != nil {
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

	err, data := auth.GetAll(db)

	db.Close()

	if err != nil {
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
	defer db.Close()
	userId := rest.Params(r).GetInt(VarID)

	err, ifExists := auth.IdExists(db, userId)

	// Check if error occurred on user id check
	if err != nil {
		logger.GetLogger().Error(err)
		rest.Error(err).Write(&w)
		return
	}

	// Return 404 if user isn't exists
	if !ifExists {
		rest.HttpErrorFromString("User doesn't exists", 404).Write(&w)
		return
	}

	// Extract request data
	user := auth.User{}
	decoder := json.NewDecoder(r.Body)
	decodeErr := decoder.Decode(&user)
	defer r.Body.Close()

	if decodeErr != nil {
		rest.HttpError(decodeErr, http.StatusBadRequest).Write(&w)
		return
	}

	// Modify data in DB
	createErr := auth.UpdateUser(db, &user)

	// Write error if occurred
	if createErr != nil {
		logger.GetLogger().Error(createErr)
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
	defer db.Close()
	userId := rest.Params(r).GetInt(VarID)

	err, ifExists := auth.IdExists(db, userId)

	// Check if error occurred on user id check
	if err != nil {
		logger.GetLogger().Error(err)
		rest.Error(err).Write(&w)
		return
	}

	// Return 404 if user isn't exists
	if !ifExists {
		rest.HttpErrorFromString("User doesn't exists", 404).Write(&w)
		return
	}

	delErr := auth.Delete(db, userId)

	if delErr != nil {
		rest.Error(delErr).Write(&w)
		return
	}

	rest.Echo("Success")
}

// Add new user
// (POST /api/users/)
func AddUser(w http.ResponseWriter, r *http.Request) {
	// Source model
	var user auth.User

	// Extract request data
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	defer r.Body.Close()

	if err != nil {
		rest.HttpError(err, http.StatusBadRequest).Write(&w)
		return
	}

	validate := validator.New()

	// Validate
	validErrors := validate.Struct(&user)

	if validErrors != nil {
		rest.HttpError(validErrors, http.StatusBadRequest).Write(&w)
		return
	}

	// New sql connection
	db := database.GetInstance()
	defer db.Close()

	// Check if user exists
	err, exists := auth.MailExists(db, user.Email)

	if err != nil {
		// If query error occurred - return error
		rest.Error(err).Write(&w)
		return
	}

	if exists {
		// If user already exists - return error
		rest.ErrorFromString("User already exists", http.StatusConflict).Write(&w)
		return
	}

	// Create new user
	err = auth.AddUser(db, &user)

	if err != nil {
		rest.Error(err).Write(&w)
	} else {
		rest.Echo("Success").Write(&w)
	}
}
