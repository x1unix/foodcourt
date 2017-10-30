package app

import (
	"github.com/gorilla/mux"
	"github.com/op/go-logging"
	"time"
	"net/http"
	"fmt"
	"../database"
)

import _ "github.com/go-sql-driver/mysql"

type Application struct {
	Router *mux.Router
	Log *logging.Logger
}

func (app *Application) Run(httpHost string, httpPost string) {
	httpAddr := httpHost + ":" + httpPost
	app.testSQLConnection()
	app.initializeHTTPServer(httpAddr)
}

// Initialize HTTP server
func (app *Application) initializeHTTPServer(httpAddr string) {
	server := &http.Server{
		Handler: app.Router,
		Addr: httpAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	app.Log.Info(fmt.Sprintf("HTTP server started at %s", httpAddr))

	app.Log.Fatal(server.ListenAndServe())
}

func (app *Application) testSQLConnection() {
	app.Log.Info("Testing connection to the MySQL database...")

	result, err := database.TestConnection()

	if (!result) {
		msg := fmt.Sprintf("Failed to connect to the database: %s. Application will be terminated.", err.Error())
		app.Log.Error(msg)
		panic(msg)
	} else {
		app.Log.Info("Connection successful, everything is fine.")
	}

}
