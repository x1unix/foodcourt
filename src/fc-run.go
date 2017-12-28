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

	cmd.HandleFunc("orders:send", "Some test cmd", sender.SendLunchOrders)

	cmd.Run()
}