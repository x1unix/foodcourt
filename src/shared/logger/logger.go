package logger

import (
	"os"
	"path/filepath"
	"../config"
	"github.com/op/go-logging"
	"fmt"
)

// Log file name
const LOG_FILE_NAME = "voracity.log"
const LOG_NAME = "voracity"

// Console output formatter
var consoleFormatter = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{module}:%{shortfunc}/%{level:.4s}: %{color:reset} %{message}`,
)

// File output
var fileFormatter = logging.MustStringFormatter(
	`%{time:15:04:05.000} %{module}:%{shortfunc}/%{level:.4s}: %{message}`,
)

// Gets location of logs directory
func GetLogPath() string {
	// Get log path from the config
	logPathOrigin := config.Get("LOG_PATH", ".")

	// Get absolute log path
	logPath, _ := filepath.Abs(logPathOrigin)

	return logPath
}

// Gets full log file name and path
func GetLogFilePath() string {
	logPath := GetLogPath()

	return fmt.Sprintf("%s/%s", logPath, LOG_FILE_NAME)
}

// Initialize the application logger
func Bootstrap() *logging.Logger {

	// Logger instance
	var log = logging.MustGetLogger(LOG_NAME)

	// Get log file name and path
	logFilePath := GetLogFilePath()

	// Log file writer instance
	logFileWriter := LogFileWriter{
		logFilePath,
	}

	// Purge log file if it exists
	logFileWriter.PurgeFile()

	// Log file writer backend
	fileBackend := logging.NewLogBackend(logFileWriter, "", 0)

	// Console output writer backend
	consoleBackend := logging.NewLogBackend(os.Stderr, "", 0)

	// For messages written to console we want to add some additional
	// information to the output, including the used log level and the name of
	// the function.
	consoleBackendFormatter := logging.NewBackendFormatter(consoleBackend, consoleFormatter)
	fileBackendFormatter := logging.NewBackendFormatter(fileBackend, fileFormatter)

	// Only errors and more severe messages should be sent to log file
	fileBackendLeveled := logging.AddModuleLevel(fileBackendFormatter)
	fileBackendLeveled.SetLevel(logging.ERROR, "")

	// Set the backends to be used.
	logging.SetBackend(fileBackendLeveled, consoleBackendFormatter)

	log.Info("info")
	log.Notice("notice")
	log.Warning("warning")
	log.Error("err")
	log.Critical("crit")

	return log
}

func GetLogger() *logging.Logger {
	return logging.MustGetLogger("voracity")
}
