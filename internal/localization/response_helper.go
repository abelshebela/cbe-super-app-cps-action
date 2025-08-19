package localization

import (
	"encoding/json"
	"net/http"
)

// StandardResponse represents the standardized API response structure
type StandardResponse struct {
	Ok      bool         `json:"ok"`
	Status  int          `json:"status"`
	Message string       `json:"message"`
	Data    interface{}  `json:"data,omitempty"`
	Error   *ErrorDetail `json:"error,omitempty"`
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
		Ok:      true,
		Status:  responseCode.StatusCode,
		Message: responseCode.Message,
		Data:    data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Fallback to error response if encoding fails
		SendErrorResponse(w, ErrorUnexpectedError, nil, nil)
	}
}

// SendErrorResponse sends a standardized error response
func SendErrorResponse(w http.ResponseWriter, responseCode ResponseCode, fieldErrors []FieldError, details map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseCode.StatusCode)

	errorDetail := &ErrorDetail{
		Code:        responseCode.Code,
		Message:     responseCode.Message,
		StatusCode:  responseCode.StatusCode,
		Type:        responseCode.Type,
		FieldErrors: fieldErrors,
		Details:     details,
	}

	response := StandardResponse{
		Ok:      false,
		Status:  responseCode.StatusCode,
		Message: responseCode.Message,
		Error:   errorDetail,
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

// SendUnauthorizedResponse sends an unauthorized error response
func SendUnauthorizedResponse(w http.ResponseWriter, message string) {
	if message == "" {
		message = MsgUserUnauthorized
	}

	customResponseCode := ResponseCode{
		Code:       "ERROR_UNAUTHORIZED",
		StatusCode: StatusUnauthorized,
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

// SendErrorByCode sends an error response by fetching ResponseCode based on the given code
func SendErrorByCode(w http.ResponseWriter, code string, fieldErrors []FieldError, details map[string]interface{}) {
	responseCode := GetResponseCodeByCode(code)
	SendErrorResponse(w, responseCode, fieldErrors, details)
}

// GetResponseCodeByCode fetches the ResponseCode based on the given code string
func GetResponseCodeByCode(code string) ResponseCode {
	switch code {
	case ErrorUserNotFound.Code:
		return ErrorUserNotFound
	case ErrorUserAlreadyExists.Code:
		return ErrorUserAlreadyExists
	case ErrorUserUnauthorized.Code:
		return ErrorUserUnauthorized
	case ErrorUserForbidden.Code:
		return ErrorUserForbidden
	case ErrorUserInvalidCredentials.Code:
		return ErrorUserInvalidCredentials
	case ErrorUserAccountBlocked.Code:
		return ErrorUserAccountBlocked
	case ErrorUserSessionExpired.Code:
		return ErrorUserSessionExpired
	case ErrorUserTooManyLoginAttempts.Code:
		return ErrorUserTooManyLoginAttempts
	case ErrorUserDeviceMismatch.Code:
		return ErrorUserDeviceMismatch
	case ErrorUserDeviceNotFound.Code:
		return ErrorUserDeviceNotFound
	case ErrorUserDeviceAlreadyExists.Code:
		return ErrorUserDeviceAlreadyExists
	case ErrorOTPNotFound.Code:
		return ErrorOTPNotFound
	case ErrorOTPExpired.Code:
		return ErrorOTPExpired
	case ErrorOTPInvalid.Code:
		return ErrorOTPInvalid
	case ErrorOTPAlreadyExists.Code:
		return ErrorOTPAlreadyExists
	case ErrorOTPTooManyAttempts.Code:
		return ErrorOTPTooManyAttempts
	case ErrorOTPSendFailed.Code:
		return ErrorOTPSendFailed
	case ErrorPINInvalid.Code:
		return ErrorPINInvalid
	case ErrorPINMismatch.Code:
		return ErrorPINMismatch
	case ErrorPINTooWeak.Code:
		return ErrorPINTooWeak
	case ErrorPINInHistory.Code:
		return ErrorPINInHistory
	case ErrorPINLengthInvalid.Code:
		return ErrorPINLengthInvalid
	case ErrorPINOnlyDigits.Code:
		return ErrorPINOnlyDigits
	case ErrorPINSequential.Code:
		return ErrorPINSequential
	case ErrorPINRedundant.Code:
		return ErrorPINRedundant
	case ErrorSamePIN.Code:
		return ErrorSamePIN
	case ErrorSessionNotFound.Code:
		return ErrorSessionNotFound
	case ErrorSessionExpired.Code:
		return ErrorSessionExpired
	case ErrorSessionInvalid.Code:
		return ErrorSessionInvalid
	case ErrorSessionCreationFailed.Code:
		return ErrorSessionCreationFailed
	case ErrorSessionUpdateFailed.Code:
		return ErrorSessionUpdateFailed
	case ErrorFileTooLarge.Code:
		return ErrorFileTooLarge
	case ErrorFileInvalidType.Code:
		return ErrorFileInvalidType
	case ErrorFileUploadFailed.Code:
		return ErrorFileUploadFailed
	case ErrorFileNotFound.Code:
		return ErrorFileNotFound
	case ErrorFileDeleteFailed.Code:
		return ErrorFileDeleteFailed
	case ErrorValidationFailed.Code:
		return ErrorValidationFailed
	case ErrorRequiredFieldMissing.Code:
		return ErrorRequiredFieldMissing
	case ErrorInvalidFormat.Code:
		return ErrorInvalidFormat
	case ErrorInvalidEmail.Code:
		return ErrorInvalidEmail
	case ErrorInvalidPhoneNumber.Code:
		return ErrorInvalidPhoneNumber
	case ErrorInvalidDate.Code:
		return ErrorInvalidDate
	case ErrorFieldTooLong.Code:
		return ErrorFieldTooLong
	case ErrorFieldTooShort.Code:
		return ErrorFieldTooShort
	case ErrorInternalServerError.Code:
		return ErrorInternalServerError
	case ErrorServiceUnavailable.Code:
		return ErrorServiceUnavailable
	case ErrorDatabaseError.Code:
		return ErrorDatabaseError
	case ErrorNetworkError.Code:
		return ErrorNetworkError
	case ErrorTimeoutError.Code:
		return ErrorTimeoutError
	case ErrorUnexpectedError.Code:
		return ErrorUnexpectedError
	case ErrorConfigurationError.Code:
		return ErrorConfigurationError
	case ErrorExternalServiceError.Code:
		return ErrorExternalServiceError
	case ErrorInsufficientBalance.Code:
		return ErrorInsufficientBalance
	case ErrorTransactionFailed.Code:
		return ErrorTransactionFailed
	case ErrorLimitExceeded.Code:
		return ErrorLimitExceeded
	case ErrorOperationNotAllowed.Code:
		return ErrorOperationNotAllowed
	case ErrorResourceBusy.Code:
		return ErrorResourceBusy
	case ErrorMaintenanceMode.Code:
		return ErrorMaintenanceMode
	default:
		return ErrorUnexpectedError
	}
}
