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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	ResponseMaker(map[string]interface{}{
		"code":    successDef.Code,
		"success": true,
		"status":  statusCode,
		"message": successDef.Message,
		"data":    data,
	}, w)
}

// ErrorResponse formats an error response
func ErrorResponse(w http.ResponseWriter, statusCode int, errorKey string, data interface{}) {
	var errorDef ErrorDefinition
	switch {
	case DefineError.Budget[errorKey].Code != "":
		errorDef = DefineError.Budget[errorKey]
	default:
		errorDef = DefineError.General[errorKey]
	}
	w.WriteHeader(statusCode)
	ResponseMaker(map[string]interface{}{
		"code":    errorDef.Code,
		"success": false,
		"status":  statusCode,
		"message": errorDef.Message,
		"data":    data,
	}, w)
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
