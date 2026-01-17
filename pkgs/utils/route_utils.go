package utils

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"errors"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation"
)

func GetParam(r *http.Request, key string) (string, bool) {
	value := chi.URLParam(r, key)
	if value == "" {
		return "", false
	}
	return value, true
}

var mongoIDRegex = regexp.MustCompile(`^[a-fA-F0-9]{24}$`)

func ValidateMongoID(id string) error {
	return validation.Validate(
		id,
		validation.Required,
		validation.Match(mongoIDRegex).Error("invalid MongoDB ObjectID"),
	)
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

// type nopFile struct {
// 	*bytes.Reader
// }

// func (n nopFile) Close() error { return nil }

// func ParseMultipartFormFile(
// 	r *http.Request,
// 	key string,
// 	maxMemory int64,
// ) (multipart.File, *multipart.FileHeader, error) {

// 	if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
// 		return nil, nil, errors.New(localization.ErrorMissingFile.Code)
// 	}

// 	if err := r.ParseMultipartForm(maxMemory); err != nil {
// 		return nil, nil, errors.New(localization.ErrorFileParseFailed.Code)
// 	}

// 	origFile, origHeader, err := r.FormFile(key)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	defer origFile.Close()

// 	// --- gzip compression ---
// 	var buf bytes.Buffer
// 	gw := gzip.NewWriter(&buf)

// 	if _, err := io.Copy(gw, origFile); err != nil {
// 		gw.Close()
// 		return nil, nil, err
// 	}

// 	if err := gw.Close(); err != nil {
// 		return nil, nil, err
// 	}

// 	// Create a multipart.File-compatible reader
// 	reader := bytes.NewReader(buf.Bytes())
// 	compressedFile := nopFile{Reader: reader}

// 	// --- modify FileHeader safely ---
// 	newHeader := &multipart.FileHeader{
// 		Filename: origHeader.Filename + ".gz",
// 		Size:     int64(buf.Len()),
// 		Header:   make(textproto.MIMEHeader),
// 	}

// 	newHeader.Header.Set("Content-Type", "application/gzip")
// 	newHeader.Header.Set("Content-Encoding", "gzip")

// 	return compressedFile, newHeader, nil
// }
