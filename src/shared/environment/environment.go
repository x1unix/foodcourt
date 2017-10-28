package environment

import (
	"runtime"
	"strings"
	"path"
)

const DIR_PUBLIC = "public"



// Get the application root path
func GetRoot() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}

	return path.Dir(filename)
}

// Get path of file(s) in 'public' directory.
// If subPath is undefined, will return path to 'public' directory
func GetPublic(subPath ...string) string {
	publicPath := GetRoot() + "/" + DIR_PUBLIC

	if (len(subPath) > 0) {
		publicPath += strings.Join(subPath, "/")
	}

	return publicPath
}