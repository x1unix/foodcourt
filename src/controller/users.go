package controller

import (
	"net/http"
	"../shared/rest"
	"../shared/database"
	"../shared/logger"
	"../model"
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
