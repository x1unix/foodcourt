package vault

import (
	"time"
)

type SessionData struct {
	TTL int64
	Authorized bool
	UserId int
}

func (s *SessionData) IsValid() bool {
	now := time.Now().Unix()

	return s.TTL < now
}
