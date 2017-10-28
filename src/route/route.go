package route

import (
	"github.com/gorilla/mux"
	"../controller"
	"../shared/environment"
	"net/http"
)

func Bootstrap() *mux.Router {
	r := mux.NewRouter()

	dirStatic := "./" + environment.DIR_PUBLIC

	// GET /
	r.HandleFunc("/api/foo", controller.HandleHomeGET).Methods("GET")

	// Static
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(dirStatic))))


	return r
}
