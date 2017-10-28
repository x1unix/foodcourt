package main

import (
	"./shared/logger"
	"./shared/config"
	"github.com/op/go-logging"
)

func main() {
	logger.Bootstrap()
	config.Bootstrap()

	var log = logging.MustGetLogger("voracity")
	log.Info(config.Get("HTTP_PORT", "80"))
	log.Info("Foo Bar");
}
