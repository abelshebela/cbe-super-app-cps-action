package utils

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"errors"
	"mime/multipart"
	"net/http"
	"regexp"

	"bytes"
	"encoding/base64"
	"io"
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

// func ParseMultipartFormFile(r *http.Request, key string, maxMemory int64) (multipart.File, *multipart.FileHeader, error) {
// 	// if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
// 	// 	return nil, nil, errors.New(localization.ErrorMissingFile.Code)
// 	// }
// 	if err := r.ParseMultipartForm(maxMemory); err != nil {
// 		return nil, nil, errors.New(localization.ErrorFileParseFailed.Code)
// 	}

// 	file, fileHeader, err := r.FormFile(key)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	if !IsValidImage(fileHeader) {
// 		return nil, nil, errors.New(localization.ErrorInvalidFileUpload.Code)
// 	}

// 	return file, fileHeader, nil
// }
// ...existing code...

type nopFile struct {
	*bytes.Reader
}

func (n nopFile) Close() error { return nil }

func ParseMultipartFormFile(r *http.Request, key string, maxMemory int64) (multipart.File, *multipart.FileHeader, error) {
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		return nil, nil, errors.New(localization.ErrorFileParseFailed.Code)
	}
	file, fileHeader, err := r.FormFile(key)
	if err != nil {
		return nil, nil, err
	}

	// Read the file content to check for base64 encoding
	var buf bytes.Buffer
	_, copyErr := io.Copy(&buf, file)
	if copyErr != nil {
		return nil, nil, copyErr
	}
	content := buf.Bytes()

	// Try to decode as base64
	decoded, decodeErr := base64.StdEncoding.DecodeString(strings.TrimSpace(string(content)))
	if decodeErr == nil {
		// If base64 decoding succeeds, use the decoded content
		file = nopFile{bytes.NewReader(decoded)}
		fileHeader.Size = int64(len(decoded))
		fileHeader.Header.Set("Content-Type", http.DetectContentType(decoded))
		content = decoded
	} else {
		// If not base64, reset file to original content
		file = nopFile{bytes.NewReader(content)}
	}

	if !IsValidImage(fileHeader) {
		return nil, nil, errors.New(localization.ErrorInvalidFileUpload.Code)
	}
	return file, fileHeader, nil
}

// ...existing code...

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
