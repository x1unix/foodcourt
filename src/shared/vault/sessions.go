package vault

import (
	"github.com/gorilla/sessions"
	"net/http"
	"../logger"
	"fmt"
	"time"
	"errors"
	"github.com/gorilla/context"
)

// Token length
const SessionTokenLength = 32

const ContextTokenKey = "token"
const KeyLifetime = "TTL"
const KeyUserId = "userId"
const KeyAuthorized = "Authorized"
const DefaultTTL = 1209600 // 2 weeks

// Global session store
var Store *sessions.FilesystemStore

// Bootstrap sessions vault
func Bootstrap(sessionsLocation string, encryptKey string) {
	Store = sessions.NewFilesystemStore(sessionsLocation, []byte(encryptKey))
}

// Get raw value from the session
func GetValue(session *sessions.Session, key string) (interface{}, error) {
	val := session.Values[key]

	if (val == nil) {
		return nil, errors.New("Session key not found: " + key)
	}

	return val, nil
}

// Get string value from the session
func GetString(session *sessions.Session, key string) (string, error) {
	raw, err := GetValue(session, key)

	if (err != nil) {
		return "", err
	}

	val, ok := raw.(string)

	if (!ok) {
		return "", errors.New("Failed to extract value")
	}

	return val, err
}

// Get int64 value from the session
func GetInt64(session *sessions.Session, key string) (int64, error) {
	raw, err := GetValue(session, key)

	if (err != nil) {
		return 0, err
	}

	val, ok := raw.(int64)

	if (!ok) {
		return 0, errors.New("Failed to extract value")
	}

	return val, err
}

// Get boolean value from the session
func GetBool(session *sessions.Session, key string) (bool, error) {
	raw, err := GetValue(session, key)

	if (err != nil) {
		return false, err
	}

	val, ok := raw.(bool)

	if (!ok) {
		return false, errors.New("Failed to extract value")
	}

	return val, err
}

// Checks if the session is valid
func IsSessionValid(r *http.Request, token string) (*sessions.Session, bool) {
	session, err := Store.Get(r, token)

	// check session extract errors
	if err != nil {
		logger.GetLogger().Warning(fmt.Sprintf("Failed to extract session: %s", err.Error()))
		return session, false
	}

	// extract ttl
	ttl, err := GetInt64(session, KeyLifetime)

	if (err != nil) {
		logger.GetLogger().Warning(fmt.Sprintf("Failed to extract session TTL: %s", err.Error()))
		return session, false
	}


	// Check ttl
	now := time.Now().Unix()

	if ttl >= now {
		return session, false
	}

	return session, true
}

// Build session data object from session
func BuildSessionData(session *sessions.Session) (*SessionData, error) {
	ttl, ttlErr := GetInt64(session, KeyLifetime)
	authStat, authStatErr := GetBool(session, KeyAuthorized)

	if ttlErr != nil || authStatErr != nil {
		return nil, errors.New("Failed to extract TTL or AuthStatus from the session")
	}

	userId, ok := 0, false

    rawUser := session.Values[KeyUserId]

	if userId, ok = rawUser.(int); !ok {
		return nil, errors.New("Failed to extract user data from the session")
	}

	data := SessionData{
		ttl,
		authStat,
		userId,
	}

	return &data, nil
}

// Register new session and get session data
func NewSessionTicket(r *http.Request, w *http.ResponseWriter, userId int) (*SessionData, error) {
	token := context.Get(r, ContextTokenKey).(string)
	ttl := time.Now().Unix() + DefaultTTL

	session, err := Store.Get(r, token)

	if (err == nil) {
		session.Values[KeyLifetime] = ttl
		session.Values[KeyAuthorized] = true
		session.Values[KeyUserId] = userId

		data := SessionData {
			ttl,
			true,
			userId,
		}

		session.Save(r, *w)
		return &data, err
	}

	return nil, err
}

