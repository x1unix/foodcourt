package rest

import (
	"io"
	"encoding/json"
)

func ReadJSONBody(body *io.ReadCloser, target *interface{}) error {
	b := *body
	decoder := json.NewDecoder(b)
	err := decoder.Decode(target)

	defer b.Close()

	return err
}