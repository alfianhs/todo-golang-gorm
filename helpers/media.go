package helpers

import "fmt"

func GetExtensionFromMimeType(mimeType string) (string, error) {
	switch mimeType {
	case "image/jpeg":
		return "jpg", nil
	case "image/png":
		return "png", nil
	case "image/gif":
		return "gif", nil
	case "image/webp":
		return "webp", nil
	case "application/pdf":
		return "pdf", nil
	default:
		return "", fmt.Errorf("unsupported mime type: %s", mimeType)
	}
}

var allowedMimeTypes = map[string][]string{
	"image":    {"image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp"},
	"document": {"application/pdf"},
}

func IsMimeTypeAllowed(mimeType string, category string) bool {
	allowedTypes, exists := allowedMimeTypes[category]
	if !exists {
		return false
	}
	for _, t := range allowedTypes {
		if t == mimeType {
			return true
		}
	}
	return false
}
