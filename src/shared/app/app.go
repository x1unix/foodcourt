package app

import (
	"github.com/gorilla/mux"
	"github.com/op/go-logging"
	"time"
	"net/http"
	"database/sql"
	"../config"
	"fmt"
)

import _ "github.com/go-sql-driver/mysql"

const SQL_CON_ERR = "Failed to connect to the database ('%s@%s:%s/%s'): %s"
const SQL_CON_SUCCESS = "SQL connection established ('%s@%s:%s/%s')"

type Application struct {
	Router *mux.Router
	DB *sql.DB
	Log *logging.Logger
}

func (app *Application) Run(httpHost string, httpPost string) {
	app.initializeDatabase()

	httpAddr := httpHost + ":" + httpPost
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

// Initialize database connection
func (app *Application) initializeDatabase() {
	host, port, dbase, user, pass :=
	config.Get(config.DB_HOST, "localhost"),
	config.Get(config.DB_PORT, "3306"),
	config.Get(config.DB_NAME, "foodcourt"),
	config.Get(config.DB_USER, "foodcourt"),
	config.Get(config.DB_PASS, "")

	sqlDsn := user + ":" + pass + "@" + host + ":" + port + "/" + dbase

	db, err := sql.Open("mysql", sqlDsn)

	// Throw error if occurred
	if (err != nil) {
		errMsg := fmt.Sprintf(SQL_CON_ERR, user, host, port, dbase, err.Error())
		app.Log.Error(errMsg)
		panic(errMsg)
	}

	// Write log message
	app.Log.Info(fmt.Sprintf(SQL_CON_SUCCESS, user, host, port, dbase))

	// Save db instance
	app.DB = db


}
