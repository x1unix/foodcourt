package vault

import (
	"github.com/gorilla/sessions"
)

// Token length
const SessionTokenLength = 32

// Global session store
var Store *sessions.FilesystemStore

func Bootstrap(sessionsLocation string, encryptKey string) {
	Store = sessions.NewFilesystemStore(sessionsLocation, []byte(encryptKey))
}


