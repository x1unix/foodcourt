package controller

import (
	"net/http"
)

func HandleHomeGET(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "text/html")
	html := `
		<img src="https://http.cat/200" />
	`

	w.Write([]byte(html))
}
