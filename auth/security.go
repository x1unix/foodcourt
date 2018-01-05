package auth

import (
	"foodcourt/config"
	"encoding/hex"
	"golang.org/x/crypto/scrypt"
)

func Password(pass string) string {
	key := config.Get(config.APP_KEY, "")
	salt := []byte(key)

	// hash darling, hash
	cryptedBytes, _ := scrypt.Key([]byte(pass), salt, 16384, 8, 1, 32)
	hexPass := hex.EncodeToString(cryptedBytes)

	return hexPass
}
