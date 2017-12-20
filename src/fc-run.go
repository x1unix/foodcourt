package main

import (
	"./worker/cmd"
	"./shared/logger"
	"./shared/config"
	"./worker/sender"
)

func main() {
	config.Bootstrap(true)
	logger.Bootstrap("fc-worker")
	cmd.Bootstrap()

	cmd.HandleFunc("sendLunchOrders", "Some test cmd", sender.SendLunchOrders)

	cmd.Run()
}