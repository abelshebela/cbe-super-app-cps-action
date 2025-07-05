package users

import "errors"

var (
	ErrNotFound = errors.New("not found")

	// Registration errors
	ErrPhoneAlreadyExists        = errors.New("PHONE_ALREADY_EXISTS")
	ErrDeviceAlreadyRegistered   = errors.New("DEVICE_ALREADY_REGISTERED")
	ErrInvalidPhoneNumber        = errors.New("INVALID_PHONE_NUMBER")
	ErrInvalidPlatform           = errors.New("INVALID_PLATFORM")
	ErrInvalidDeviceUUID         = errors.New("INVALID_DEVICE_UUID")
	ErrRegistrationInProgress    = errors.New("REGISTRATION_IN_PROGRESS")
	ErrRegistrationRateLimited   = errors.New("REGISTRATION_RATE_LIMITED")
	ErrRegistrationFailed        = errors.New("REGISTRATION_FAILED")
	ErrRegistrationNotFound      = errors.New("REGISTRATION_NOT_FOUND")
	ErrRegistrationExpired       = errors.New("REGISTRATION_EXPIRED")
	ErrInvalidRegistrationStatus = errors.New("INVALID_REGISTRATION_STATUS")
	ErrUserCreationFailed        = errors.New("USER_CREATION_FAILED")

	// Login errors
	ErrUserNotFound          = errors.New("USER_NOT_FOUND")
	ErrInvalidPin            = errors.New("INVALID_PIN")
	ErrAccountBlocked        = errors.New("ACCOUNT_BLOCKED")
	ErrAccountDeleted        = errors.New("ACCOUNT_DELETED")
	ErrDeviceNotLinked       = errors.New("DEVICE_NOT_LINKED")
	ErrTooManyLoginAttempts  = errors.New("TOO_MANY_LOGIN_ATTEMPTS")
	ErrPinNotSet             = errors.New("PIN_NOT_SET")
	ErrTokenGenerationFailed = errors.New("TOKEN_GENERATION_FAILED")

	// PIN Reset errors
	ErrPinResetSessionNotFound   = errors.New("PIN_RESET_SESSION_NOT_FOUND")
	ErrPinResetSessionExpired    = errors.New("PIN_RESET_SESSION_EXPIRED")
	ErrPinResetSessionInvalid    = errors.New("PIN_RESET_SESSION_INVALID")
	ErrPinResetTooManyAttempts   = errors.New("PIN_RESET_TOO_MANY_ATTEMPTS")
	ErrPinResetRateLimited       = errors.New("PIN_RESET_RATE_LIMITED")
	ErrPinResetFailed            = errors.New("PIN_RESET_FAILED")
	ErrPinResetOTPInvalid        = errors.New("PIN_RESET_OTP_INVALID")
	ErrPinResetOTPExpired        = errors.New("PIN_RESET_OTP_EXPIRED")
	ErrPinResetUserNotFound      = errors.New("PIN_RESET_USER_NOT_FOUND")
	ErrPinResetDeviceMismatch    = errors.New("PIN_RESET_DEVICE_MISMATCH")
	ErrPinResetAlreadyInProgress = errors.New("PIN_RESET_ALREADY_IN_PROGRESS")
)
