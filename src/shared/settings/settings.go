package settings

import (
	"../storage"
)

type Settings struct {
	BaseURL string 			`msgpack:"baseUrl" json:"baseUrl"`
	SMTP 	SMTPSettings 	`msgpack:"smtp" json:"smtp"`
	Sender  SenderSettings	`msgpack:"sender" json:"sender"`
}

const StorageKey = "settings"

func GetSettings() *Settings {
	data := storage.Get(StorageKey, Settings{})
	settings := data.(Settings)
	return &settings
}

func SetSettings(data *Settings) error {
	return storage.Set(StorageKey, *data)
}
