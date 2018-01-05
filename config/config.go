package config

import (
	"fmt"
	"os"

	"foodcourt/environment"
	"github.com/joho/godotenv"
	"strconv"
	"errors"
)

// EnvFileName file name
const EnvFileName = ".env"

// Dotenv file location
var fileDir = fmt.Sprintf("%s/%s", environment.GetRoot(), EnvFileName)

const errPropNotFound = "Property variable '%s' is not defined. Please define the value in '.env' file or in ENV variables"

// Bootstrap loads configuration from dotenv file
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

// Gets multiple items and returns and error if one of the elements wasn't found
func GetMultiple(output *map[string] string, keys ...string) error {
	dict := *output

	for _, key := range keys {
		if val := os.Getenv(key); len(val) > 0 {
			dict[key] = val
		} else {
			return errors.New(fmt.Sprintf(errPropNotFound, key))
		}
	}

	return nil
}

func GetUnsafe(key string) (error, string) {
	val := os.Getenv(key)

	if len(val) == 0 {
		return errors.New(fmt.Sprintf(errPropNotFound, key)), ""
	}

	return nil, val
}

// Get param value
func Get(key string, otherwise string) string {
	val := os.Getenv(key)

	if len(val) == 0 {
		return otherwise
	}

	return val
}

func GetInt(key string) int {
	val := Get(key, "0")

	if i, e := strconv.Atoi(val); e != nil {
		return 0
	} else {
		return i
	}
}