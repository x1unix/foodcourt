package main

import (
	"./shared/logger"
	"./shared/config"
	"./shared/app"
	"./route"
	"fmt"
)

func main() {
	config.Bootstrap()

	log := logger.Bootstrap()
	router := route.Bootstrap()

	httpPort := config.Get(config.HTTP_PORT, "80")
	httpHost := config.Get(config.HTTP_HOST, "0.0.0.0")

	log.Info(fmt.Sprintf("HTTP server started at %s:%s", httpHost, httpPort))

	appInstance := app.Application {
		router,
		nil,
		log,
	}

	appInstance.Run(httpHost, httpPort)
}
