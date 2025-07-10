package common

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Response[T any] struct {
	ResponseWriter http.ResponseWriter
	Status         int
	Data           T
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (r *Response[T]) SendJSON() {
	r.ResponseWriter.Header().Set("Content-Type", "application/json")
	r.ResponseWriter.WriteHeader(r.Status)
	if err := json.NewEncoder(r.ResponseWriter).Encode(r.Data); err != nil {
		http.Error(
			r.ResponseWriter,
			`{"error": "failed to encode response"}`,
			http.StatusInternalServerError,
		)
	}
}

// successResponse defines the structure of a successful response
func SuccessResponse(w http.ResponseWriter, statusCode int, successKey string, data interface{}) {
	var successDef SuccessDefinition
	successDef = DefineSuccess.General[successKey]
	r := Response[map[string]interface{}]{
		ResponseWriter: w,
		Status:         statusCode,
		Data: map[string]interface{}{
			"code":    successDef.Code,
			"success": true,
			"status":  statusCode,
			"message": successDef.Message,
			"data":    data,
		},
	}
	r.SendJSON()
}

// ErrorResponse formats an error response
func ErrorResponse(w http.ResponseWriter, statusCode int, errorKey string, data interface{}) {
	var errorDef ErrorDefinition
	switch {
	case DefineError.Account[errorKey].Code != "":
		errorDef = DefineError.Account[errorKey]
	case DefineError.Action[errorKey].Code != "":
		errorDef = DefineError.Action[errorKey]
	case DefineError.Auth[errorKey].Code != "":
		errorDef = DefineError.Auth[errorKey]
	case DefineError.Budget[errorKey].Code != "":
		errorDef = DefineError.Budget[errorKey]
	case DefineError.User[errorKey].Code != "":
		errorDef = DefineError.User[errorKey]
	case DefineError.Transaction[errorKey].Code != "":
		errorDef = DefineError.Transaction[errorKey]
	case DefineError.OTP[errorKey].Code != "":
		errorDef = DefineError.OTP[errorKey]
	case DefineError.File[errorKey].Code != "":
		errorDef = DefineError.File[errorKey]
	default:
		errorDef = DefineError.General[errorKey]
	}
	r := Response[map[string]interface{}]{
		ResponseWriter: w,
		Status:         statusCode,
		Data: map[string]interface{}{
			"code":    errorDef.Code,
			"success": false,
			"status":  statusCode,
			"message": errorDef.Message,
			"data":    data,
		},
	}
	r.SendJSON()
}

func FormatValidationErrors(err error) []ValidationError {
	var errors []ValidationError

	if validationErrors, ok := err.(validation.Errors); ok {
		for field, err := range validationErrors {
			errors = append(errors, ValidationError{
				Field:   field,
				Message: err.Error(),
			})
		}
	} else if err != nil {
		// Handle non-validation errors
		errors = append(errors, ValidationError{
			Field:   "",
			Message: err.Error(),
		})
	}

	return errors
}

func convertValidationError(err error) []FieldError {
	if err == nil {
		return nil
	}

	var errors []FieldError
	if errs, ok := err.(validation.Errors); ok {
		for field, e := range errs {
			errors = append(errors, FieldError{
				Field:   field,
				Message: e.Error(),
				Code:    "VALIDATION_ERROR",
			})
		}
	} else {
		errors = append(errors, FieldError{
			Message: err.Error(),
			Code:    "VALIDATION_ERROR",
		})
	}
	return errors
}
