package common

type SuccessDefinition struct {
	Code    string
	Message string
}

type SuccessGroup map[string]SuccessDefinition

type SuccessDefinitions struct {
	General     SuccessGroup
	Auth        SuccessGroup
	User        SuccessGroup
	Transaction SuccessGroup
}

var DefineSuccess = SuccessDefinitions{
	General: SuccessGroup{
		"CONFLICT_KEY": {
			Code:    "GEN_001",
			Message: "Duplicate key success: The specified field already exists.",
		},
	},
	Auth: SuccessGroup{
		"AUTH_USER_OTP_SENT_BPS": {
			Code:    "AUTH_SUCCESS_001",
			Message: "Hello BPS user your are logged in for the first time it needs to verified the OTP is sent to your phone",
		},
		"AUTH_USER_BPS_OTP_VERIFIED": {
			Code:    "AUTH_SUCCESS_002",
			Message: "Welcome back to the Supper App BPS",
		},
		"SET_NEW_PASS": {
			Code:    "AUTH_014",
			Message: "Set your new password",
		},
		"NEW_PASS_SETS": {
			Code:    "AUTH_014",
			Message: "Password successfuly changed",
		},
		"DEVICE_FOUND": {
			Code:    "AUTH_015",
			Message: "Your Device is already registered, please login with your pin",
		},
		"DEVICE_NOT_FOUND": {
			Code:    "AUTH_016",
			Message: "Your Device is not registered, please register your device",
		},
		"AUTH_USER_OTP_SENT_USER": {
			Code:    "AUTH_SUCCESS_017",
			Message: "Hello user OTP is sent to your phone",
		},
	},
	Transaction: SuccessGroup{
		"TRANSACTION_NOT_FOUND": {
			Code:    "TXN_001",
			Message: "Transaction not found.",
		},
	},
}
