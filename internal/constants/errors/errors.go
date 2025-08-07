package errors

import (
	"errors"
	"net/http"
)

var (
	ErrAccBlocked                   = errors.New("account is blocked")
	ErrOldDevice                    = errors.New("your device app is older version")
	ErrIdEmpty                      = errors.New("id can not be empty")
	ErrUnexpected                   = errors.New("unexpected error")
	ErrInternalServerError          = errors.New("internal server error")
	ErrRequestTimeout               = errors.New("request timeout")
	ErrAccountNotFound              = errors.New("account not found")
	ErrBadRequest                   = errors.New("bad request")
	ErrInvalidData                  = errors.New("invalid data")
	ErrUnauthorized                 = errors.New("unauthorized")
	ErrActionNotAllowed             = errors.New("action not allowed")
	ErrMoneyRequestNotFound         = errors.New("money request not found")
	ErrDeviceMismatch               = errors.New("device mismatch")
	ErrTimeout                      = errors.New("request timeout")
	ErrDonationNotFound             = errors.New("donation not found")
	ErrUserNotFound                 = errors.New("user not found")
	ErrOTPNotFound                  = errors.New("otp not found")
	ErrFileTooLarge                 = errors.New("file too large")
	ErrInvalidFileType              = errors.New("invalid file type")
	ErrBranchNotFound               = errors.New("branch not found")
	ErrPhoneNumberAlreadyExists     = errors.New("phone number already exists")
	ErrRequestFailed                = errors.New("request failed")
	ErrUserAccountBlocked           = errors.New("account is blocked")
	ErrDeviceNotFound               = errors.New("device not found")
	ErrDeviceIDRequired             = errors.New("device id is required")
	ErrInvalidFormData              = errors.New("invalid form data")
	ErrHQDataNotFound               = errors.New("hq data not found")
	ErrPINMismatch                  = errors.New("old pin mismatch")
	ErrInvalidPIN                   = errors.New("invalid pin")
	ErrPINOnlyDigit                 = errors.New("pin only digit")
	ErrWeakPIN                      = errors.New("weak pin")
	ErrPINRedundant                 = errors.New("pin redundant")
	ErrPINSeq                       = errors.New("pin seq")
	ErrPinInHistory                 = errors.New("the pin you use is recently used")
	ErrSamePIN                      = errors.New("same pin")
	ErrPINLength                    = errors.New(" PIN must be exactly 6 digits")
	ErrOtpExpired                   = errors.New("otp expired")
	ErrInvalidOTP                   = errors.New("invalid otp")
	ErrRegistrationExpired          = errors.New("registration expired")
	ErrInvalidRegistrationStatus    = errors.New("invalid registration status")
	ErrSessionNotFound              = errors.New("session not found")
	ErrPinResetAlreadyInProgress    = errors.New("pin reset already in progress")
	ErrUserAlreadyEmail             = errors.New("user already has email")
	ErrEmailInuse                   = errors.New("email in use")
	ErrWaitForPreviousOTPExpiration = errors.New("wait for previous otp expiration")
	ErrTooManyLoginAttempts         = errors.New("too many login attempts")
	ErrUserAlreadyExists            = errors.New("user already exists")
	ErrDeviceAleadyExists           = errors.New("device already exists")
	ErrRegistrationInProgress       = errors.New("registration already in progress")
	ErrSessionExpired               = errors.New("session expired")
	ErrInvalidSession               = errors.New("invalid session")
	ErrTooManyResetAttempts         = errors.New("too many reset attempts")
	ErrNoLinkedDevices              = errors.New("no linked devices")
	ErrSMSSendFailure               = errors.New("sms send failure")
	ErrInvalidPassword              = errors.New("invalid password")
	ErrEmptyFilterParam             = errors.New("empty filter param")
	ErrTryToSaveEmptyUser           = errors.New("error saving empty user data")
	ErrNoMongoDocument              = errors.New("mongo: no documents in result")
	ErrPhoneNumberCanNotBeEmpty     = errors.New("phone number can not be empty")
	ErrDeviceUUIDCanNotBeNull       = errors.New("device uuid can not be null")

	ErrUserBlockedByCps          = errors.New("user is blocked by cps")
	ErrUserAccountBlockedByCps   = errors.New("user account is blocked")
	ErrInvalidKey                = errors.New("invalid key or IV size")
	ErrInvalidEncData            = errors.New("invalid encrypted data length")
	ErrInvalidPadding            = errors.New("invalid padding")
	ErrFailedOtpCreation         = errors.New("failed to create OTP record")
	ErrOTPAlreadyExist           = errors.New("otp Already exists wait until it expires")
	ErrFailedToPrepareOTP        = errors.New("failed to prepare SMS payload: %w")
	ErrFailedSMSApiCall          = errors.New("SMS API call failed: %w")
	ErrFailedHttpCall            = errors.New("failed to create HTTP request: %w")
	ErrPhoneNotFound             = errors.New("phone number not found")
	ErrRegistrationFailedExpired = errors.New("registration failed")
	ErrPinResetDeviceMismatch    = errors.New("pin reset device mismatch")
	ErrPinResetFailed            = errors.New("pin reset failed")

	ErrPinNeedActivation          = errors.New("pin reset session %s is not verified please visit nearest branch")
	ErrPinResetSessionInvalid     = errors.New("invalid pin reset session")
	ErrPinResetSessionNotFound    = errors.New("pin reset session not found")
	ErrPinResetSessionExpired     = errors.New("pin reset session Expired")
	ErrOldPinMismatch             = errors.New("old Pin Mismatch")
	ErrFailedToCreateTemp         = errors.New("failed to create temp file")
	ErrFailedToUpload             = errors.New("failed to upload file")
	ErrProfileSet                 = errors.New("failed To Set Profile on user")
	ErrInvalidDateFormat          = errors.New("invalid date format sent from the client")
	ErrDeviceDiffInstallationDate = errors.New("the divice have different installation date")
	ErrNoResetSession             = errors.New("no reset session found")
	ErrDeviceUUidExist            = errors.New("device uuid already exist")
	ErrDeviceDataExist            = errors.New("data already exist")
	ErrEmptyEmptyData             = errors.New("trying to update empty data")
	ErrTooManyAttempt             = errors.New("too many login attempts. Please try again later")
	ErrGotoBranchToEnable         = errors.New("you have to go branch to enable this action")
)

var ErrorMap = map[error]int{
	ErrGotoBranchToEnable:           http.StatusBadRequest,
	ErrTooManyAttempt:               http.StatusBadRequest,
	ErrEmptyEmptyData:               http.StatusBadRequest,
	ErrDeviceDataExist:              http.StatusBadRequest,
	ErrNoResetSession:               http.StatusNoContent,
	ErrProfileSet:                   http.StatusInternalServerError,
	ErrFailedToUpload:               http.StatusInternalServerError,
	ErrFailedToCreateTemp:           http.StatusInternalServerError,
	ErrDeviceUUidExist:              http.StatusBadRequest,
	ErrOldPinMismatch:               http.StatusBadRequest,
	ErrPinResetSessionExpired:       http.StatusBadRequest,
	ErrPinResetSessionNotFound:      http.StatusBadRequest,
	ErrPinResetSessionInvalid:       http.StatusBadRequest,
	ErrPinNeedActivation:            http.StatusBadRequest,
	ErrPhoneNotFound:                http.StatusBadRequest,
	ErrPinResetFailed:               http.StatusBadRequest,
	ErrPinResetDeviceMismatch:       http.StatusBadRequest,
	ErrUserAccountBlockedByCps:      http.StatusBadRequest,
	ErrFailedOtpCreation:            http.StatusBadRequest,
	ErrOTPAlreadyExist:              http.StatusBadRequest,
	ErrInvalidPadding:               http.StatusBadRequest,
	ErrInvalidEncData:               http.StatusBadRequest,
	ErrInvalidKey:                   http.StatusBadRequest,
	ErrUserBlockedByCps:             http.StatusBadRequest,
	ErrDeviceDiffInstallationDate:   http.StatusBadRequest,
	ErrFailedHttpCall:               http.StatusInternalServerError,
	ErrFailedSMSApiCall:             http.StatusInternalServerError,
	ErrFailedToPrepareOTP:           http.StatusInternalServerError,
	ErrUnexpected:                   http.StatusInternalServerError,
	ErrInternalServerError:          http.StatusInternalServerError,
	ErrRequestTimeout:               http.StatusRequestTimeout,
	ErrAccountNotFound:              http.StatusNotFound,
	ErrNoMongoDocument:              http.StatusNotFound,
	ErrAccBlocked:                   http.StatusBadRequest,
	ErrOldDevice:                    http.StatusBadRequest,
	ErrDeviceUUIDCanNotBeNull:       http.StatusBadRequest,
	ErrIdEmpty:                      http.StatusBadRequest,
	ErrBadRequest:                   http.StatusBadRequest,
	ErrPhoneNumberCanNotBeEmpty:     http.StatusBadRequest,
	ErrTryToSaveEmptyUser:           http.StatusBadRequest,
	ErrInvalidData:                  http.StatusBadRequest,
	ErrUnauthorized:                 http.StatusUnauthorized,
	ErrActionNotAllowed:             http.StatusForbidden,
	ErrMoneyRequestNotFound:         http.StatusNotFound,
	ErrDeviceMismatch:               http.StatusConflict,
	ErrTimeout:                      http.StatusGatewayTimeout,
	ErrDonationNotFound:             http.StatusNotFound,
	ErrUserNotFound:                 http.StatusNotFound,
	ErrOTPNotFound:                  http.StatusNotFound,
	ErrFileTooLarge:                 http.StatusRequestEntityTooLarge,
	ErrInvalidFileType:              http.StatusBadRequest,
	ErrBranchNotFound:               http.StatusNotFound,
	ErrPhoneNumberAlreadyExists:     http.StatusBadRequest,
	ErrRequestFailed:                http.StatusBadRequest,
	ErrUserAccountBlocked:           http.StatusBadRequest,
	ErrDeviceNotFound:               http.StatusNotFound,
	ErrDeviceIDRequired:             http.StatusBadRequest,
	ErrHQDataNotFound:               http.StatusNotFound,
	ErrPINMismatch:                  http.StatusBadRequest,
	ErrInvalidPIN:                   http.StatusBadRequest,
	ErrPINOnlyDigit:                 http.StatusBadRequest,
	ErrWeakPIN:                      http.StatusBadRequest,
	ErrPINSeq:                       http.StatusBadRequest,
	ErrPinInHistory:                 http.StatusBadRequest,
	ErrSamePIN:                      http.StatusBadRequest,
	ErrPINLength:                    http.StatusBadRequest,
	ErrOtpExpired:                   http.StatusBadRequest,
	ErrInvalidOTP:                   http.StatusBadRequest,
	ErrRegistrationExpired:          http.StatusBadRequest,
	ErrRegistrationFailedExpired:    http.StatusBadRequest,
	ErrSessionNotFound:              http.StatusNotFound,
	ErrPinResetAlreadyInProgress:    http.StatusBadRequest,
	ErrUserAlreadyEmail:             http.StatusBadRequest,
	ErrEmailInuse:                   http.StatusBadRequest,
	ErrWaitForPreviousOTPExpiration: http.StatusBadRequest,
	ErrTooManyLoginAttempts:         http.StatusBadRequest,
	ErrUserAlreadyExists:            http.StatusBadRequest,
	ErrDeviceAleadyExists:           http.StatusBadRequest,
	ErrRegistrationInProgress:       http.StatusBadRequest,
	ErrSessionExpired:               http.StatusBadRequest,
	ErrInvalidSession:               http.StatusBadRequest,
	ErrTooManyResetAttempts:         http.StatusBadRequest,
	ErrNoLinkedDevices:              http.StatusBadRequest,
	ErrInvalidPassword:              http.StatusBadRequest,
	ErrEmptyFilterParam:             http.StatusBadRequest,
}
