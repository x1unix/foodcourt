package settings

import (
	"../storage"
	"github.com/vmihailenco/msgpack"
)

type Settings struct {
	BaseURL string 			`msgpack:"baseUrl" json:"baseUrl" validate:"required,url"`
	SMTP 	SMTPSettings 	`msgpack:"smtp" json:"smtp" validate:"required"`
	Sender  SenderSettings	`msgpack:"sender" json:"sender" validate:"required"`
}

const StorageKey = "settings"

func GetSettings() (*Settings, error) {
	var settings Settings
	data, err := storage.Get(StorageKey)

	if err != nil {
		return &settings, err
	}

	err = msgpack.Unmarshal([]byte(data), &settings);

	return &settings, err
}

func SetSettings(data *Settings) error {
	settings := &data
	if data, err := msgpack.Marshal(settings); err != nil {
		return err
	} else {
		return storage.SetRaw(StorageKey, data)
	}
}
