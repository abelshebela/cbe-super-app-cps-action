package localization

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// StandardResponse represents the standardized API response structure
type StandardResponse struct {
	Ok        bool         `json:"ok"`
	Status    int          `json:"status"`
	TimeStamp time.Time    `json:"timestamp,omitempty"`
	Message   string       `json:"message"`
	Data      interface{}  `json:"data,omitempty"`
	Error     *ErrorDetail `json:"error,omitempty"`
}

// ErrorDetail represents error details in the response
type ErrorDetail struct {
	Code        string                 `json:"code"`
	Message     string                 `json:"message"`
	StatusCode  int                    `json:"status_code"`
	Type        string                 `json:"type"`
	FieldErrors []FieldError           `json:"field_errors,omitempty"`
	Details     map[string]interface{} `json:"details,omitempty"`
}

// FieldError represents validation field errors
type FieldError struct {
	Field      string `json:"field"`
	Message    string `json:"message"`
	Value      string `json:"value,omitempty"`
	Constraint string `json:"constraint,omitempty"`
}

// SendSuccessResponse sends a standardized success response
func SendSuccessResponse(w http.ResponseWriter, responseCode ResponseCode, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseCode.StatusCode)

	response := StandardResponse{
		Ok:        true,
		Status:    responseCode.StatusCode,
		TimeStamp: time.Now(),
		Message:   responseCode.Message,
		Data:      data,
	}
	// if responseCode.Type == "error" {
	// 	response.Ok = false
	// }

	if err := json.NewEncoder(w).Encode(response); err != nil {
		SendErrorResponse(w, ErrorUnexpectedError, nil, nil)
	}
}

// SendErrorResponse sends a standardized error response
func SendErrorResponse(w http.ResponseWriter, responseCode ResponseCode, fieldErrors []FieldError, details map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseCode.StatusCode)

	response := StandardResponse{
		Ok:        false,
		TimeStamp: time.Now(),
		Status:    responseCode.StatusCode,
		Message:   responseCode.Message,
		// Error:   errorDetail,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Fallback to basic error response if encoding fails
		w.WriteHeader(StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ok":      false,
			"status":  StatusInternalServerError,
			"message": "Failed to encode response",
		})
	}
}

// SendValidationErrorResponse sends a validation error response
func SendValidationErrorResponse(w http.ResponseWriter, fieldErrors []FieldError) {
	SendErrorResponse(w, ErrorValidationFailed, fieldErrors, nil)
}

// SendErrorByCodeResponse sends a validation error response
func SendErrorByCodeResponse(w http.ResponseWriter, code string) {
	responseCode, ok := GetResponseCodeByCode(code)
	if !ok {
		SendErrorResponse(w, ErrorValidationFailed, nil, nil)
		return
	}
	SendSuccessResponse(w, responseCode, nil)
}

// SendUnauthorizedResponse sends an unauthorized error response
func SendUnauthorizedResponse(w http.ResponseWriter, message string) {
	if message == "" {
		message = MsgUserUnauthorized
	}

	customResponseCode := ResponseCode{
		Code:       "ERROR_UNAUTHORIZED",
		StatusCode: StatusUnauthorized,
		TimeStamp:  time.Now(),
		Message:    message,
		Type:       "error",
	}

	SendErrorResponse(w, customResponseCode, nil, nil)
}

// SendForbiddenResponse sends a forbidden error response
func SendForbiddenResponse(w http.ResponseWriter, message string) {
	if message == "" {
		message = MsgUserForbidden
	}

	customResponseCode := ResponseCode{
		Code:       "ERROR_FORBIDDEN",
		TimeStamp:  time.Now(),
		StatusCode: StatusForbidden,
		Message:    message,
		Type:       "error",
	}

	SendErrorResponse(w, customResponseCode, nil, nil)
}

// SendNotFoundResponse sends a not found error response
func SendNotFoundResponse(w http.ResponseWriter, message string) {
	if message == "" {
		message = MsgUserNotFound
	}

	customResponseCode := ResponseCode{
		Code:       "ERROR_NOT_FOUND",
		TimeStamp:  time.Now(),
		StatusCode: StatusNotFound,
		Message:    message,
		Type:       "error",
	}

	SendErrorResponse(w, customResponseCode, nil, nil)
}

// SendBadRequestResponse sends a bad request error response
func SendBadRequestResponse(w http.ResponseWriter, message string) {
	if message == "" {
		message = MsgBadRequest
	}

	customResponseCode := ResponseCode{
		Code:       "ERROR_BAD_REQUEST",
		TimeStamp:  time.Now(),
		StatusCode: StatusBadRequest,
		Message:    message,
		Type:       "error",
	}

	SendErrorResponse(w, customResponseCode, nil, nil)
}

// SendInternalServerErrorResponse sends an internal server error response
func SendInternalServerErrorResponse(w http.ResponseWriter, message string) {
	if message == "" {
		message = MsgInternalServerError
	}

	customResponseCode := ResponseCode{
		Code:       "ERROR_INTERNAL_SERVER_ERROR",
		TimeStamp:  time.Now(),
		StatusCode: StatusInternalServerError,
		Message:    message,
		Type:       "error",
	}

	SendErrorResponse(w, customResponseCode, nil, nil)
}

// CreateFieldError creates a field error for validation
func CreateFieldError(field, message, value, constraint string) FieldError {
	return FieldError{
		Field:      field,
		Message:    message,
		Value:      value,
		Constraint: constraint,
	}
}

// CreateFieldErrors creates multiple field errors
func CreateFieldErrors(errors ...FieldError) []FieldError {
	return errors
}

// GetResponseCodeByCode fetches a ResponseCode by its Code field from a predefined set of response codes.
func GetResponseCodeByCode(code string) (ResponseCode, bool) {
	// List all response codes to search through.

	for _, rc := range ResponseCodesList {
		if rc.Code == code {
			return rc, true
		}
	}
	return ResponseCode{}, false
}

func ErrorToResponseCode(err string, statusCode int, message string) ResponseCode {
	resp, ok := GetResponseCodeByCode(err)
	if !ok {
		new_resp := ResponseCode{
			Code:       strings.ToUpper(strings.Join(strings.Split(err, " "), "_")),
			TimeStamp:  time.Now(),
			StatusCode: statusCode,
			Message:    message,
			Type:       "error",
		}
		ResponseCodesList = append(ResponseCodesList, new_resp)
		resp = new_resp
	}
	return resp

}
