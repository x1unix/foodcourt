package helpers

import (
	"strconv"
	"strings"
)

// Check if value is not empty
func NotEmpty(item interface{}) bool {
	switch item.(type) {
	case int:
		return item.(int) > 0

	case string:
		return len(item.(string)) > 0
	default:
		return true
	}
}

func EscapeString(origin string) string {
	out := strings.Replace(origin, "'", "\\'", -1)
	return strconv.Quote(out)
}
