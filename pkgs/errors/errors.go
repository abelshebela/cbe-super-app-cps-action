package errors

import (
	"fmt"
)

// ServiceError defines a structured error with a code, message, and optional cause.
type ServiceError struct {
	Code    string
	Message string
	Cause   error
}

// Error implements the error interface.
func (e *ServiceError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s | cause: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying cause for errors.Is and errors.As.
func (e *ServiceError) Unwrap() error {
	return e.Cause
}

// New creates a new ServiceError without an underlying cause.
func New(code, message string) *ServiceError {
	return &ServiceError{
		Code:    code,
		Message: message,
	}
}

// Wrap creates a new ServiceError wrapping an underlying error.
func Wrap(code, message string, cause error) *ServiceError {
	return &ServiceError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// Predefined error constants.
var (
	ErrNotFound                = New("NOT_FOUND", "resource not found")
	ErrUploadFailed            = New("UPLOAD_FAILED", "upload failed")
	ErrUnhandledServerError    = New("UNHANDLED_SERVER_ERROR", "unhandled server error")
	ErrInvalidInput            = New("INVALID_INPUT", "invalid input provided")
	ErrInvalidPhoneNumber      = New("INVALID_PHONE_NUMBER", "invalid phone number")
	ErrInvalidDeviceUUID       = New("INVALID_DEVICE_UUID", "invalid device uuid")
	ErrInvalidPlatform         = New("INVALID_PLATFORM", "invalid platform")
	ErrInvalidPin              = New("INVALID_PIN", "invalid pin")
	ErrInvalidOTP              = New("INVALID_OTP", "invalid otp")
	ErrPinLimit                = New("PIN_LIMIT", "pin attempt limit reached")
	ErrPinOnlyDigit            = New("PIN_ONLY_DIGIT", "pin must contain only digits")
	ErrPinRedundant            = New("PIN_REDUNDANT", "pin is redundant")
	ErrPinSeq                  = New("PIN_SEQ", "pin has sequential digits")
	ErrPinInHistory            = New("PIN_IN_HISTORY", "pin was used previously")
	ErrOldPinMismatch          = New("OLD_PIN_MISMATCH", "old pin mismatch")
	ErrSamePin                 = New("SAME_PIN", "new pin is same as old pin")
	ErrTokenGenerationFailed   = New("TOKEN_GENERATION_FAILED", "token generation failed")
	ErrOtpCreationFailed       = New("OTP_CREATION_FAILED", "otp creation failed")
	ErrOtpEncryptionFailed     = New("OTP_ENCRYPTION_FAILED", "otp encryption failed")
	ErrOtpExpired              = New("EXPIRED_OTP", "otp expired")
	ErrOtpNotFound             = New("OTP_NOT_FOUND", "otp not found")
	ErrOtpInvalid              = New("INVALID_OTP", "invalid otp")
	ErrOtpWaitPrevious         = New("WAIT_FOR_PREVIOUS_OTP_EXPIRATION", "wait for previous otp expiration")
	ErrEmailInUse              = New("EMAIL_IN_USE", "email already in use")
	ErrUserAlreadyHasEmail     = New("USER_ALREADY_HAS_EMAIL", "user already has an email")
	ErrRegistrationInProgress  = New("REGISTRATION_IN_PROGRESS", "registration in progress")
	ErrRegistrationFailed      = New("REGISTRATION_FAILED", "registration failed")
	ErrDeviceAlreadyRegistered = New("DEVICE_ALREADY_REGISTERED", "device already registered")
	ErrPhoneAlreadyExists      = New("PHONE_ALREADY_EXISTS", "phone already exists")
	ErrAccountBlocked          = New("ACCOUNT_BLOCKED", "account is blocked")
	ErrTooManyLoginAttempts    = New("TOO_MANY_LOGIN_ATTEMPTS", "too many login attempts")
	ErrPinNotSet               = New("PIN_NOT_SET", "pin not set")
	ErrUserNotFound            = New("USER_NOT_FOUND", "user not found")
)
