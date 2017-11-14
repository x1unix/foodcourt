package helpers

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
