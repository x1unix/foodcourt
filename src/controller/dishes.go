package controller

import (
	"net/http"
	"../shared/database"
	"../shared/rest"
	"../shared/dishes"
	"../shared/logger"
	"encoding/json"
	"gopkg.in/go-playground/validator.v9"
)

var dishValidator = validator.New()

// Get all dishes
// (GET /api/dishes)
func GetDishes(w http.ResponseWriter, r *http.Request) {
	db := database.GetInstance()
	defer db.Close()

	items := []dishes.Dish{}

	err := dishes.All(&items, db)

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
		logger.GetLogger().Error(searchErr)
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
		logger.GetLogger().Error(err)
		rest.Error(err).Write(&w)
	} else {
		rest.Success(dish).Write(&w)
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