package controller

import (
	"encoding/json"
	"fmt"
	"foodcourt/database"
	"foodcourt/dishes"
	"foodcourt/orders"
	"foodcourt/rest"
	"net/http"
	"strconv"
)

const userDebugErrorTpl = "User: %d; Date: %d; Error: %s"

const queryDateFrom = "from"
const queryDateTill = "till"

const errBadDateFormat = "bad date format"

// Order items from menu for specific user
// (POST /api/orders/{date:[0-9]+}/users/{userId:[0-9]+})
func OrderDishes(w http.ResponseWriter, r *http.Request) {
	// Extract date and user data
	params := rest.Params(r)
	date := params.GetInt(paramDate)
	targetUser := params.GetInt(paramUserId)

	// Check if menu is locked
	isWritable := CheckMenuPermissions(date, &w)

	if !isWritable {
		// Break if menu is in read-only mode. Response already built
		return
	}

	// Extract request payload
	var dishes []int

	// Extract request JSON data
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&dishes)
	defer r.Body.Close()

	// Return error on deserialization error
	if err != nil {
		rest.BadRequest(&w, err.Error())
		return
	}

	// Check if order list is empty
	if len(dishes) == 0 {
		rest.BadRequest(&w, "order list cannot be empty")
		return
	}

	// Get DB connection
	db := database.GetInstance()
	defer db.Close()

	err = orders.OrderDishes(dishes, date, targetUser, db)

	if err != nil {
		rest.Error(err).Write(&w)
		return
	}

	rest.Ok(&w)
}

// Commit bulk order for several days
// (POST /api/orders/users/{userId:[0-9]+})
func MakeBulkOrder(w http.ResponseWriter, r *http.Request) {
	// Get user id
	userId := rest.Params(r).GetInt(paramUserId)

	// Extract request payload
	var ordersBundle orders.BulkOrderBundle
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&ordersBundle)
	defer r.Body.Close()

	if err != nil {
		rest.BadRequest(&w, err.Error())
		return
	}

	// Get DB connection
	db := database.GetInstance()
	defer db.Close()

	if err = orders.BulkOrderDishes(&ordersBundle, userId, db); err != nil {
		log.Error(err.Error())
		rest.Error(err).Write(&w)
		return
	}

	rest.Ok(&w)
}

// Gets orders for period for specific user
// (GET /api/orders/users/{userId:[0-9]+}?from=YYYYMMDD&till=YYYYMMDD)
func GetOrdersForPeriod(w http.ResponseWriter, r *http.Request) {
	// Get user id
	userId := rest.Params(r).GetInt(paramUserId)

	// Extract query params
	params := rest.QueryParams(r)

	if !params.Has("from") || !params.Has("till") {
		rest.BadRequest(&w, "no start or end date provided")
		return
	}

	// Extract range points
	dateFrom := params.GetInt("from")
	dateTill := params.GetInt("till")

	db := database.GetInstance()
	defer db.Close()

	// Try to get data
	orders, err := orders.GetUserOrdersForPeriod(userId, dateFrom, dateTill, db)

	// Handle error
	if err != nil {
		log.Error("failed to get period order [%d - %d] (u:%d): %v", dateFrom, dateTill, userId, err)
		rest.Error(err).Write(&w)
		return
	}

	rest.Success(*orders).Write(&w)
}

// Get ordered dish ids for specific user
// (GET /api/orders/{date:[0-9]+}/users/{userId:[0-9]+})
func GetOrderedMenuItems(w http.ResponseWriter, r *http.Request) {
	// Get route params
	params := rest.Params(r)
	date := params.GetInt(paramDate)
	userId := params.GetInt(paramUserId)

	db := database.GetInstance()
	defer db.Close()

	var ids []int

	err := orders.GetUserOrderMenuItems(&ids, userId, date, db)

	if err != nil {
		log.Error(err.Error())
		rest.Error(err).Write(&w)
	} else {
		rest.Success(ids).Write(&w)
	}
}

// Delete order for specific user
// (GET /api/orders/{date:[0-9]+}/users/{userId:[0-9]+})
func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	// Get route params
	params := rest.Params(r)
	date := params.GetInt(paramDate)
	userId := params.GetInt(paramUserId)

	// Check if menu and orders are locked
	isWritable := CheckMenuPermissions(date, &w)

	if !isWritable {
		// Break if menu is in read-only mode. Response already built
		return
	}

	db := database.GetInstance()
	defer db.Close()

	err := orders.DeleteOrder(date, userId, db)

	if err != nil {
		log.Error(fmt.Sprintf(userDebugErrorTpl, userId, date, err.Error()))
		rest.Error(err).Write(&w)
		return
	}

	rest.Ok(&w)
}

// Get list of ordered dishes
// (GET /api/orders/{date:[0-9]+}/users/{userId:[0-9]+}/dishes)
func GetOrderedDishes(w http.ResponseWriter, r *http.Request) {
	// Get route params
	params := rest.Params(r)
	date := params.GetInt(paramDate)
	userId := params.GetInt(paramUserId)

	db := database.GetInstance()
	defer db.Close()

	// Output data
	var dishesList []dishes.Dish

	// Get data
	err := orders.GetOrderedDishes(&dishesList, date, userId, db)

	if err != nil {
		log.Error(fmt.Sprintf(userDebugErrorTpl, userId, date, err.Error()))
		rest.Error(err).Write(&w)
		return
	}

	rest.Success(&dishesList).Write(&w)
}

// Get order report for date range
// (GET /api/orders/report?from=[0-9]+&till=[0-9]+)
func GetOrdersReport(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	strDateFrom, strDateTill := queryParams.Get(queryDateFrom), queryParams.Get(queryDateTill)

	if !(rest.DateFormatValid(strDateFrom) && rest.DateFormatValid(strDateTill)) {
		rest.BadRequest(&w, errBadDateFormat)
		return
	}

	// Convert string args to int
	dateFrom, dfErr := strconv.Atoi(strDateFrom)
	dateTill, dtErr := strconv.Atoi(strDateTill)

	if (dtErr != nil) || (dfErr != nil) {
		rest.BadRequest(&w, errBadDateFormat)
		return
	}

	db := database.GetInstance()
	defer db.Close()

	var orderStats []orders.UserOrderCounter

	err := orders.GetOrderStatsBetweenDates(&orderStats, dateFrom, dateTill, db)

	if err != nil {
		log.Error(fmt.Sprintf("%s (range: %d - %d)", err.Error(), dateFrom, dateTill))
		rest.Error(err).Write(&w)
		return
	}

	formatted := orders.FormatOrdersReport(&orderStats)

	rest.Success(formatted).Write(&w)
}
