package controller

import (
	"net/http"
)

func HandleHomeGET(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "text/html")
	html := `
		<b>Hello</b>
	`

	w.Write([]byte(html))
}
