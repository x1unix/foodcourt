package app

import (
	"github.com/gorilla/mux"
	"github.com/op/go-logging"
	"time"
	"net/http"
	"fmt"
)

import _ "github.com/go-sql-driver/mysql"

type Application struct {
	Router *mux.Router
	Log *logging.Logger
}

func (app *Application) Run(httpHost string, httpPost string) {
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
