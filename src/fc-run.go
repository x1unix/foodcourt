package main

import (
	"./worker/cmd"
	"./shared/logger"
	"./shared/config"
	"fmt"
)

func main() {
	config.Bootstrap(true)
	logger.Bootstrap("fc-worker")
	cmd.Bootstrap()

	cmd.HandleFunc("sendLunchOrders", "Some test cmd", func() (bool, error) {
		fmt.Println("Hello world!")
		return true, nil
	})

	cmd.Run()
}