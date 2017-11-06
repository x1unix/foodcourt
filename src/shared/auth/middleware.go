package auth

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
)

// List of ignored routes by auth middleware
var ProtectedRoutes = []string{
	"/api/token",
}

func AuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mr := mux.CurrentRoute(r)
		if (mr != nil) {
			fmt.Println(mr.GetPathTemplate())
		}
		h.ServeHTTP(w, r)
	})
}


