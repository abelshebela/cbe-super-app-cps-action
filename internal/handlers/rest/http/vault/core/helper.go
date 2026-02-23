package core

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"fmt"
	"mime"
	"mime/multipart"
	"net/http"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func ParseMultipartFormFile(r *http.Request, key string, maxMemory int64, isRequired bool, logger utils.Logger) (multipart.File, *multipart.FileHeader, error) {
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		return nil, nil, fmt.Errorf(localization.MsgMissingContentTypeHeader)
	}

	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil || mediaType != "multipart/form-data" {
		return nil, nil, fmt.Errorf("invalid Content-Type: expected multipart/form-data, got %s", contentType)
	}

	if err := r.ParseMultipartForm(maxMemory); err != nil {
		return nil, nil, fmt.Errorf("failed to parse multipart form: %w", err)
	}

	file, fileHeader, err := r.FormFile(key)
	if err != nil {
		if err == http.ErrMissingFile {
			if isRequired {
				logger.Errorf("[VaultHelper][ParseFile] required file missing: %s", key)
				return nil, nil, fmt.Errorf("file '%s' is required", key)
			}
			logger.Infof("[VaultHelper][ParseFile] optional file skipped: %s", key)
			return nil, nil, nil
		}
		return nil, nil, fmt.Errorf("error retrieving file '%s': %w", key, err)
	}

	return file, fileHeader, nil
}
