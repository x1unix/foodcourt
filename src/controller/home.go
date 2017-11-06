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

func GetToken(w http.ResponseWriter, r *http.Request) {
	token, err := vault.NewToken()

	if (err != nil) {
		rest.Error(err).Write(&w)
		return
	}

	rest.Echo(token).Write(&w)
}

func TestToken(w http.ResponseWriter, r *http.Request) {

}