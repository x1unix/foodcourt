package main

import (
	"./worker/cmd"
	"./shared/logger"
	"./shared/config"
	"./worker/sender"
)

func main() {
	config.Bootstrap(true)
	logger.Bootstrap("fc-run")
	cmd.Bootstrap()

	cmd.HandleFunc("orders:send", "sends order report to all users", sender.SendLunchOrders)

	cmd.HandleFunc("orders:report", "sends list of all ordered dishes to the food provider", sender.SendOrderReport)

	cmd.Run()
}