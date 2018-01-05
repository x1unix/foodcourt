package controller

import (
	"foodcourt/environment"
	"foodcourt/logger"
	"foodcourt/rest"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const UploadDir = "photos"
const ParamName = "image"

// Upload a new photo
// (POST - /api/media)
func UploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile(ParamName)

	// try to read input file
	if err != nil {
		rest.Error(err).Write(&w)
		return
	}

	defer file.Close()

	// Check image MIME type
	isImage, ext := isImageType(handler)

	if !isImage {
		rest.ErrorFromString("Provided file is not an image", 400).Write(&w)
		return
	}

	// Generate random file name
	fileName := generateFileName(handler.Filename) + ext

	// Define file URL and path
	furl := "/photos/" + fileName
	flocal := "./public" + furl

	// Ensure that upload dir exists
	checkUploadEnv()

	f, err := os.OpenFile(flocal, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("Failed to create upload file '%s': %s", flocal, err.Error()))
		rest.Error(err).Write(&w)
		return
	}

	defer f.Close()

	_, copyErr := io.Copy(f, file)

	if copyErr != nil {
		logger.GetLogger().Error(fmt.Sprintf("Failed to save uploaded file '%s': %s", flocal, copyErr.Error()))
		rest.Error(copyErr).Write(&w)
		return
	}

	// Send new file URL if it's uploaded
	rest.Echo(furl).Write(&w)

	return
}

// Generate hashed file name
func generateFileName(originalFileName string) string {
	// Current date
	now := strconv.FormatInt(time.Now().Unix(), 10)

	// Random number based on file name
	rndNum := strconv.Itoa(rand.Intn(len(originalFileName)))

	// file hash salt
	salt := []byte(originalFileName + ":" + now + ":" + rndNum)

	hasher := md5.New()

	// hash darling, hash
	hasher.Write(salt)
	return hex.EncodeToString(hasher.Sum(nil))
}

// Assert image type
func isImageType(handler *multipart.FileHeader) (bool, string) {
	mime := handler.Header.Get("Content-Type")
	isImageMime := strings.Contains(mime, "image/")
	return isImageMime, filepath.Ext(handler.Filename)
}

// Check if upload directory exists
func checkUploadEnv() {
	dirUploads := filepath.ToSlash("./" + environment.DIR_PUBLIC + "/" + UploadDir)

	if _, err := os.Stat(dirUploads); err != nil {
		if os.IsNotExist(err) {
			logger.GetLogger().Info("Images upload directory doesn't exists. A new one will be created.")
			logger.GetLogger().Debug(fmt.Sprintf("Images directory: %s", dirUploads))
			os.MkdirAll(dirUploads, 0700)
		}
	}
}
