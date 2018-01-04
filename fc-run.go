package main

import (
	"./worker/cmd"
	"./shared/logger"
	"./shared/config"
	"./shared/cache"
	"./worker/sender"
	"./worker/control"
)

func main() {
	config.Bootstrap(true)
	cache.Bootstrap()
	logger.Bootstrap("fc-run")
	cmd.Bootstrap()

	cmd.HandleFunc("orders:send", "sends order report to all users", sender.SendLunchOrders)

	cmd.HandleFunc("orders:report", "sends list of all ordered dishes to the food provider", sender.SendOrderReport)

	cmd.HandleFunc("menu:lockNextDay", "blocks a menu for the next day for edit (read-only)", control.LockMenu)

	cmd.Run()
}