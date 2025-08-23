package localization

// Success Messages
const (
	// User related success messages
	MsgUserCreatedSuccessfully   = "User created successfully"
	MsgUserUpdatedSuccessfully   = "User updated successfully"
	MsgUserDeletedSuccessfully   = "User deleted successfully"
	MsgUserRetrievedSuccessfully = "User retrieved successfully"
	MsgUserLoginSuccessfully     = "User logged in successfully"
	MsgUserLogoutSuccessfully    = "User logged out successfully"
	MsgUserProfileUpdated        = "User profile updated successfully"
	MsgUserPasswordChanged       = "Password changed successfully"
	MsgUserPINChanged            = "PIN changed successfully"
	MsgUserDeviceLinked          = "Device linked successfully"
	MsgUserDeviceUnlinked        = "Device unlinked successfully"

	// OTP related success messages
	MsgOTPSentSuccessfully     = "OTP sent successfully"
	MsgOTPVerifiedSuccessfully = "OTP verified successfully"
	MsgOTPResentSuccessfully   = "OTP resent successfully"

	// Session related success messages
	MsgSessionCreatedSuccessfully    = "Session created successfully"
	MsgSessionRefreshedSuccessfully  = "Session refreshed successfully"
	MsgSessionTerminatedSuccessfully = "Session terminated successfully"

	// File related success messages
	MsgFileUploadedSuccessfully  = "File uploaded successfully"
	MsgFileDeletedSuccessfully   = "File deleted successfully"
	MsgFileRetrievedSuccessfully = "File retrieved successfully"

	// General success messages
	MsgOperationCompleted        = "Operation completed successfully"
	MsgDataRetrievedSuccessfully = "Data retrieved successfully"
	MsgDataSavedSuccessfully     = "Data saved successfully"
	MsgDataUpdatedSuccessfully   = "Data updated successfully"
	MsgDataDeletedSuccessfully   = "Data deleted successfully"
	MsgValidationPassed          = "Validation passed successfully"
	MsgHealthCheckPassed         = "Health check passed"
	MsgBadRequest                = "Bad request"
)

// Error Messages
const (
	// User related error messages
	MsgUserNotFound             = "User not found"
	MsgUserAlreadyExists        = "User already exists"
	MsgUserAccountBlocked       = "User account is blocked"
	MsgUserUnauthorized         = "User is not authorized"
	MsgUserForbidden            = "User access is forbidden"
	MsgUserInvalidCredentials   = "Invalid credentials"
	MsgUserSessionExpired       = "User session has expired"
	MsgUserTooManyLoginAttempts = "Too many login attempts"
	MsgUserDeviceMismatch       = "Device mismatch detected"
	MsgUserDeviceNotFound       = "User device not found"
	MsgUserDeviceAlreadyExists  = "User device already exists"
	MsgSamePIN                  = "You used the same pin as the old one"
	// OTP related error messages
	MsgOTPNotFound        = "OTP not found"
	MsgOTPExpired         = "OTP has expired"
	MsgOTPInvalid         = "Invalid OTP"
	MsgOTPAlreadyExists   = "OTP already exists, wait until it expires"
	MsgOTPTooManyAttempts = "Too many OTP attempts"
	MsgOTPSendFailed      = "Failed to send OTP"

	// PIN related error messages
	MsgPINInvalid       = "Invalid PIN"
	MsgPINMismatch      = "PIN mismatch"
	MsgPINTooWeak       = "PIN is too weak"
	MsgPINInHistory     = "PIN was recently used"
	MsgPINLengthInvalid = "PIN must be exactly 6 digits"
	MsgPINOnlyDigits    = "PIN must contain only digits"
	MsgPINSequential    = "PIN cannot be sequential"
	MsgPINRedundant     = "PIN cannot be redundant"

	// Session related error messages
	MsgSessionNotFound       = "Session not found"
	MsgSessionExpired        = "Session has expired"
	MsgSessionInvalid        = "Invalid session"
	MsgSessionCreationFailed = "Failed to create session"
	MsgSessionUpdateFailed   = "Failed to update session"

	// File related error messages
	MsgFileTooLarge     = "File size is too large"
	MsgFileInvalidType  = "Invalid file type"
	MsgFileUploadFailed = "File upload failed"
	MsgFileNotFound     = "File not found"
	MsgFileDeleteFailed = "File deletion failed"

	// Validation error messages
	MsgValidationFailed     = "Validation failed"
	MsgRequiredFieldMissing = "Required field is missing"
	MsgInvalidFormat        = "Invalid format"
	MsgInvalidEmail         = "Invalid email format"
	MsgInvalidPhoneNumber   = "Invalid phone number format"
	MsgInvalidDate          = "Invalid date format"
	MsgInvalidID			= "Invalid ID format"
	MsgFieldTooLong         = "Field value is too long"
	MsgFieldTooShort        = "Field value is too short"

	// System error messages
	MsgInternalServerError  = "Internal server error occurred"
	MsgServiceUnavailable   = "Service is temporarily unavailable"
	MsgDatabaseError        = "Database operation failed"
	MsgNetworkError         = "Network error occurred"
	MsgTimeoutError         = "Request timeout occurred"
	MsgUnexpectedError      = "An unexpected error occurred"
	MsgConfigurationError   = "Configuration error"
	MsgExternalServiceError = "External service error"

	// Business logic error messages
	MsgInsufficientBalance = "Insufficient balance"
	MsgTransactionFailed   = "Transaction failed"
	MsgLimitExceeded       = "Limit exceeded"
	MsgOperationNotAllowed = "Operation not allowed"
	MsgResourceBusy        = "Resource is busy"
	MsgMaintenanceMode     = "System is in maintenance mode"
)
