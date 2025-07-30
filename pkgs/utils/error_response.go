package utils

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/http"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/common"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type APIErrorResponse struct {
	Message       string `json:"message"`
	Code          string `json:"code,omitempty"`
	AccountLocked bool   `json:"accountLocked,omitempty"`
	Auth          bool   `json:"auth,omitempty"`
	Errors        any    `json:"errors,omitempty"`
}

func formatValidationErrors(ve validation.Errors) map[string]string {
	result := make(map[string]string)
	for field, fieldErr := range ve {
		def := lookupErrorDefinition(fieldErr.Error())
		result[field] = def.Message
	}
	return result
}

func SendErrorResponse(w http.ResponseWriter, errorKey any, statusCode int, additionalData map[string]any) {
	var (
		code              string
		message           string
		errorKeyString    string
		errors            any
		isValidationError bool
	)

	switch v := errorKey.(type) {
	case string:
		def := lookupErrorDefinition(v)
		code = def.Code
		message = def.Message
		errorKeyString = v

	case error:
		if ve, ok := v.(validation.Errors); ok {
			code = "GEN_113"
			message = "One or more required fields are invalid."
			errors = formatValidationErrors(ve)
			errorKeyString = "GEN_113"
			isValidationError = true
		} else {
			def := lookupErrorDefinition(v.Error())
			code = def.Code
			message = def.Message
			errorKeyString = def.Code
		}

	default:
		code = "GEN_UNKNOWN"
		message = fmt.Sprintf("%v", v)
		errorKeyString = "GEN_UNKNOWN"
	}

	status := httpStatusOrDefault(statusCode, errorKeyString, isValidationError)
	response := map[string]any{
		"status":  status,
		"message": message,
		"data": map[string]any{
			"code": code,
		},
	}

	if errors != nil {
		response["errors"] = errors
	}

	maps.Copy(response, additionalData)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Printf("Failed to write error response: %v\n", err)
	}
}

func lookupErrorDefinition(key string) common.ErrorDefinition {
	errorGroups := []common.ErrorGroup{
		common.DefineError.General,
		common.DefineError.Auth,
		common.DefineError.User,
		common.DefineError.Transaction,
		common.DefineError.Account,
		common.DefineError.OTP,
		common.DefineError.File,
		common.DefineError.Branch,
		common.DefineError.Bank,
		common.DefineError.Department,
		common.DefineError.Wallet,
		common.DefineError.Action,
		common.DefineError.AD,
		common.DefineError.Permission,
		common.DefineError.MiniApp,
		common.DefineError.Event,
	}

	for _, group := range errorGroups {
		if def, ok := group[key]; ok {
			return def
		}
	}

	return common.ErrorDefinition{
		Code:    key,
		Message: key,
		Status:  http.StatusInternalServerError,
	}
}

func httpStatusOrDefault(providedStatus int, errorCode string, isValidation bool) int {
	if isValidation {
		return http.StatusBadRequest
	}
	if providedStatus != 0 {
		return providedStatus
	}

	def := lookupErrorDefinition(errorCode)
	return def.Status
}
