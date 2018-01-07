package settings

import (
	"foodcourt/storage"
	"foodcourt/logger"

	"github.com/vmihailenco/msgpack"
)

type Settings struct {
	BaseURL string         `msgpack:"baseUrl" json:"baseUrl" validate:"required,url"`
	SMTP    SMTPSettings   `msgpack:"smtp" json:"smtp" validate:"required"`
	Sender  SenderSettings `msgpack:"sender" json:"sender" validate:"required"`
}

const StorageKey = "settings"

var DefaultSettings = Settings{
	BaseURL: "http://localhost",
	Sender: SenderSettings{
		Enable: false,
	},
}

func GetSettings() (*Settings, error) {
	var settings Settings

	if !storage.Exists(StorageKey) {
		logger.GetLogger().Warning("Settings are not initialized, please configure the app in system settings")
		return &DefaultSettings, nil
	}

	data, err := storage.Get(StorageKey)

	if err != nil {
		return &settings, err
	}

	err = msgpack.Unmarshal([]byte(data), &settings)

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
