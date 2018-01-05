package cache

import (
	"foodcourt/config"
	"github.com/go-redis/redis"
)

const DefaultHost = "localhost"
const DefaultPort = "6379"
const DefaultPass = ""
const DefaultDB = 0

var Client *redis.Client

// Bootstrap cache
func Bootstrap() {
	Client = redis.NewClient(getConnectionOptions())
}

// Test connection to the Redis cache
func TestConnection() error {
	_, err := Client.Ping().Result()
	return err
}

// Get redis cache connection options
func getConnectionOptions() *redis.Options {
	host := config.Get(config.REDIS_HOST, DefaultHost)
	port := config.Get(config.REDIS_PORT, DefaultPort)
	pass := config.Get(config.REDIS_PASSWORD, DefaultPass)

	return &redis.Options{
		Addr:     host + ":" + port,
		Password: pass,
		DB:       DefaultDB,
	}
}
