package utils

import "strings"

func GetContentTypeByFileName(filename string) (string, string) {
	switch {
	case strings.Contains(filename, ".png"):
		return "image/png", ".png"
	case strings.Contains(filename, ".jpg"):
		return "image/jpeg", ".jpg"
	case strings.Contains(filename, ".jpeg"):
		return "image/jpeg", ".jpeg"
	case strings.Contains(filename, ".gif"):
		return "image/gif", ".gif"
	case strings.Contains(filename, ".webp"):
		return "image/webp", ".webp"
	}

	return "", ""
}
