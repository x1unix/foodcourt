package route

import (
	"../controller"
	"../shared/environment"
	"../shared/rest"
	"github.com/gorilla/mux"
	"net/http"
)

func Bootstrap() *mux.Router {
	r := mux.NewRouter()

	dirStatic := environment.DIR_PUBLIC

	// == COMMON ==

	// Get server api information
	r.HandleFunc("/api", controller.GetServerInfo).Methods("GET")




	// == AUTH ==

	// Login
	r.HandleFunc("/api/login", controller.Login).Methods("POST")

	// Logout
	r.HandleFunc("/api/logout", RequireAuth(controller.Logout)).Methods("POST")

	// Session info
	r.HandleFunc("/api/session", RequireAuth(controller.GetSessionInfo)).Methods("GET")





	// === USERS ===

	// Delete a list of users
	r.HandleFunc("/api/users/purge", RequireAdmin(controller.PurgeUsers)).Methods("POST")

	// Get user by id
	r.HandleFunc("/api/users/{id:[0-9]+}", RequireAuth(controller.GetUserById)).Methods("GET")

	// Get all users
	r.HandleFunc("/api/users", RequireAdmin(controller.GetUsers)).Methods("GET")

	// Add new user
	r.HandleFunc("/api/users", RequireAdmin(controller.AddUser)).Methods("POST")

	// Delete a user
	r.HandleFunc("/api/users/{id:[0-9]+}", RequireAdmin(controller.DropUser)).Methods("DELETE")

	// Update a user
	r.HandleFunc("/api/users/{id:[0-9]+}", RequireAuth(controller.UpdateUser)).Methods("PUT")





	// == Dishes ==

	// Get all dishes
	r.HandleFunc("/api/dishes", RequireAuth(controller.GetDishes)).Methods("GET")

	// Get by id
	r.HandleFunc("/api/dishes/{id:[0-9]+}", RequireAuth(controller.GetDishById)).Methods("GET")

	// Add new dish
	r.HandleFunc("/api/dishes", RequireAdmin(controller.AddDish)).Methods("POST")

	// Delete by id
	r.HandleFunc("/api/dishes/{id:[0-9]+}", RequireAuth(controller.DeleteDishById)).Methods("DELETE")

	// Delete multiple dishes
	r.HandleFunc("/api/dishes/purge", RequireAuth(controller.DeleteMultipleDishes)).Methods("POST")






	// == Menu ==

	// Get dishes by day
	r.HandleFunc("/api/menu/{date:[0-9]+}/dishes", RequireAuth(RequireValidDate(controller.GetMenuForTheDay))).Methods("GET")

	// Add a single dish to menu for specific date
	r.HandleFunc("/api/menu/{date:[0-9]+}/dishes", RequireAdmin(RequireValidDate(controller.AddMenuItem))).Methods("PUT")

	// Update/set a list of dishes for specific date
	r.HandleFunc("/api/menu/{date:[0-9]+}/dishes", RequireAdmin(RequireValidDate(controller.SetMenuItems))).Methods("POST")

	// Clear menu for specific date
	r.HandleFunc("/api/menu/{date:[0-9]+}", RequireAdmin(RequireValidDate(controller.ClearMenu))).Methods("DELETE")

	// Get menu status
	r.HandleFunc("/api/menu/{date:[0-9]+}/status", RequireAuth(RequireValidDate(controller.GetMenuLockState))).Methods("GET")





	// == Orders ==

	// Make an order
	r.HandleFunc("/api/orders/{date:[0-9]+}/users/{userId:[0-9]+}", RequireAuth(RequireValidDate(OnlySelfOrManager(controller.OrderDishes)))).Methods("POST")

	// Get ordered dish ids
	r.HandleFunc("/api/orders/{date:[0-9]+}/users/{userId:[0-9]+}", RequireAuth(RequireValidDate(OnlySelfOrManager(controller.GetOrderedMenuItems)))).Methods("GET")

	// Delete order
	r.HandleFunc("/api/orders/{date:[0-9]+}/users/{userId:[0-9]+}", RequireAuth(RequireValidDate(OnlySelfOrManager(controller.DeleteOrder)))).Methods("DELETE")

	// Get ordered dishes
	r.HandleFunc("/api/orders/{date:[0-9]+}/users/{userId:[0-9]+}/dishes", RequireAuth(RequireValidDate(OnlySelfOrManager(controller.GetOrderedDishes)))).Methods("GET")


	// Get order stats for period
	r.HandleFunc("/api/orders/report", RequireAuth(controller.GetOrdersReport)).Methods("GET")


	// == Files ==

	// Upload an image
	r.HandleFunc("/api/media", RequireAdmin(controller.UploadFile)).Methods("POST")






	// === ETC ===

	// Serve static files
	r.PathPrefix("/").Handler(SpaFileServer(http.Dir(dirStatic), HandleNotFound))

	return r
}

func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	apiToken := rest.GetToken(r)

	if len(apiToken) > 0 {
		// If token is defined - sent API error
		rest.HttpErrorFromString("Not Found", http.StatusNotFound).Write(&w)
	} else {
		// Otherwise - redirect to SPA
		indexFile := "./" + environment.DIR_PUBLIC + "/index.html"
		http.ServeFile(w, r, indexFile)
	}
}
