package imageUtil

import (
	"net/http"
	"path/filepath"
	"strings"
)

func GetImageContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".bmp":
		return "image/bmp"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	case ".ico":
		return "image/x-icon"
	case ".tiff", ".tif":
		return "image/tiff"
	default:
		return "image/png"
	}
}

func DetectImageFormat(data []byte) string {
	contentType := http.DetectContentType(data)

	switch {
	case strings.Contains(contentType, "jpeg"):
		return "image/jpeg"
	case strings.Contains(contentType, "png"):
		return "image/png"
	case strings.Contains(contentType, "gif"):
		return "image/gif"
	case strings.Contains(contentType, "bmp"):
		return "image/bmp"
	case strings.Contains(contentType, "webp"):
		return "image/webp"
	default:
		return contentType
	}
}
