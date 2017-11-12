package controller

import (
	"net/http"
	"path/filepath"
	"io"
	"os"
	"encoding/hex"
	"crypto/md5"
	"../shared/environment"
	"../shared/rest"
	"../shared/logger"
	"mime/multipart"
	"gopkg.in/h2non/filetype.v1"
	"fmt"
)

const UploadDir = "photos"
const ParamName = "image"


// Upload a new image
// (POST - /api/media)
func UploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile(ParamName)

	// try to read input file
	if err != nil {
		rest.Error(err).Write(&w)
		return
	}

	// Close the file at the end
	defer file.Close()

	// try to check MIME
	isImage, ext := isImage(file)

	if !isImage {
		rest.ErrorFromString("Provided file is not an image", 400).Write(&w)
		return
	}

	// Get file local location and URL
	local, url, ferr := getFileLocation(&file, ext)

	if ferr != nil {
		logger.GetLogger().Error(fmt.Sprintf("Failed to determine uploaded file destination: %s", local, ferr.Error()))
		rest.Error(ferr).Write(&w)
		return
	}

	// Check if image exists
	imageExists := imageExists(local)

	if imageExists {
		// If image exists - just return current location
		rest.Echo(url).Write(&w)
		return
	}

	// Ensure that upload dir exists
	checkUploadEnv()

	// Open output stream for the image
	f, openErr := os.OpenFile(local, os.O_WRONLY|os.O_CREATE, 0666)

	// Report error if failed to create a file
	if openErr != nil {
		logger.GetLogger().Error(fmt.Sprintf("Failed to prepare file to save '%s': %s", local, openErr.Error()))
		rest.Error(openErr).Write(&w)
		return
	}

	defer f.Close()

	// Try to save the file
	_, copyErr := io.Copy(f, file)

	if copyErr != nil {
		logger.GetLogger().Error(fmt.Sprintf("Failed to save uploaded file '%s': %s", local, copyErr.Error()))
		rest.Error(openErr).Write(&w)
		return
	}

	// Send new file URL if it's uploaded
	rest.Echo(url).Write(&w)
}

// Get file location
func getFileLocation(file *multipart.File, ext string) (string, string, error) {
	hash, err := md5FromFile(file)
	var local string
	var url string

	if err == nil {
		fname := fmt.Sprintf("%s.%s", hash, ext)
		url = "/" + UploadDir + "/" + fname
		local = filepath.ToSlash("./" + environment.DIR_PUBLIC + url)
	}

	return local, url, err
}

// Check if provided file is image
func isImage(file multipart.File) (bool, string) {
	buff := make([]byte, 512)
	_, err := file.Read(buff)

	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("Failed to check file MIME: %s", err.Error()))
		return false, "";
	}

	isImage := filetype.IsImage(buff)
	var ext string

	if isImage {
		kind, _ := filetype.Match(buff)

		ext = kind.Extension
	}

	return isImage, ext
}


// Generate MD5 from file
func md5FromFile(file *multipart.File) (string, error) {
	var returnMD5String string

	hash := md5.New()

	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}

	//Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:16]

	//Convert the bytes to a string
	returnMD5String = hex.EncodeToString(hashInBytes)

	return returnMD5String, nil
}

func imageExists(fileName string) bool {
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

func checkUploadEnv() {
	dirUploads := filepath.ToSlash("./" + environment.DIR_PUBLIC + "/" + UploadDir)

	if _, err := os.Stat(dirUploads); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(dirUploads, 655)
		}
	}
}
