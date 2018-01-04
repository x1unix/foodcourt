package controller

import (
	"../shared/database"
	"../shared/dishes"
	"../shared/logger"
	"../shared/rest"
	"encoding/json"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

var dishValidator = validator.New()

// Get all dishes
// (GET /api/dishes)
func GetDishes(w http.ResponseWriter, r *http.Request) {
	db := database.GetInstance()
	defer db.Close()

	querySelector := database.QuerySelector{SearchKey: "label"}
	items := []dishes.Dish{}

	// Read search query params
	parseQuerySelector(&querySelector, r)

	err := dishes.All(&items, &querySelector, db)

	if err != nil {
		rest.Error(err).Write(&w)
	} else {
		rest.Success(&items).Write(&w)
	}
}

// Get dish by id
// (GET /api/dishes/{id:[0-9]+})
func GetDishById(w http.ResponseWriter, r *http.Request) {
	db := database.GetInstance()
	defer db.Close()

	id := rest.Params(r).GetString(VarID)

	// Check if dish exists
	searchErr, exists := dishes.Exists(id, db)

	if searchErr != nil {
		logger.GetLogger().Error(searchErr.Error())
		rest.Error(searchErr).Write(&w)
		return
	}

	if !exists {
		rest.NotFound(&w)
		return
	}

	var dish dishes.Dish
	err := dishes.FindById(id, &dish, db)

	if err != nil {
		logger.GetLogger().Error(err.Error())
		rest.Error(err).Write(&w)
	} else {
		rest.Success(dish).Write(&w)
	}
}

// Delete dish by id
// (DELETE /api/dishes/{id:[0-9]+})
func DeleteDishById(w http.ResponseWriter, r *http.Request) {
	db := database.GetInstance()
	defer db.Close()

	id := rest.Params(r).GetString(VarID)
	err := dishes.Delete(id, db)

	if err != nil {
		logger.GetLogger().Error(err.Error())
		rest.Error(err).Write(&w)
	} else {
		rest.Echo("OK").Write(&w)
	}
}

// Delete multiple dishes by id
// (DELETE /api/dishes)
func DeleteMultipleDishes(w http.ResponseWriter, r *http.Request) {
	// ID's
	var ids []int

	// Extract request JSON data
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&ids)
	defer r.Body.Close()

	if err != nil {
		rest.BadRequest(&w, err.Error())
		return
	}

	// Create DB connection
	db := database.GetInstance()
	defer db.Close()

	// Try to delete items
	delErr := dishes.DeleteIds(ids, db)

	// Handle errors
	if delErr != nil {
		rest.Error(delErr).Write(&w)
	} else {
		rest.Echo("OK").Write(&w)
	}
}

// Add new dish
// (POST /api/users/)
func AddDish(w http.ResponseWriter, r *http.Request) {
	// Source model
	var dish dishes.Dish

	// Extract request JSON data
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&dish)
	defer r.Body.Close()

	if err != nil {
		rest.BadRequest(&w, err.Error())
		return
	}

	// Validate data
	validErrors := dishValidator.Struct(&dish)

	if validErrors != nil {
		rest.BadRequest(&w, validErrors.Error())
		return
	}

	db := database.GetInstance()
	defer db.Close()

	// Try to create
	createErr := dishes.Add(&dish, db)

	// Handle errors
	if createErr != nil {
		rest.Error(createErr).Write(&w)
	} else {
		rest.Echo("Success").Write(&w)
	}
}
