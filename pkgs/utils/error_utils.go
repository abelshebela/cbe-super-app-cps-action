package utils

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"database/sql"
	"errors"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// HandleDBError checks the error type and returns appropriate localized errors.
// If the error is mongo.ErrNoDocuments or sql.ErrNoRows, it returns ErrorResourceNotFound.
// For all other errors, it returns ErrorUnexpectedError.
func HandleDBError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, mongo.ErrNoDocuments) || errors.Is(err, sql.ErrNoRows) {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	return errors.New(localization.ErrorUnexpectedError.Code)
}
