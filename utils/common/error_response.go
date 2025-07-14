package common

import (
	"encoding/json"
	"log"
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

	var errorDef ErrorDefinition
	switch {
	case DefineError.General[errorKey].Code != "":
		errorDef = DefineError.General[errorKey]
	case DefineError.Auth[errorKey].Code != "":
		errorDef = DefineError.Auth[errorKey]
	case DefineError.User[errorKey].Code != "":
		errorDef = DefineError.User[errorKey]
	case DefineError.Transaction[errorKey].Code != "":
		errorDef = DefineError.Transaction[errorKey]
	case DefineError.OTP[errorKey].Code != "":
		errorDef = DefineError.OTP[errorKey]
	default:
		errorDef = DefineError.General[errorKey]
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

	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Printf("failed to marshal response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	var responseMap map[string]interface{}

	if err := json.Unmarshal(responseBytes, &responseMap); err != nil {
		log.Printf("failed to unmarshal response bytes: %v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(responseMap); err != nil {
		log.Printf("failed to encode error response: %v", err)
	}
}

func SendSuccessResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to encode success response: %v", err)
	}
}
