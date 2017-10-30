package controller

import (
	"net/http"
	"../shared/rest"
	"../shared/database"
	"../model"
)

const USERS_VAR_UID = "id"

func GetUserById(w http.ResponseWriter, r *http.Request) {
	con := database.GetInstance()
	userId := rest.Params(r).GetInt(USERS_VAR_UID)

	mod := model.User{Connection: con}
	err, data := mod.FindById(userId);

	if (err != nil) {
		rest.Error(err).Write(&w)
	} else {
		rest.Success(data).Write(&w)
	}

	mod.Dispose()
}
