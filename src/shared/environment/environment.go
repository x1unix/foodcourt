package environment

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
)

const DIR_PUBLIC = "public"

// Get the application root path
func GetRoot() string {
	filename, err := filepath.Abs(".")

	if err != nil {
		panic(fmt.Sprintf("Failed to get document root path: %s", err.Error()))
	}

	return path.Dir(filename)
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
