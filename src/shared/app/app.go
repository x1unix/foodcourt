package app

import (
	"github.com/gorilla/mux"
	"database/sql"
	"github.com/op/go-logging"
	"time"
	"net/http"
)

type Application struct {
	Router *mux.Router
	DB *sql.DB
	Log *logging.Logger
}

func (app *Application) Run(httpHost string, httpPost string) {
	httpAddr := httpHost + ":" + httpPost

	server := &http.Server{
		Handler: app.Router,
		Addr: httpAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	app.Log.Fatal(server.ListenAndServe())
}
