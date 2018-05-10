package auth

import (
	"time"
	"foodcourt/cache"
	"github.com/go-redis/redis"
)

const recoveryTokenMinutesTTL = 15
const resetTokenKeyPrefix = "resettoken__"

type ResetTokenPair struct {
	CoreKey string
	dynamicKey string
}

func (t *ResetTokenPair) GetRedisKey() string {
	return "resettoken__" + t.dynamicKey + ":" + t.CoreKey
}

func (t *ResetTokenPair) Exists() (bool, error) {
	val, err := cache.Client.Exists(t.GetRedisKey()).Result()

	if err != nil {
		return false, err
	}

	return (val > 0), nil
}

func (t *ResetTokenPair) Reveal() {
	cache.Client.Del(t.GetRedisKey()).Result()
}

// GetTokenOwner gets token's owner email
func (t *ResetTokenPair) GetTokenOwner() (string, error) {
	email, err := cache.Client.Get(t.GetRedisKey()).Result()

	if err == redis.Nil {
		return "Invalid token", nil
	} else {
		return "", err
	}

	return email, nil
}

func (t *ResetTokenPair) SetDynamicKey(userAgent string, ip string) *ResetTokenPair {
	t.dynamicKey = EncryptString(userAgent + ip)
	return t
}

func (t *ResetTokenPair) Save(ownerEmail string) error {
	return cache.Client.Set(t.GetRedisKey(), ownerEmail, time.Duration(time.Minute * recoveryTokenMinutesTTL)).Err()
}

func NewTokenPair(coreKey string, userAgent string, ip string) *ResetTokenPair {
	pair := ResetTokenPair{}
	pair.CoreKey = coreKey
	pair.SetDynamicKey(userAgent, ip)

	return &pair
}

func BuildResetTokenPair(email string, code string, agent string, ip string) *ResetTokenPair {
	return NewTokenPair(EncryptString(email + code), agent, ip)
}
