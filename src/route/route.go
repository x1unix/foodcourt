package route

import (
	"github.com/gorilla/mux"
	"../controller"
	"../shared/environment"
	"net/http"
)

func Bootstrap() *mux.Router {
	r := mux.NewRouter()

	//dirStatic := "./" + environment.DIR_PUBLIC
	dirStatic := "./" + environment.DIR_PUBLIC

	// === USERS ===

	// Get user by id
	r.HandleFunc("/api/users/{id:[0-9]+}", controller.GetUserById).Methods("GET")

	// Get all users
	r.HandleFunc("/api/users/", controller.GetUsers).Methods("GET")

	// Add new user
	r.HandleFunc("/api/users/", controller.AddUser).Methods("POST")

	// Delete a user
	r.HandleFunc("/api/users/{id:[0-9]+}", controller.DropUser).Methods("DELETE")

	// Static
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(dirStatic))))


	return r
}
