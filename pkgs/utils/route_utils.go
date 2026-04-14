package utils

import (
	"bytes"
	"cbe-super-app-cps-action/internal/constants/localization"
	"encoding/base64"
	"errors"
	"io"
	"log"
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

type nopFile struct {
	*bytes.Reader
}

func (n nopFile) Close() error { return nil }

func ParseMultipartFormFile(r *http.Request, key string, maxMemory int64) (multipart.File, *multipart.FileHeader, error) {
	log.Printf("[DEBUG] ParseMultipartFormFile called with key=%s, maxMemory=%d", key, maxMemory)
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		log.Printf("[ERROR] ParseMultipartForm error: %v", err)
		return nil, nil, errors.New(localization.ErrorFileParseFailed.Code)
	}
	file, fileHeader, err := r.FormFile(key)
	if err != nil {
		log.Printf("[ERROR] FormFile error: %v", err)
		return nil, nil, err
	}

	// Read the file content to check for base64 encoding
	var buf bytes.Buffer
	_, copyErr := io.Copy(&buf, file)
	if copyErr != nil {
		log.Printf("[ERROR] io.Copy error: %v", copyErr)
		return nil, nil, copyErr
	}
	content := buf.Bytes()

	// Try to decode as base64
	decoded, decodeErr := base64.StdEncoding.DecodeString(strings.TrimSpace(string(content)))
	if decodeErr == nil {
		log.Printf("[DEBUG] Base64 decode succeeded for key=%s, size=%d", key, len(decoded))
		// If base64 decoding succeeds, use the decoded content
		file = nopFile{bytes.NewReader(decoded)}
		fileHeader.Size = int64(len(decoded))
		fileHeader.Header.Set("Content-Type", http.DetectContentType(decoded))
		content = decoded
	} else {
		log.Printf("[DEBUG] Base64 decode failed for key=%s: %v (treating as binary)", key, decodeErr)
		// If not base64, reset file to original content
		file = nopFile{bytes.NewReader(content)}
	}

	if !IsValidImage(fileHeader) {
		log.Printf("[ERROR] IsValidImage failed for key=%s, filename=%s", key, fileHeader.Filename)
		return nil, nil, errors.New(localization.ErrorInvalidFileUpload.Code)
	}
	log.Printf("[DEBUG] ParseMultipartFormFile succeeded for key=%s, filename=%s, size=%d", key, fileHeader.Filename, fileHeader.Size)
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
