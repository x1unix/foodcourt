package logger

import (
	"os"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("voracity")

var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{module}:%{shortfunc}/%{level:.4s}: %{color:reset} %{message}`,
)

// Initialize the application logger
func Bootstrap() {

	consoleBackend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(consoleBackend, format)

	backendLeveled := logging.AddModuleLevel(consoleBackend)
	backendLeveled.SetLevel(logging.ERROR, "")

	logging.SetBackend(consoleBackend, backendFormatter)
}
