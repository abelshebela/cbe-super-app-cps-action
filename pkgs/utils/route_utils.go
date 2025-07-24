package utils

import (
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"

	"go.mongodb.org/mongo-driver/bson/primitive"

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

func UserContextToModel(userContext ctx_util.UserContext) model.User {
	return model.User{
		UserCode:    userContext.UserCode,
		FullName:    userContext.FullName,
		PhoneNumber: userContext.PhoneNumber,
	}
}

func ParsePrimitiveObjectID(ID string) (primitive.ObjectID, error) {
	objectID, err := primitive.ObjectIDFromHex(ID)

	if err != nil {
		return primitive.ObjectID{}, fmt.Errorf(InvalidID)
	}
	return objectID, nil
}
