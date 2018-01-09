package environment

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const DIR_PUBLIC = "public"
const ResourcesDir = "/resources"

// Get the application root path
func GetRoot() string {

	exePath, err := os.Executable()

	if err != nil {
		panic(fmt.Sprintf("Failed to get executable path: %v", err))
	}

	return filepath.Dir(exePath)
}

// Get path of file(s) in 'public' directory.
// If subPath is undefined, will return path to 'public' directory
func GetPublic(subPath ...string) string {
	publicPath := GetRoot() + "/" + DIR_PUBLIC

	if len(subPath) > 0 {
		publicPath += strings.Join(subPath, "/")
	}

	return publicPath
}

func GetResourcePath(path string) string {
	return GetRoot() + "/" + ResourcesDir + "/" + path
}