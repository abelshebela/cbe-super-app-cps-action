package shared

import (
	"errors"
)

type ErrorDefinition struct {
	Code    string
	Message string
}

func (e ErrorDefinition) Error() string {
	return e.Message
}

var (
	ErrNotFound = errors.New("not found")
)

type ErrorGroup map[string]ErrorDefinition

type ErrorDefinitions struct {
	Account ErrorGroup
	OTP     ErrorGroup
}

var DefineError = ErrorDefinitions{
	Account: ErrorGroup{
		"PHONE_LOOKUP_FAILED": {
			Code:    "ACC_001",
			Message: "Failed to perform phone number lookup.",
		},
		"MOCK_DATA_FETCH_FAILED": {
			Code:    "ACC_002",
			Message: "Failed to fetch mock user data.",
		},
		"API_REQUEST_FAILED": {
			Code:    "ACC_003",
			Message: "Failed to communicate with external API.",
		},
	},
	OTP: ErrorGroup{
		"INVALID_OTP": {
			Code:    "OTP_001",
			Message: "The provided OTP is invalid",
		},
		"EXPIRED_OTP": {
			Code:    "OTP_002",
			Message: "The OTP has expired",
		},
		"EMAIL_IN_USE": {
			Code:    "OTP_003",
			Message: "Email address is already registered",
		},
	},
}