package middleware

import (
	"encoding/json"
	"errors"
	"net/http"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/utils"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
)

type FieldError struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func ErrorHandler(w http.ResponseWriter, err error) {
	var errorResponse utils.ErrorDefinition

	if _, ok := err.(validation.Errors); ok {
		fieldErr := ErrorFields(err)
		res := common.Response[[]FieldError]{
			ResponseWriter: w,
			Status:         http.StatusBadRequest,
			Data:           fieldErr,
		}

		res.SendJSON()
		return
	}

	err = errors.Unwrap(err)

	if err := json.Unmarshal([]byte(err.Error()), &errorResponse); err != nil {
		http.Error(w, "failed to unmarshal error", http.StatusInternalServerError)
		return
	}

	res := common.Response[utils.ErrorDefinition]{
		ResponseWriter: w,
		Status:         errorResponse.Code,
		Data:           errorResponse,
	}

	res.SendJSON()
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

		return errs
	}

	return nil
}
