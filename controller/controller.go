package controller

import (
	"../shared/logger"
	"../shared/rest"
	"../shared/manifest"
	"net/http"
)

const paramUserId = "userId"
const paramDate = "date"

var log = logger.GetLogger()

// For route testing purposes only
func Echo(w http.ResponseWriter, r *http.Request) {
	rest.Ok(&w)
}

// Application meta information
// (GET /api)
func GetServerInfo(w http.ResponseWriter, r *http.Request) {
	info := manifest.GetApplicationInfo()
	rest.Success(info).Write(&w)
}