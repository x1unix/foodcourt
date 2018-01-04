package vault

import (
	"crypto/rand"
	"encoding/base64"
)

// Generate new bytes sequence
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)

	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Generate new token
func NewToken() (string, error) {
	bytes, err := generateRandomBytes(SessionTokenLength)

	return base64.URLEncoding.EncodeToString(bytes), err
}
