package config

import (
	"../environment"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

const ENV_NAME = ".env"

// Dotenv file location
var fileDir = fmt.Sprintf("%s/%s", environment.GetRoot(), ENV_NAME)

// Load configuration from dotenv file
func Bootstrap() {
	err := godotenv.Load()
	if err != nil {
		panic(fmt.Sprintf("Environment file does not exist (%s).", fileDir))
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
