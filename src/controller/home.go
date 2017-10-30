package controller

import (
	"net/http"
	"../shared/rest"
)

func HandleHomeGET(w http.ResponseWriter, r *http.Request) {
	rest.Echo("Hello World").Write(&w)
}
