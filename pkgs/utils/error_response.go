package utils

import (
	"cbe-super-app-member-users/pkgs/common"
	"encoding/json"
	"net/http"
)

type APIErrorResponse struct {
	Message       string      `json:"message"`
	Code          string      `json:"code,omitempty"`
	AccountLocked bool        `json:"accountLocked,omitempty"`
	Auth          bool        `json:"auth,omitempty"`
	Errors        interface{} `json:"errors,omitempty"`
}

func SendErrorResponse(w http.ResponseWriter, errorKey string, statusCode int, additionalData map[string]interface{}) {

	var errorDef common.ErrorDefinition
	switch {
	case common.DefineError.General[errorKey].Code != "":
		errorDef = common.DefineError.General[errorKey]
	case common.DefineError.Auth[errorKey].Code != "":
		errorDef = common.DefineError.Auth[errorKey]
	case common.DefineError.User[errorKey].Code != "":
		errorDef = common.DefineError.User[errorKey]
	case common.DefineError.Transaction[errorKey].Code != "":
		errorDef = common.DefineError.Transaction[errorKey]
	case common.DefineError.OTP[errorKey].Code != "":
		errorDef = common.DefineError.OTP[errorKey]
	default:
		errorDef = common.DefineError.General[errorKey]
	}

	response := APIErrorResponse{
		Message: errorDef.Message,
		Code:    errorDef.Code,
	}

	if additionalData != nil {
		if accountLocked, ok := additionalData["accountLocked"].(bool); ok {
			response.AccountLocked = accountLocked
		}
		if auth, ok := additionalData["auth"].(bool); ok {
			response.Auth = auth
		}
		if errors, exists := additionalData["errors"]; exists {
			response.Errors = errors
		}
	}

	responseBytes, _ := json.Marshal(response)
	var responseMap map[string]interface{}
	json.Unmarshal(responseBytes, &responseMap)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	ResponseMaker(responseMap, w)
}

/*
Example usage:

1. Basic error response:
   SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, nil)

2. Account locked error:
   SendErrorResponse(w, "AUTH_USER_RESET_PASSWORD_REQUIRED", http.StatusUnauthorized,
       map[string]interface{}{"accountLocked": true})

3. Validation error with details:
   SendErrorResponse(w, "FAILD_VALIDATION", http.StatusBadRequest,
       map[string]interface{}{"errors": validationErrors})

4. OTP error:
   SendErrorResponse(w, "AUTH_EXPIRED_OTP", http.StatusNotFound,
       map[string]interface{}{"auth": false})
*/
