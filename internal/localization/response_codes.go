package localization

// ResponseCode represents a standardized response code with status and message
// type ResponseCode struct {
// 	Code       string `json:"code"`
// 	StatusCode int    `json:"status_code"`
// 	Message    string `json:"message"`
// 	Type       string `json:"type"` // "success", "error", "warning", "info"
// }

// // Success Response Codes
// var (
// 	SuccessUserCreated = ResponseCode{
// 		Code:       "SUCCESS_USER_CREATED",
// 		StatusCode: StatusCreated,
// 		Message:    MsgUserCreatedSuccessfully,
// 		Type:       "success",
// 	}

// 	SuccessUserUpdated = ResponseCode{
// 		Code:       "SUCCESS_USER_UPDATED",
// 		StatusCode: StatusOK,
// 		Message:    MsgUserUpdatedSuccessfully,
// 		Type:       "success",
// 	}

// 	SuccessUserDeleted = ResponseCode{
// 		Code:       "SUCCESS_USER_DELETED",
// 		StatusCode: StatusOK,
// 		Message:    MsgUserDeletedSuccessfully,
// 		Type:       "success",
// 	}

// 	SuccessUserRetrieved = ResponseCode{
// 		Code:       "SUCCESS_USER_RETRIEVED",
// 		StatusCode: StatusOK,
// 		Message:    MsgUserRetrievedSuccessfully,
// 		Type:       "success",
// 	}

// 	SuccessUserLogin = ResponseCode{
// 		Code:       "SUCCESS_USER_LOGIN",
// 		StatusCode: StatusOK,
// 		Message:    MsgUserLoginSuccessfully,
// 		Type:       "success",
// 	}
	
// 	SuccessUserLogout = ResponseCode{
// 		Code:       "SUCCESS_USER_LOGOUT",
// 		StatusCode: StatusOK,
// 		Message:    MsgUserLogoutSuccessfully,
// 		Type:       "success",
// 	}

// 	SuccessOTPSent = ResponseCode{
// 		Code:       "SUCCESS_OTP_SENT",
// 		StatusCode: StatusOK,
// 		Message:    MsgOTPSentSuccessfully,
// 		Type:       "success",
// 	}

// 	SuccessOTPVerified = ResponseCode{
// 		Code:       "SUCCESS_OTP_VERIFIED",
// 		StatusCode: StatusOK,
// 		Message:    MsgOTPVerifiedSuccessfully,
// 		Type:       "success",
// 	}

// 	SuccessOperationCompleted = ResponseCode{
// 		Code:       "SUCCESS_OPERATION_COMPLETED",
// 		StatusCode: StatusOK,
// 		Message:    MsgOperationCompleted,
// 		Type:       "success",
// 	}

// 	SuccessDataRetrieved = ResponseCode{
// 		Code:       "SUCCESS_DATA_RETRIEVED",
// 		StatusCode: StatusOK,
// 		Message:    MsgDataRetrievedSuccessfully,
// 		Type:       "success",
// 	}

// 	SuccessDataSaved = ResponseCode{
// 		Code:       "SUCCESS_DATA_SAVED",
// 		StatusCode: StatusCreated,
// 		Message:    MsgDataSavedSuccessfully,
// 		Type:       "success",
// 	}

// 	SuccessDataUpdated = ResponseCode{
// 		Code:       "SUCCESS_DATA_UPDATED",
// 		StatusCode: StatusOK,
// 		Message:    MsgDataUpdatedSuccessfully,
// 		Type:       "success",
// 	}

// 	SuccessDataDeleted = ResponseCode{
// 		Code:       "SUCCESS_DATA_DELETED",
// 		StatusCode: StatusOK,
// 		Message:    MsgDataDeletedSuccessfully,
// 		Type:       "success",
// 	}

// 	SuccessHealthCheck = ResponseCode{
// 		Code:       "SUCCESS_HEALTH_CHECK",
// 		StatusCode: StatusOK,
// 		Message:    MsgHealthCheckPassed,
// 		Type:       "success",
// 	}

// 	SuccessUserPINChanged = ResponseCode{
// 		Code:       "SUCCESS_USER_PIN_CHANGED",
// 		StatusCode: StatusOK,
// 		Message:    MsgUserPINChanged,
// 		Type:       "success",
// 	}
// )

// // Error Response Codes
// var (
// 	ErrorUserNotFound = ResponseCode{
// 		Code:       "ERROR_USER_NOT_FOUND",
// 		StatusCode: StatusNotFound,
// 		Message:    MsgUserNotFound,
// 		Type:       "error",
// 	}

// 	ErrorUserAlreadyExists = ResponseCode{
// 		Code:       "ERROR_USER_ALREADY_EXISTS",
// 		StatusCode: StatusConflict,
// 		Message:    MsgUserAlreadyExists,
// 		Type:       "error",
// 	}

// 	ErrorUserUnauthorized = ResponseCode{
// 		Code:       "ERROR_USER_UNAUTHORIZED",
// 		StatusCode: StatusUnauthorized,
// 		Message:    MsgUserUnauthorized,
// 		Type:       "error",
// 	}

// 	ErrorUserForbidden = ResponseCode{
// 		Code:       "ERROR_USER_FORBIDDEN",
// 		StatusCode: StatusForbidden,
// 		Message:    MsgUserForbidden,
// 		Type:       "error",
// 	}

// 	ErrorUserInvalidCredentials = ResponseCode{
// 		Code:       "ERROR_USER_INVALID_CREDENTIALS",
// 		StatusCode: StatusUnauthorized,
// 		Message:    MsgUserInvalidCredentials,
// 		Type:       "error",
// 	}

// 	ErrorUserAccountBlocked = ResponseCode{
// 		Code:       "ERROR_USER_ACCOUNT_BLOCKED",
// 		StatusCode: StatusForbidden,
// 		Message:    MsgUserAccountBlocked,
// 		Type:       "error",
// 	}

// 	ErrorUserSessionExpired = ResponseCode{
// 		Code:       "ERROR_USER_SESSION_EXPIRED",
// 		StatusCode: StatusUnauthorized,
// 		Message:    MsgUserSessionExpired,
// 		Type:       "error",
// 	}

// 	ErrorUserTooManyLoginAttempts = ResponseCode{
// 		Code:       "ERROR_USER_TOO_MANY_LOGIN_ATTEMPTS",
// 		StatusCode: StatusTooManyRequests,
// 		Message:    MsgUserTooManyLoginAttempts,
// 		Type:       "error",
// 	}

// 	ErrorUserDeviceMismatch = ResponseCode{
// 		Code:       "ERROR_USER_DEVICE_MISMATCH",
// 		StatusCode: StatusConflict,
// 		Message:    MsgUserDeviceMismatch,
// 		Type:       "error",
// 	}

// 	ErrorUserDeviceNotFound = ResponseCode{
// 		Code:       "ERROR_USER_DEVICE_NOT_FOUND",
// 		StatusCode: StatusNotFound,
// 		Message:    MsgUserDeviceNotFound,
// 		Type:       "error",
// 	}

// 	ErrorUserDeviceAlreadyExists = ResponseCode{
// 		Code:       "ERROR_USER_DEVICE_ALREADY_EXISTS",
// 		StatusCode: StatusConflict,
// 		Message:    MsgUserDeviceAlreadyExists,
// 		Type:       "error",
// 	}

// 	ErrorOTPNotFound = ResponseCode{
// 		Code:       "ERROR_OTP_NOT_FOUND",
// 		StatusCode: StatusNotFound,
// 		Message:    MsgOTPNotFound,
// 		Type:       "error",
// 	}

// 	ErrorOTPExpired = ResponseCode{
// 		Code:       "ERROR_OTP_EXPIRED",
// 		StatusCode: StatusBadRequest,
// 		Message:    MsgOTPExpired,
// 		Type:       "error",
// 	}

// 	ErrorOTPInvalid = ResponseCode{
// 		Code:       "ERROR_OTP_INVALID",
// 		StatusCode: StatusBadRequest,
// 		Message:    MsgOTPInvalid,
// 		Type:       "error",
// 	}

// 	ErrorOTPAlreadyExists = ResponseCode{
// 		Code:       "ERROR_OTP_ALREADY_EXISTS",
// 		StatusCode: StatusBadRequest,
// 		Message:    MsgOTPAlreadyExists,
// 		Type:       "error",
// 	}

// 	ErrorOTPTooManyAttempts = ResponseCode{
// 		Code:       "ERROR_OTP_TOO_MANY_ATTEMPTS",
// 		StatusCode: StatusTooManyRequests,
// 		Message:    MsgOTPTooManyAttempts,
// 		Type:       "error",
// 	}

// 	ErrorOTPSendFailed = ResponseCode{
// 		Code:       "ERROR_OTP_SEND_FAILED",
// 		StatusCode: StatusInternalServerError,
// 		Message:    MsgOTPSendFailed,
// 		Type:       "error",
// 	}

// 	ErrorPINInvalid = ResponseCode{
// 		Code:       "ERROR_PIN_INVALID",
// 		StatusCode: StatusBadRequest,
// 		Message:    MsgPINInvalid,
// 		Type:       "error",
// 	}

// 	ErrorPINMismatch = ResponseCode{
// 		Code:       "ERROR_PIN_MISMATCH",
// 		StatusCode: StatusBadRequest,
// 		Message:    MsgPINMismatch,
// 		Type:       "error",
// 	}

// 	ErrorPINTooWeak = ResponseCode{
// 		Code:       "ERROR_PIN_TOO_WEAK",
// 		StatusCode: StatusBadRequest,
// 		Message:    MsgPINTooWeak,
// 		Type:       "error",
// 	}

// 	ErrorPINInHistory = ResponseCode{
// 		Code:       "ERROR_PIN_IN_HISTORY",
// 		StatusCode: StatusBadRequest,
// 		Message:    MsgPINInHistory,
// 		Type:       "error",
// 	}

// 	ErrorPINLengthInvalid = ResponseCode{
// 		Code:       "ERROR_PIN_LENGTH_INVALID",
// 		StatusCode: StatusBadRequest,
// 		Message:    MsgPINLengthInvalid,
// 		Type:       "error",
// 	}

// 	ErrorPINOnlyDigits = ResponseCode{
// 		Code:       "ERROR_PIN_ONLY_DIGITS",
// 		StatusCode: StatusBadRequest,
// 		Message:    MsgPINOnlyDigits,
// 		Type:       "error",
// 	}

// 	ErrorPINSequential = ResponseCode{
// 		Code:       "ERROR_PIN_SEQUENTIAL",
// 		StatusCode: StatusBadRequest,
// 		Message:    MsgPINSequential,
// 		Type:       "error",
// 	}

// 	ErrorPINRedundant = ResponseCode{
// 		Code:       "ERROR_PIN_REDUNDANT",
// 		StatusCode: StatusBadRequest,
// 		Message:    MsgPINRedundant,
// 		Type:       "error",
// 	}

// 	ErrorSamePIN = ResponseCode{
// 		Code:       "ERROR_SAME_PIN",
// 		StatusCode: StatusBadRequest,
// 		Message:    MsgSamePIN,
// 		Type:       "error",
// 	}

// 	ErrorSessionNotFound = ResponseCode{
// 		Code:       "ERROR_SESSION_NOT_FOUND",
// 		StatusCode: StatusNotFound,
// 		Message:    MsgSessionNotFound,
// 		Type:       "error",
// 	}

// 	ErrorSessionExpired = ResponseCode{
// 		Code:       "ERROR_SESSION_EXPIRED",
// 		StatusCode: StatusUnauthorized,
// 		Message:    MsgSessionExpired,
// 		Type:       "error",
// 	}

// 	ErrorSessionInvalid = ResponseCode{
// 		Code:       "ERROR_SESSION_INVALID",
// 		StatusCode: StatusBadRequest,
// 		Message:    MsgSessionInvalid,
// 		Type:       "error",
// 	}

// 	ErrorSessionCreationFailed = ResponseCode{
// 		Code:       "ERROR_SESSION_CREATION_FAILED",
// 		StatusCode: StatusInternalServerError,
// 		Message:    MsgSessionCreationFailed,
// 		Type:       "error",
// 	}

// 	ErrorSessionUpdateFailed = ResponseCode{
// 		Code:       "ERROR_SESSION_UPDATE_FAILED",
// 		StatusCode: StatusInternalServerError,
// 		Message:    MsgSessionUpdateFailed,
// 		Type:       "error",
// 	}

// 	ErrorFileTooLarge = ResponseCode{
// 		Code:       "ERROR_FILE_TOO_LARGE",
// 		StatusCode: StatusRequestEntityTooLarge,
// 		Message:    MsgFileTooLarge,
// 		Type:       "error",
// 	}

// 	ErrorFileInvalidType = ResponseCode{
// 		Code:       "ERROR_FILE_INVALID_TYPE",
// 		StatusCode: StatusBadRequest,
// 		Message:    MsgFileInvalidType,
// 		Type:       "error",
// 	}

// 	ErrorFileUploadFailed = ResponseCode{
// 		Code:       "ERROR_FILE_UPLOAD_FAILED",
// 		StatusCode: StatusInternalServerError,
// 		Message:    MsgFileUploadFailed,
// 		Type:       "error",
// 	}

// 	ErrorFileNotFound = ResponseCode{
// 		Code:       "ERROR_FILE_NOT_FOUND",
// 		StatusCode: StatusNotFound,
// 		Message:    MsgFileNotFound,
// 		Type:       "error",
// 	}

// 	ErrorFileDeleteFailed = ResponseCode{
// 		Code:       "ERROR_FILE_DELETE_FAILED",
// 		StatusCode: StatusInternalServerError,
// 		Message:    MsgFileDeleteFailed,
// 		Type:       "error",
// 	}

// 	ErrorValidationFailed = ResponseCode{
// 		Code:       "ERROR_VALIDATION_FAILED",
// 		StatusCode: StatusBadRequest,
// 		Message:    MsgValidationFailed,
// 		Type:       "error",
// 	}

// 	ErrorRequiredFieldMissing = ResponseCode{
// 		Code:       "ERROR_REQUIRED_FIELD_MISSING",
// 		StatusCode: StatusBadRequest,
// 		Message:    MsgRequiredFieldMissing,
// 		Type:       "error",
// 	}

// 	ErrorInvalidFormat = ResponseCode{
// 		Code:       "ERROR_INVALID_FORMAT",
// 		StatusCode: StatusBadRequest,
// 		Message:    MsgInvalidFormat,
// 		Type:       "error",
// 	}

// 	ErrorInvalidEmail = ResponseCode{
// 		Code:       "ERROR_INVALID_EMAIL",
// 		StatusCode: StatusBadRequest,
// 		Message:    MsgInvalidEmail,
// 		Type:       "error",
// 	}

// 	ErrorInvalidPhoneNumber = ResponseCode{
// 		Code:       "ERROR_INVALID_PHONE_NUMBER",
// 		StatusCode: StatusBadRequest,
// 		Message:    MsgInvalidPhoneNumber,
// 		Type:       "error",
// 	}

// 	ErrorInvalidDate = ResponseCode{
// 		Code:       "ERROR_INVALID_DATE",
// 		StatusCode: StatusBadRequest,
// 		Message:    MsgInvalidDate,
// 		Type:       "error",
// 	}

// 	ErrorInvalidID = ResponseCode{	
// 		Code:       "ERROR_INVALID_ID",
// 		StatusCode: StatusBadRequest,
// 		Message:    MsgInvalidID,
// 		Type:       "error",
// 	}

// 	ErrorFieldTooLong = ResponseCode{
// 		Code:       "ERROR_FIELD_TOO_LONG",
// 		StatusCode: StatusBadRequest,
// 		Message:    MsgFieldTooLong,
// 		Type:       "error",
// 	}

// 	ErrorFieldTooShort = ResponseCode{
// 		Code:       "ERROR_FIELD_TOO_SHORT",
// 		StatusCode: StatusBadRequest,
// 		Message:    MsgFieldTooShort,
// 		Type:       "error",
// 	}

// 	ErrorInternalServerError = ResponseCode{
// 		Code:       "ERROR_INTERNAL_SERVER_ERROR",
// 		StatusCode: StatusInternalServerError,
// 		Message:    MsgInternalServerError,
// 		Type:       "error",
// 	}

// 	ErrorServiceUnavailable = ResponseCode{
// 		Code:       "ERROR_SERVICE_UNAVAILABLE",
// 		StatusCode: StatusServiceUnavailable,
// 		Message:    MsgServiceUnavailable,
// 		Type:       "error",
// 	}

// 	ErrorDatabaseError = ResponseCode{
// 		Code:       "ERROR_DATABASE_ERROR",
// 		StatusCode: StatusInternalServerError,
// 		Message:    MsgDatabaseError,
// 		Type:       "error",
// 	}

// 	ErrorNetworkError = ResponseCode{
// 		Code:       "ERROR_NETWORK_ERROR",
// 		StatusCode: StatusInternalServerError,
// 		Message:    MsgNetworkError,
// 		Type:       "error",
// 	}

// 	ErrorTimeoutError = ResponseCode{
// 		Code:       "ERROR_TIMEOUT_ERROR",
// 		StatusCode: StatusRequestTimeout,
// 		Message:    MsgTimeoutError,
// 		Type:       "error",
// 	}

// 	ErrorUnexpectedError = ResponseCode{
// 		Code:       "ERROR_UNEXPECTED_ERROR",
// 		StatusCode: StatusInternalServerError,
// 		Message:    MsgUnexpectedError,
// 		Type:       "error",
// 	}

// 	ErrorConfigurationError = ResponseCode{
// 		Code:       "ERROR_CONFIGURATION_ERROR",
// 		StatusCode: StatusInternalServerError,
// 		Message:    MsgConfigurationError,
// 		Type:       "error",
// 	}

// 	ErrorExternalServiceError = ResponseCode{
// 		Code:       "ERROR_EXTERNAL_SERVICE_ERROR",
// 		StatusCode: StatusInternalServerError,
// 		Message:    MsgExternalServiceError,
// 		Type:       "error",
// 	}

// 	ErrorInsufficientBalance = ResponseCode{
// 		Code:       "ERROR_INSUFFICIENT_BALANCE",
// 		StatusCode: StatusBadRequest,
// 		Message:    MsgInsufficientBalance,
// 		Type:       "error",
// 	}

// 	ErrorTransactionFailed = ResponseCode{
// 		Code:       "ERROR_TRANSACTION_FAILED",
// 		StatusCode: StatusInternalServerError,
// 		Message:    MsgTransactionFailed,
// 		Type:       "error",
// 	}

// 	ErrorLimitExceeded = ResponseCode{
// 		Code:       "ERROR_LIMIT_EXCEEDED",
// 		StatusCode: StatusBadRequest,
// 		Message:    MsgLimitExceeded,
// 		Type:       "error",
// 	}

// 	ErrorOperationNotAllowed = ResponseCode{
// 		Code:       "ERROR_OPERATION_NOT_ALLOWED",
// 		StatusCode: StatusForbidden,
// 		Message:    MsgOperationNotAllowed,
// 		Type:       "error",
// 	}

// 	ErrorResourceBusy = ResponseCode{
// 		Code:       "ERROR_RESOURCE_BUSY",
// 		StatusCode: StatusConflict,
// 		Message:    MsgResourceBusy,
// 		Type:       "error",
// 	}

// 	ErrorMaintenanceMode = ResponseCode{
// 		Code:       "ERROR_MAINTENANCE_MODE",
// 		StatusCode: StatusServiceUnavailable,
// 		Message:    MsgMaintenanceMode,
// 		Type:       "error",
// 	}
// )
