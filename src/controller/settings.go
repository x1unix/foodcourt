package controller

import (
	"../shared/settings"
	"../shared/rest"
	"encoding/json"
	"net/http"
	"fmt"
)

const errSettings = "Failed to %s settings: %v"

// Save settings
// (POST - /api/settings)
func SaveSettings(w http.ResponseWriter, r *http.Request) {
	// Settings body
	var data settings.Settings

	// Extract request data
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	defer r.Body.Close()

	if err != nil {
		rest.HttpError(err, http.StatusBadRequest).Write(&w)
		return
	}

	if err = settings.SetSettings(&data); err != nil {
		log.Error(fmt.Sprintf(errSettings, "save", err))
		rest.HttpError(err, http.StatusInternalServerError).Write(&w)
		return
	}

	rest.Ok(&w)
}

// Get settings
// (GET - /api/settings)
func GetSettings(w http.ResponseWriter, r *http.Request) {
	data := settings.GetSettings()

	rest.Success(data).Write(&w)
}