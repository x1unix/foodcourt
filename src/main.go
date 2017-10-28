package main

import (
	"./shared/logger"
	"./shared/config"
	"./route"
	"net/http"
	"fmt"
	"time"
)

func main() {
	config.Bootstrap()

	log := logger.Bootstrap()
	router := route.Bootstrap()

	httpPort := config.Get("HTTP_PORT", "80")
	httpHost := config.Get("HTTP_HOST", "0.0.0.0")
	httpAddr := httpHost + ":" + httpPort

	log.Info(fmt.Sprintf("Starting HTTP server at %s", httpAddr))

	server := &http.Server{
		Handler: router,
		Addr: httpAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
