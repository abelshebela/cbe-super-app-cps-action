package utils

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"errors"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func GetParam(r *http.Request, key string) (string, bool) {
	value := chi.URLParam(r, key)
	if value == "" {
		return "", false
	}
	return value, true
}

func ParseMultipartFormFile(r *http.Request, key string, maxMemory int64) (multipart.File, *multipart.FileHeader, error) {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		return nil, nil, errors.New(localization.ErrorMissingFile.Code)
	}
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		return nil, nil, errors.New(localization.ErrorFileParseFailed.Code)
	}

	file, fileHeader, err := r.FormFile(key)
	if err != nil {
		return nil, nil, err
	}

	return file, fileHeader, nil
}
