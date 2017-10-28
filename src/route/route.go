package route

import (
	"github.com/gorilla/mux"
	"../controller"
)

func Bootstrap() *mux.Router {
	r := mux.NewRouter()

	// GET /
	r.HandleFunc("/", controller.HandleHomeGET).Methods("GET")

	return r
}
