package storage

import (
	"../cache"
	"../logger"
	"fmt"
	"strconv"
	"github.com/vmihailenco/msgpack"
	"github.com/go-redis/redis"
	"errors"
)

const keyPrefix = "com.foodcourt."

const errChkExists = "Failed to check if value exists on key '%s': %v"
const errGet = "Failed to get value on key '%s': %v"
const errDeserialize = "Failed to deserialize value to '%s' (key: %s) (data: %s): %v"
const warnNoKey = "Requested key value doesn't exists ('%s'), a default value provided"

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

func GetString(key string, otherwise string) string {
	if raw, err := cache.Client.Get(prop(key)).Result(); err != nil {
		log := logger.GetLogger()
		if err != redis.Nil {
			log.Warning(fmt.Sprintf(warnNoKey, key))
		} else {
			log.Error(fmt.Sprintf(errGet, key, err))
		}
		return otherwise
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

func Get(key string) (string, error) {
	rdKey := prop(key)
	data := ""
	if raw, err := cache.Client.Get(rdKey).Result(); err != nil {
		if err != redis.Nil {
			return data, err
		} else {
			return data, errors.New(fmt.Sprintf("Requested key value not exists: '%s'", rdKey))
		}
	} else {
		return raw, nil
	}
}

func SetString(key string, val string) error {
	return cache.Client.Set(prop(key), val, 0).Err()
}

func SetBoolean(key string, val bool) error {
	return SetString(key, strconv.FormatBool(val))
}

func SetRaw(key string, val []byte) error {
	return cache.Client.Set(prop(key), val, 0).Err()
}

func Set(key string, val interface{}) error {
	if data, err := msgpack.Marshal(val); err != nil {
		return err
	} else {
		return cache.Client.Set(prop(key), data, 0).Err()
	}
}