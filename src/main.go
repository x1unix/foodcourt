package main

import (
	"./shared/logger"
	"./shared/config"
	"./shared/app"
	"./route"
)

func main() {
	config.Bootstrap()

	log := logger.Bootstrap()
	router := route.Bootstrap()

	httpPort := config.Get(config.HTTP_PORT, "80")
	httpHost := config.Get(config.HTTP_HOST, "0.0.0.0")

	log.Info("Starting application...")

	appInstance := app.Application {
		router,
		log,
	}

	appInstance.Run(httpHost, httpPort)
}
