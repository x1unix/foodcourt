package config

import (
	"os"
	"fmt"
	"../environment"
	"github.com/joho/godotenv"
)

const ENV_NAME = ".env"

// Dotenv file location
var fileDir = fmt.Sprintf("%s/%s", environment.GetRoot(), ENV_NAME)

// Load configuration from dotenv file
func Bootstrap() {
	err := godotenv.Load()
	if err != nil {
		panic(fmt.Sprintf("Environment file not exists (%s).", fileDir))
	}
}

// Get param value
func Get(key string, otherwize string) string {
	val := os.Getenv(key)

	if (len(val) == 0) {
		return otherwize
	}

	return val
}
