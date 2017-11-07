package controller

import (
	"net/http"
	"../shared/rest"
	"../shared/vault"
)

const TOKEN_PARAM = "token"

func HandleHomeGET(w http.ResponseWriter, r *http.Request) {
	rest.Echo("Hello World").Write(&w)
}
