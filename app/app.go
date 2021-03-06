package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/op/go-logging"
	"net/http"
	"time"
	"foodcourt/config"
	"foodcourt/database"
	"foodcourt/cache"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

type Application struct {
	Router *mux.Router
	Log    *logging.Logger
}

// Run app
func (app *Application) Run(httpHost string, httpPort string) {
	httpAddr := httpHost + ":" + httpPort

	// Check app key
	app.checkAppKey()

	// Initialize sessions
	app.initializeSessions()

	// Check for DB_TEST_CONNECTION param
	checkConnection, _ := strconv.ParseBool(config.Get(config.DB_TEST_CONNECTION, "true"))

	if checkConnection {
		// Test db connection if required
		app.testSQLConnection()
	}

	app.initializeCache(checkConnection)

	app.initializeHTTPServer(httpAddr)
}

// Initialize HTTP server
func (app *Application) initializeHTTPServer(httpAddr string) {
	server := &http.Server{
		Handler:      app.Router,
		Addr:         httpAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	app.Log.Info(fmt.Sprintf("HTTP server started at %s", httpAddr))

	app.Log.Fatal(server.ListenAndServe())
}

func (app *Application) checkAppKey() {
	key := config.Get(config.APP_KEY, ".-.-.-.-.-.-.")

	if len(key) == 0 {
		msg := "Application key is not defined. Please define it at APP_KEY in .env file"
		app.Log.Error(msg)
		panic(msg)
	}
}

func (app *Application) initializeCache(checkConnection bool) {
	cache.Bootstrap()

	if !checkConnection {
		return
	}

	app.Log.Info("Checking connection to the Redis cache...")
	err := cache.TestConnection()

	if err != nil {
		msg := fmt.Sprintf("Failed to connect to the Redis: %s. Application will be terminated.", err.Error())
		app.Log.Error(msg)
		panic(msg)
	}
}

// Bootstrap sessions vault
func (app *Application) initializeSessions() {
	//sessionsPath := config.Get(config.TMP_PATH, "/tmp")

	// Check if sessions dir exists and is writable
	/*if unix.Access(sessionsPath, unix.W_OK) != nil {
		msg := fmt.Sprintf(
			"Failed to initialize sessions storage, directory doesn't exists or not writable: %s",
			sessionsPath,
		)

		app.Log.Error(msg)
		panic(msg)
	}*/

	// Get session encrypt key

	// Check if key is empty

	// Bootstrap sessions vault
	//vault.Bootstrap(sessionsPath, sessKey)
}

// Check connection to the main MySQL database
func (app *Application) testSQLConnection() {
	app.Log.Info("Checking connection to the MySQL database...")

	result, err := database.TestConnection()

	if !result {
		msg := fmt.Sprintf("Failed to connect to the database: %s. Application will be terminated.", err.Error())
		app.Log.Error(msg)
		panic(msg)
	}

}
