package logger

import (
	"fmt"
	"foodcourt/config"
	"foodcourt/environment"
	"os"
	"path/filepath"
	"strings"
	"github.com/op/go-logging"
)

// Log file name
var LogFileName = "fc.log"
var LogName = "fc"

// Console output formatter
var consoleFormatter = logging.MustStringFormatter(
	`%{color}[%{time:2006-01-02 15:04}] %{module}:%{shortfunc}/%{level:.6s}: %{color:reset} %{message}`,
)

// File output
var fileFormatter = logging.MustStringFormatter(
	`[%{time:2006-01-02 15:04}] %{module}:%{shortfunc}/%{level:.6s}: %{message}`,
)

var LogLevels = []string{"critical", "error", "warning", "notice", "info", "debug"}

// Gets log level from the configuration
func GetLogLevel() logging.Level {
	logLevel := strings.ToLower(config.Get(config.LOG_LEVEL, "ALL"))

	switch logLevel {
	case "info":
		return logging.INFO
	case "notice":
		return logging.NOTICE
	case "warning":
		return logging.WARNING
	case "error":
		return logging.ERROR
	case "critical":
		return logging.CRITICAL
	default:
		return logging.DEBUG
	}
}

// Gets location of logs directory
func GetLogPath() string {
	// Get log path from the config
	logPathOrigin := config.Get(config.LOG_PATH, environment.GetRoot())

	// Get absolute log path
	logPath, _ := filepath.Abs(logPathOrigin)

	return logPath
}

// Gets full log file name and path
func GetLogFilePath() string {
	logPath := GetLogPath()

	return fmt.Sprintf("%s/%s", logPath, LogFileName)
}

// Initialize the application logger
func Bootstrap(logName string) *logging.Logger {

	LogName = logName
	LogFileName = logName + ".log"

	// Logger instance
	var log = logging.MustGetLogger(logName)

	// Get log file name and path
	logFilePath := GetLogFilePath()

	// Log file writer instance
	logFileWriter := LogFileWriter{
		logFilePath,
	}

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

	logLevel := GetLogLevel()
	fileBackendLeveled.SetLevel(logLevel, "")

	// Set the backends to be used.
	logging.SetBackend(fileBackendLeveled, consoleBackendFormatter)

	return log
}

func GetLogger() *logging.Logger {
	return logging.MustGetLogger(LogName)
}
