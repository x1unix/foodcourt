package controller

import (
	"net/http"
	"encoding/json"
	"strconv"
	"../shared/logger"
	"../shared/rest"
	"../shared/database"
	"../shared/menu"
	"../shared/dishes"
	"fmt"
)

const menuParamDate = "date"

// Add a single menu item for specific day.
// (POST /api/menu/{date: [0-9]{8}+}/items)
func AddMenuItem(w http.ResponseWriter, r *http.Request) {
	// Menu date
	date := rest.Params(r).GetInt(menuParamDate)

	// Request body
	var menuItem menu.MenuItem

	// Extract request JSON data
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&menuItem)
	defer r.Body.Close()

	if err != nil {
		rest.BadRequest(&w, err.Error())
		return
	}

	// Create DB connection
	db := database.GetInstance()
	defer db.Close()

	// Get logger instance
	log := logger.GetLogger()

	// Extract and validate dish id
	dishId := menuItem.DishId

	// Check if dish exists
	exErr, exists := dishes.Exists(strconv.Itoa(dishId), db)

	// Handle exists query error
	if exErr != nil {
		log.Error(exErr)
		rest.Error(exErr).Write(&w)
		return
	}

	// Drop error if item not exists
	if !exists {
		rest.ErrorFromString(fmt.Sprintf("Provided dish ID doesn't exists: %d", dishId), http.StatusBadRequest).
			Write(&w)
		return
	}

	// Check if item already added to the menu
	itemExErr, itemExists := menu.DishExistsInMenu(dishId, date, db)

	// Check query errors
	if itemExErr != nil {
		log.Error(itemExErr)
		rest.Error(itemExErr).Write(&w)
		return
	}

	// handle item exists
	if itemExists {
		rest.ErrorFromString(fmt.Sprintf("Provided dish ID already exists in menu: %d", dishId), http.StatusBadRequest).
			Write(&w)
		return
	}

	// And only now, we can try to add the item
	err = menu.AddDish(dishId, date, db)

	// Handle errors
	if err != nil {
		log.Error(err)
		rest.Error(err).Write(&w)
		return
	}

	// Write success message
	rest.Echo("OK").Write(&w)
}