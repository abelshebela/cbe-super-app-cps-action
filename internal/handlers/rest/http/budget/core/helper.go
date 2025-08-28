package core

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func ParseMultipartFormFile(r *http.Request, key string, maxMemory int64) (multipart.File, *multipart.FileHeader, error) {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		return nil, nil, errors.New(localization.ErrorInvalidRequest.Code)
	}
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		return nil, nil, errors.New(localization.ErrorBankFileParseFailed.Code)
	}

	file, fileHeader, err := r.FormFile(key)
	if err != nil {
		if err == http.ErrMissingFile {
			return nil, nil, errors.New(localization.ErrorMissingFile.Code)
		}
		return nil, nil, fmt.Errorf("missing or invalid file for key '%s': %w", key, err)
	}

	return file, fileHeader, nil
}

func FileValidator(w http.ResponseWriter, file multipart.File, fileHeader multipart.FileHeader, logger utils.Logger) error {
	defer file.Close()

	allowedImageTypes := map[string]bool{
		"image/jpeg":    true,
		"image/png":     true,
		"image/gif":     true,
		"image/webp":    true,
		"image/bmp":     true,
		"image/svg+xml": true,
	}

	if !allowedImageTypes[fileHeader.Header.Get("Content-Type")] {
		logger.Errorf("invalid file type: %v", fileHeader.Header.Get("Content-Type"))
		localization.SendErrorResponse(w, localization.ErrorInvalidFormat, nil, nil)
		return errors.New(localization.ErrorInvalidFormat.Code)
	}
	return nil
}
