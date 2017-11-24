package controller

import (
	"net/http"
	"../shared/rest"
	"../shared/orders"
	"../shared/database"
	"encoding/json"
)


// Order items from menu
// (POST /api/orders/{userId: [0-9}+/{date: [0-9]{8}+})
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
