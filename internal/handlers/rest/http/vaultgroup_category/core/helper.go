package core

import (
	helper "cbe-super-app-cps-action/internal/constants/dto/vaultgroup_category"
	"cbe-super-app-cps-action/internal/constants/localization"
	"fmt"
	"mime"
	"mime/multipart"
	"net/http"

	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func ToDomainCreateVaultGroupCategoryRequest(req helper.CreateVaultGroupCategoryRequest) *model.VaultCategory {
	return &model.VaultCategory{
		Name:     req.Name,
		IsActive: false,
	}
}

func ToDomainUpdateVaultGroupCategoryRequest(req helper.UpdateVaultGroupCategoryRequest) *model.VaultCategory {
	return &model.VaultCategory{
		Name: func() string {
			if req.Name != "" {
				return req.Name
			}
			return ""
		}(),
		UpdatedAt: time.Now().UTC(),
	}
}

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
				logger.Errorf("Required file '%s' is missing", key)
				return nil, nil, fmt.Errorf("file '%s' is required", key)
			}
			logger.Infof("Optional file '%s' not provided - skipping", key)
			return nil, nil, nil
		}
		return nil, nil, fmt.Errorf("error retrieving file '%s': %w", key, err)
	}

	return file, fileHeader, nil
}
