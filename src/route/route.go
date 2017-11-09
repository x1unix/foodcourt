package route

import (
	"github.com/gorilla/mux"
	"../controller"
	"../shared/auth"
	"../shared/environment"
	"../shared/rest"
	"net/http"
)

func Bootstrap() *mux.Router {
	r := mux.NewRouter()

	dirStatic := environment.DIR_PUBLIC


	// == AUTH ==

	// Login
	r.HandleFunc("/api/login", controller.Login).Methods("POST")

	// Logout
	r.HandleFunc("/api/logout", auth.RequireAuth(controller.Logout)).Methods("POST")

	// Session info
	r.HandleFunc("/api/session", auth.RequireAuth(controller.GetSessionInfo)).Methods("GET")


	// === USERS ===

	// Get user by id
	r.HandleFunc("/api/users/{id:[0-9]+}", auth.RequireAuth(controller.GetUserById)).Methods("GET")

	// Get all users
	r.HandleFunc("/api/users", auth.RequireAdmin(controller.GetUsers)).Methods("GET")

	// Add new user
	r.HandleFunc("/api/users", auth.RequireAdmin(controller.AddUser)).Methods("POST")

	// Delete a user
	r.HandleFunc("/api/users/{id:[0-9]+}", auth.RequireAdmin(controller.DropUser)).Methods("DELETE")

	// Update a user
	r.HandleFunc("/api/users/{id:[0-9]+}", auth.RequireAuth(controller.UpdateUser)).Methods("PUT")

	// == Dishes ==

	// Get all dishes
	r.HandleFunc("/api/dishes", auth.RequireAuth(controller.GetDishes)).Methods("GET")

	// Get by id
	r.HandleFunc("/api/dishes/{id:[0-9]+}", auth.RequireAuth(controller.GetDishById)).Methods("GET")

	// Add new dish
	r.HandleFunc("/api/dishes", auth.RequireAdmin(controller.AddDish)).Methods("POST")



	// === ETC ===

	// Serve static files
	r.PathPrefix("/").Handler(SpaFileServer(http.Dir(dirStatic), HandleNotFound))

	return r
}

func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	apiToken := rest.GetToken(r)

	if (len(apiToken) > 0) {
		// If token is defined - sent API error
		rest.HttpErrorFromString("Not Found", http.StatusNotFound).Write(&w)
	} else {
		// Otherwise - redirect to SPA
		indexFile := "./" + environment.DIR_PUBLIC + "/index.html"
		http.ServeFile(w, r, indexFile)
	}
}
