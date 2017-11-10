package vault

import (
	"time"
	"../auth"
)

type Session struct {
	TTL int64 `msgpack:"ttl" json:"ttl"`
	Token string `msgpack:"token" json:"token"`
	Authorized bool `msgpack:"authorized" json:"authorized"`
	UserId int `msgpack:"userId" json:"userId"`
	User *auth.User `msgpack:"user" json:"user"`
}

func (s *Session) IsValid() bool {
	now := time.Now().Unix()

	return s.TTL < now
}
