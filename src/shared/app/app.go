package app

import (
	"github.com/gorilla/mux"
	"github.com/op/go-logging"
	"time"
	"net/http"
	"fmt"
	"golang.org/x/sys/unix"
	"../database"
	"../config"
	"../vault"
)

import (
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

type Application struct {
	Router *mux.Router
	Log *logging.Logger
}

func (app *Application) Run(httpHost string, httpPost string) {
	httpAddr := httpHost + ":" + httpPost

	// Initialize sessions
	app.initializeSessions()

	// Check for DB_TEST_CONNECTION param
	checkConnection, _ := strconv.ParseBool(config.Get(config.DB_TEST_CONNECTION, "true"))

	if (checkConnection) {
		// Test db connection if required
		app.testSQLConnection()
	}

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

func (app *Application) initializeSessions() {
	sessionsPath := config.Get(config.TMP_PATH, "/tmp")

	// Check if sessions dir exists and is writable
	if unix.Access(sessionsPath, unix.W_OK) != nil {
		msg := fmt.Sprintf(
			"Failed to initialize sessions storage, directory doesn't exists or not writable: %s",
			sessionsPath,
		)

		app.Log.Error(msg)
		panic(msg)
	}

	// Get session encrypt key
	sessKey := config.Get(config.APP_KEY, "")

	// Check if key is empty
	if (len(sessKey) == 0) {
		panic("Application key is empty. Please define application key at 'APP_KEY' parameter in .env file")
	}

	// Bootstrap sessions vault
	vault.Bootstrap(sessionsPath, sessKey)
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
