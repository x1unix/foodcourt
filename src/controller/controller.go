package controller

import (
	"../shared/logger"
	"../shared/rest"
	"net/http"
)

const paramUserId = "userId"
const paramDate = "date"

var log = logger.GetLogger()

// For route testing purposes only
func Echo(w http.ResponseWriter, r *http.Request) {
	rest.Ok(&w)
}