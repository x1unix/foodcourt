package storage

import (
	"../cache"
	"../logger"
	"fmt"
	"strconv"
	"github.com/vmihailenco/msgpack"
)

const keyPrefix = "com.foodcourt."

const errChkExists = "Failed to check if value exists on key '%s': %v"
const errGet = "Failed to get value on key '%s': %v"
const errDeserialize = "Failed to deserialize value to '%s' (key: %s) (data: %s): %v"

const magicPropEmpty = "@EMPTY"

func prop(key string) string {
	return keyPrefix + key
}

// Checks if property is defined
func Exists(key string) bool {
	if val, err := cache.Client.Exists(prop(key)).Result(); err != nil {
		logger.GetLogger().Error(fmt.Sprintf(errChkExists, key, err))
		return false;
	} else {
		return val > 0;
	}
}

func GetString(key string, otherwize string) string {
	if raw, err := cache.Client.Get(prop(key)).Result(); err != nil {
		logger.GetLogger().Error(fmt.Sprintf(errGet, key, err))
		return otherwize
	} else {
		return raw
	}
}

func GetBoolean(key string, otherwize bool) bool {
	raw := GetString(key, "false")
	if val, err := strconv.ParseBool(raw); err != nil {
		logger.GetLogger().Error(fmt.Sprintf(errDeserialize, "boolean", key, err))
		return otherwize
	} else {
		return val;
	}
}

func Get(key string, otherwize interface{}) interface{} {
	raw := GetString(key, magicPropEmpty)
	var output interface{}

	if raw == magicPropEmpty {
		return otherwize
	}

	if err := msgpack.Unmarshal([]byte(raw), &output); err != nil {
		logger.GetLogger().Error(fmt.Sprintf(errDeserialize, "interface{}", key, err))
		return otherwize
	} else {
		return output
	}
}

func SetString(key string, val string) error {
	return cache.Client.Set(prop(key), val, 0).Err()
}

func SetBoolean(key string, val bool) error {
	return SetString(key, strconv.FormatBool(val))
}

func Set(key string, val interface{}) error {
	if data, err := msgpack.Marshal(val); err != nil {
		return err
	} else {
		return cache.Client.Set(prop(key), data, 0).Err()
	}
}