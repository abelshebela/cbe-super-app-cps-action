package utils

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"

	"go.mongodb.org/mongo-driver/v2/bson"

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
		return nil, nil, fmt.Errorf(ErrMissingFile)
	}
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		return nil, nil, fmt.Errorf("failed to parse multipart form: %w", err)
	}

	file, fileHeader, err := r.FormFile(key)
	if err != nil {
		if err == http.ErrMissingFile {
			return nil, nil, fmt.Errorf(ErrMissingFile)
		}
		return nil, nil, fmt.Errorf("missing or invalid file for key '%s': %w", key, err)
	}

	return file, fileHeader, nil
}

// func UserContextToModel(userContext ctx_util.UserContext) model.User {
// 	return model.User{
// 		UserCode:    userContext.UserCode,
// 		FullName:    userContext.FullName,
// 		PhoneNumber: userContext.PhoneNumber,
// 		Department:  userContext.Department,
// 	}
// }

func ParsePrimitiveObjectID(ID string) (bson.ObjectID, error) {
	objectID, err := bson.ObjectIDFromHex(ID)

	if err != nil {
		return bson.ObjectID{}, fmt.Errorf(InvalidID)
	}
	return objectID, nil
}
