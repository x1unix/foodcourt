package vault

import (
	"../auth"
	"../cache"
	"github.com/go-redis/redis"
	"github.com/vmihailenco/msgpack"
	"time"
)

// Token length
const SessionTokenLength = 32

const ContextTokenKey = "token"
const KeyLifetime = "TTL"
const KeyUserId = "userId"
const KeyAuthorized = "Authorized"
const DefaultTTL = 14 // days

const CacheKeyPrefix = "sess_"

func ssid(token string) string {
	return CacheKeyPrefix + token
}

// Register new session and get session data
func NewSession(token string, user *auth.User) (*Session, error) {
	ttlDr := time.Duration(14 * (time.Hour * 24))
	ttl := time.Now().AddDate(0, 0, DefaultTTL).Unix()

	// Create session
	session := &Session{
		TTL:        ttl,
		Token:      token,
		Authorized: true,
		UserId:     user.ID,
		User:       user,
	}

	// Serialize data
	msgp, err := msgpack.Marshal(session)

	if err != nil {
		return nil, err
	}

	sessionKey := ssid(token)

	// Save session
	err = cache.Client.Set(sessionKey, msgp, ttlDr).Err()

	return session, err
}

// Check session exists
func SessionExists(token string) (bool, error) {
	ssid := ssid(token)
	val, err := cache.Client.Exists(ssid).Result()

	exists := val > 0

	return exists, err
}

// Get session by token id
func GetSession(token string) (*Session, error) {
	ssid := ssid(token)
	raw, err := cache.Client.Get(ssid).Result()

	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		session := &Session{}
		serr := msgpack.Unmarshal([]byte(raw), session)

		return session, serr
	}
}

// Delete session and reveal token
func RevealToken(token string) {
	ssid := ssid(token)
	cache.Client.Del(ssid).Result()
}
