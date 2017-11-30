package config

import (
	"fmt"
	"os"

	"../environment"
	"github.com/joho/godotenv"
)

const ENV_NAME = ".env"

// Dotenv file location
var fileDir = fmt.Sprintf("%s/%s", environment.GetRoot(), ENV_NAME)

// Load configuration from dotenv file
func Bootstrap(quiet bool) {
	err := godotenv.Load()
	if err != nil {
		errMsg := fmt.Sprintf("Environment file does not exist (%s).", fileDir)

		if quiet {
			panic(errMsg)
		} else {
			fmt.Println("config:Bootstrap/WARN: " + errMsg)
		}

	}
}

// Get param value
func Get(key string, otherwise string) string {
	val := os.Getenv(key)

	if len(val) == 0 {
		return otherwise
	}

	return val
}
