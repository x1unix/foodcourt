package rest

import (
	"encoding/json"
	"io"
	"net/http"
	"net"
)

const dateParamLength = 8

func ReadJSONBody(body *io.ReadCloser, target *interface{}) error {
	b := *body
	decoder := json.NewDecoder(b)
	err := decoder.Decode(target)

	defer b.Close()

	return err
}

// Checks if date format valid
func DateFormatValid(dateString string) bool {
	return len(dateString) == dateParamLength;
}

func GetClientIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		return r.RemoteAddr
	}

	return ip
}

