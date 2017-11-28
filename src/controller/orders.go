package controller

import (
	"net/http"
	"../shared/rest"
	"../shared/orders"
	"../shared/database"
	"encoding/json"
)


// Order items from menu for specific user
// (POST /api/orders/{date:[0-9]+}/users/{userId:[0-9]+})
func OrderDishes(w http.ResponseWriter, r *http.Request) {
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

	// Extract date and user data
	params := rest.Params(r)
	date := params.GetInt(paramDate)
	targetUser := params.GetInt(paramUserId)

	// Create DB connection
	db := database.GetInstance()
	defer db.Close()

	err = orders.OrderDishes(dishes, date, targetUser, db)

	if err != nil {
		rest.Error(err).Write(&w)
		return
	}

	rest.Ok(&w)
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
