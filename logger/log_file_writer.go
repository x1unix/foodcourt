package logger

import (
	"os"
)

// Log file writer
type LogFileWriter struct {
	filePath string
}

// Write entry
func (w LogFileWriter) Write(p []byte) (n int, err error) {
	f, err := os.OpenFile(w.filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		defer f.Close()
		return 0, err
	}

	defer f.Close()
	return f.Write(p)
}

func (w LogFileWriter) FileExists() bool {
	// detect if file exists
	var _, err = os.Stat(w.filePath)

	return !os.IsNotExist(err)
}

func (w LogFileWriter) PurgeFile() error {

	if !w.FileExists() {
		return nil
	}

	err := os.Remove(w.filePath)

	return err
}
