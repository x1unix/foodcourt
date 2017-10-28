package main

import (
	"./shared/logger"
	"./shared/config"
	"./route"
	"net/http"
	"github.com/op/go-logging"
	"fmt"
)

func main() {
	logger.Bootstrap()
	config.Bootstrap()
	router := route.Bootstrap()

	log := logging.MustGetLogger("voracity")

	httpPort := config.Get("HTTP_PORT", "80")

	log.Info(fmt.Sprintf("Starting HTTP server at port %s", httpPort))

	err := http.ListenAndServe(":" + httpPort, router);

	if err != nil {
		msg := fmt.Sprintf("Failed to start server: %s", err.Error())
		log.Error(msg)
		panic(msg)
	}



}
