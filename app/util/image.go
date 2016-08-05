package util

// ValidFormat checks if an image format is valid
func ValidFormat(format string) bool {
	switch format {
	case "gif", "png", "jpeg":
		return true
	}
	return false
}
