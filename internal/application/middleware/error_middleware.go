package middleware

import (
	"net/http"

	"cbe-super-app-member-users/pkgs/utils"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type FieldError struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func ErrorHandler(w http.ResponseWriter, err error) {
	if validationErrs, ok := err.(validation.Errors); ok {
		fieldErrs := ErrorFields(validationErrs)
		utils.SendErrorResponse(w, "FAILED_VALIDATION", http.StatusBadRequest, map[string]interface{}{
			"errors": fieldErrs,
		})
		return
	}

	utils.SendErrorResponse(w, "UNHANDLED_SERVER_ERROR", http.StatusInternalServerError, nil)
}

func ErrorFields(err error) []FieldError {
	var errs []FieldError

	if data, ok := err.(validation.Errors); ok {
		for i, v := range data {
			nestedErrors := ErrorFields(v)
			if len(nestedErrors) > 0 {
				errs = append(errs, nestedErrors...)
			} else {
				errs = append(errs, FieldError{
					Name:        i,
					Description: v.Error(),
				})
			}
		}
	}

	return errs
}
