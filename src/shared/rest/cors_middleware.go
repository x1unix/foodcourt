package rest

import "net/http"

func AddCORS(handler RequestHandler) RequestHandler {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
