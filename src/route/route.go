package route

import (
	"github.com/gorilla/mux"
	"../controller"
	"../shared/auth"
	"../shared/environment"
	"net/http"
)

func Bootstrap() *mux.Router {
	r := mux.NewRouter()

	dirStatic := "./" + environment.DIR_PUBLIC


	// == AUTH ==

	// Get Token
	r.HandleFunc("/api/token", controller.GetToken).Methods("GET")

	// Login
	r.HandleFunc("/api/login", auth.RequireToken(controller.Login)).Methods("POST")

	// Logout
	r.HandleFunc("/api/logout", auth.RequireAuth(controller.Logout)).Methods("POST")

	// Session info
	r.HandleFunc("/api/session", auth.RequireAuth(controller.GetSessionInfo)).Methods("GET")


	// === USERS ===

	// Get user by id
	r.HandleFunc("/api/users/{id:[0-9]+}", controller.GetUserById).Methods("GET")

	// Get all users
	r.HandleFunc("/api/users/", controller.GetUsers).Methods("GET")

	// Add new user
	r.HandleFunc("/api/users/", controller.AddUser).Methods("POST")

	// Delete a user
	r.HandleFunc("/api/users/{id:[0-9]+}", controller.DropUser).Methods("DELETE")

	// Update a user
	r.HandleFunc("/api/users/{id:[0-9]+}", controller.UpdateUser).Methods("PUT")

	// Static
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(dirStatic))))


	return r
}
