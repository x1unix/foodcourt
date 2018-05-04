package controller

import (
	"encoding/json"
	"fmt"
	"foodcourt/database"
	"foodcourt/dishes"
	"foodcourt/logger"
	"foodcourt/menu"
	"foodcourt/rest"
	"net/http"
	"strconv"
)

const menuParamDate = "date"

const errMenuSetFail = "Failed to set menu for date '%d': %s"

const errMenuCommon = "Date: %d, Error: %s"

const errMenuStatus = "Failed to check menu lock status - " + errMenuCommon

func SelectRangeMenu(w http.ResponseWriter, r *http.Request) {
	params := rest.QueryParams(r)

	if !params.Has("from") || !params.Has("till") {
		rest.BadRequest(&w, "No date provided")
		return
	}

	dateFrom := params.GetInt("from")
	dateTill := params.GetInt("till")

	db := database.GetInstance()
	defer db.Close()

	items, err := menu.GetMenuForPeriod(dateFrom, dateTill, db)

	if err != nil {
		rest.Error(err).Write(&w)
		return
	}

	rest.Success(items).Write(&w)

}

// Get lock status for specified dates
// (GET /api/menu/status?dates[]=...)
func GetMenusStatus(w http.ResponseWriter, r *http.Request) {

	params := rest.QueryParams(r)

	if !params.Has("dates[]") {
		rest.BadRequest(&w, "no dates specified")
		return
	}

	dates, count := params.GetBulkInt("dates[]")

	if count == 0 {
		rest.BadRequest(&w, "no valid dates provided")
		return
	}

	status, err := menu.GetMenusLockStatus(*dates)

	if err != nil {
		log.Error("error on getting menu status for %v: %v", dates, err)
		rest.Error(err).Write(&w);
	}

	rest.Success(status).Write(&w)
}

// Get list of dishes for specific day
// (GET /api/menu/{date: [0-9]{8}+}/dishes)
func GetMenuForTheDay(w http.ResponseWriter, r *http.Request) {
	// Menu date
	date := rest.Params(r).GetInt(menuParamDate)

	// Output
	var out []dishes.Dish

	// Create DB connection
	db := database.GetInstance()
	defer db.Close()

	// Try to get menu items
	err := menu.GetDishesInMenu(&out, date, db)

	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf(errMenuCommon, date, err.Error()))
		rest.Error(err).Write(&w)
		return
	}

	rest.Success(out).Write(&w)
}

// Check if menu status
// (GET /api/menu/{date: [0-9]{8}+}/status)
func GetMenuLockState(w http.ResponseWriter, r *http.Request) {
	// Menu date
	date := rest.Params(r).GetInt(menuParamDate)

	isLocked, err := menu.GetMenuLockStatus(date)

	if err != nil {
		log.Error(fmt.Sprintf(errMenuStatus, date, err))
		rest.Error(err).Write(&w)
		return
	}

	status := menu.LockStatus{Locked: isLocked}

	rest.Success(status).Write(&w)
}

func CheckMenuPermissions(date int, w *http.ResponseWriter) bool {
	// Check if menu is locked
	isLocked, err := menu.GetMenuLockStatus(date)

	if err != nil {
		log.Error(fmt.Sprintf(errMenuCommon, date, err.Error()))
		rest.Error(err).Write(w)
		return false
	}

	if isLocked {
		rest.ErrorFromString("menu and orders for this day are not editable", http.StatusBadRequest).Write(w)
		return false
	}

	return true
}

// Clear menu for specific date
// (DELETE /api/menu/{date: [0-9]{8}+})
func ClearMenu(w http.ResponseWriter, r *http.Request) {
	// Menu date
	date := rest.Params(r).GetInt(menuParamDate)

	// Check if menu is locked
	isWritable := CheckMenuPermissions(date, &w)

	if !isWritable {
		// Break if menu is in read-only mode. Response already built
		return
	}

	// Create DB connection
	db := database.GetInstance()
	defer db.Close()

	err := menu.ClearMenu(date, db)

	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf(errMenuCommon, date, err.Error()))
		rest.Error(err).Write(&w)
		return
	}

	rest.Ok(&w)
	return
}

// Update list of dishes for specific date
// (POST /api/menu/{date: [0-9]{8}+}/dishes)
func SetMenuItems(w http.ResponseWriter, r *http.Request) {
	// Menu date
	date := rest.Params(r).GetInt(menuParamDate)

	// Check if menu is locked
	isWritable := CheckMenuPermissions(date, &w)

	if !isWritable {
		// Break if menu is in read-only mode. Response already built
		return
	}

	// Request body (dish ids)
	var dishIds []int

	// Extract request JSON data
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&dishIds)
	defer r.Body.Close()

	if err != nil {
		rest.BadRequest(&w, err.Error())
		return
	}

	if len(dishIds) == 0 {
		rest.BadRequest(&w, "Collection is empty")
		return
	}

	// Create DB connection
	db := database.GetInstance()
	defer db.Close()

	// Try to set menu
	err = menu.SetDishesForDate(dishIds, date, db)

	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf(errMenuSetFail, date, err.Error()))
		rest.Error(err).Write(&w)
		return
	}

	rest.Ok(&w)
	return
}

// Add a single menu item for specific day.
// (PUT /api/menu/{date: [0-9]{8}+}/dishes)
func AddMenuItem(w http.ResponseWriter, r *http.Request) {
	// Menu date
	date := rest.Params(r).GetInt(menuParamDate)

	// Check if menu is locked
	isWritable := CheckMenuPermissions(date, &w)

	if !isWritable {
		// Break if menu is in read-only mode. Response already built
		return
	}

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
		log.Error(exErr.Error())
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
		log.Error(itemExErr.Error())
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
		log.Error(err.Error())
		rest.Error(err).Write(&w)
		return
	}

	// Write success message
	rest.Ok(&w)
}
