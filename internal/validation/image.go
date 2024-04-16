package validation

import (
	"errors"
	"mime/multipart"
)

func ValidateImageFile(fileHeader *multipart.FileHeader, maxImageSize int64) error {
	// Allowed image types
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		// Add other allowed types here if needed
	}

	// 1. Check File Type
	if !allowedTypes[fileHeader.Header.Get("Content-Type")] {
		return errors.New("invalid image type. Only JPEG and PNG are supported")
	}

	// 2. Check File Size
	if fileHeader.Size > maxImageSize {
		return errors.New("image exceeds maximum allowed size")
	}

	// If all checks pass, there's no error
	return nil
}
