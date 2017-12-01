package main

import (
	"./route"
	"./shared/app"
	"./shared/config"
	"./shared/logger"
	"./shared/manifest"
	"fmt"
)

// Application version (defined at build args)
var version string = "1.0.0"

// Application build (defined at build args)
var buildId string = "unknown"

func main() {

	fmt.Println(fmt.Sprintf("============== FoodCourt [Version %s (build %s)] ==============", version, buildId))

	manifest.SetApplicationInfo(version, buildId)

	config.Bootstrap(true)

	log := logger.Bootstrap("fc")
	router := route.Bootstrap()

	httpPort := config.Get(config.HTTP_PORT, "80")
	httpHost := config.Get(config.HTTP_HOST, "0.0.0.0")

	log.Info("Starting application...")

	appInstance := app.Application{
		router,
		log,
	}

	appInstance.Run(httpHost, httpPort)
}
