package common

type ErrorDefinition struct {
	Code    string
	Message string
}

type ErrorGroup map[string]ErrorDefinition

type ErrorDefinitions struct {
	General     ErrorGroup
	Auth        ErrorGroup
	User        ErrorGroup
	Transaction ErrorGroup
}

var DefineError = ErrorDefinitions{
	General: ErrorGroup{
		"CONFLICT_KEY": {
			Code:    "GEN_001",
			Message: "Duplicate key error: The specified field already exists.",
		},
		"INVALID_ID": {
			Code:    "GEN_002",
			Message: "Invalid ObjectId provided.",
		},
		"INVALID_JSON_PAYLOAD": {
			Code:    "GEN_003",
			Message: "Invalid JSON payload.",
		},
		"UNHANDLED_SERVER_ERROR": {
			Code:    "GEN_004",
			Message: "An unexpected server error occurred.",
		},
		"INVALID_TOKEN_FORMAT": {
			Code:    "GEN_005",
			Message: "Invalid token format.",
		},
		"INVALID_TOKEN": {
			Code:    "GEN_006",
			Message: "Invalid token.",
		},
		"EXPIRED_TOKEN": {
			Code:    "GEN_007",
			Message: "Token has expired.",
		},
		"ENCRYPTED_PAYLOAD_REQUIRED": {
			Code:    "GEN_008",
			Message: "Encrypted payload is required.",
		},
		"INVALID_SERVER_CONFIGURATION": {
			Code:    "GEN_009",
			Message: "Invalid server configuration",
		},
		"AUTH_HEADER_MISSING": {
			Code:    "GEN_010",
			Message: "Authorization header is missing. please provide a valid bearer token.",
		},
		"INVALID_KEY_CONFIGURATION": {
			Code:    "GEN_011",
			Message: "Invalid key configuration.",
		},
		"INVALID_ENCRYPTED_PAYLOAD": {
			Code:    "GEN_012",
			Message: "The provided encrypted payload is invalid.",
		},
		"DECRYPTION_ERROR": {
			Code:    "GEN_013",
			Message: "Decryption error occurred. please check your key configuration.",
		},
		"SERVER_KEYS_NOT_CONFIGURED": {
			Code:    "GEN_015",
			Message: "server error occurred. please try again later or contact support.",
		},
		"INVALID_INPUT": {
			Code:    "GEN_016",
			Message: "the provided input is invalid.",
		},
		"UNAUTHORIZED": {
			Code:    "GEN_016",
			Message: "missing authenticated user.",
		},
		"USER_REALM_NOT_FOUND": {
			Code:    "GEN_016",
			Message: "user is not permission to do this action.",
		},
		"ACTION_NOT_ALLOWED": {
			Code:    "GEN_017",
			Message: "user is not permission to do this action.",
		},
		// =========================

		"EMPTY_ORG": {
			Code:    "GEN_021",
			Message: "Cant respond for empty organization id",
		},
		"INVALID_REQ": {
			Code:    "GEN_022",
			Message: "Invalid Request body",
		},
		"REG_FRST": {
			Code:    "GEN_023",
			Message: "please register on the Super APP and visit your nearest branch for Verification Code",
		},
		"WEAK_PIN": {
			Code:    "GEN_024",
			Message: "weak PIN used: please use a strong combination",
		},
		"UNTRUSTED": {
			Code:    "GEN_025",
			Message: "untrusted installation: forbidden access",
		},
		"DEVICE_NOT_FOUND": {
			Code:    "GEN_026",
			Message: "Device not found",
		},
		"MISSING_HEADER": {
			Code:    "GEN_027",
			Message: "Missing Requeired Header",
		},
		"USED_PIN": {
			Code:    "GEN_028",
			Message: "The PIN has been used recently. Please choose a different PIN.",
		},
		"SAME_PIN": {
			Code:    "GEN_029",
			Message: "PIN has no change: please provide a different PIN",
		},
	},
	Auth: ErrorGroup{
		"AUTH_USER_NOT_FOUND": {
			Code:    "AUTH_001",
			Message: "User not found.",
		},
		"AUTH_USER_DISABLED": {
			Code:    "AUTH_002",
			Message: "User is not allowed to login, please contact your admin!",
		},
		"AUTH_INVALID_PASSWORD": {
			Code:    "AUTH_002",
			Message: "Invalid password.",
		},
		"AUTH_USER_HAS_NO_PASSWORD": {
			Code:    "AUTH_003",
			Message: "Please reset your password. to login",
		},
		"AUTH_TOO_MANY_ATTEMPTS": {
			Code:    "AUTH_004",
			Message: `Too many attempts. Try again after ${waiting_time} minutes`,
		},
		"AUTH_INVALID_OTP": {
			Code:    "AUTH_005",
			Message: "The verification Code you entered is incorrect. Please check the Code and try again.",
		},
		"AUTH_EXPIRED_OTP": {
			Code:    "AUTH_006",
			Message: "The verification Code has expired. Please request a new Code to continue.",
		},
		"AUTH_USER_RESET_PASSWORD_REQUIRED": {
			Code:    "AUTH_007",
			Message: "Too many login attempts. Please reset your password. Your account is locked until then.",
		},
		"AUTH_USER_ALREADY_EXISTS": {
			Code:    "AUTH_08",
			Message: `user with this ${name} already exist`,
		},
		"LOGIN_PROHIBITED_FOR_15_MIN": {
			Code:    "AUTH_010",
			Message: "Incorrect credentials. Login is prohibited for 15 minutes.",
		},
		"INCORRECT_PASSWORD": {
			Code:    "AUTH_011",
			Message: `Incorrect password. ${triesLeft} login attempt(s) left.`,
		},
		"OLD_PASSWORD_SAME_AS_NEW": {
			Code:    "AUTH_010",
			Message: "New password can't be the same as old password.",
		},
		// ========================
		"PIN_OLY_DIG": {
			Code:    "AUTH_011",
			Message: "PIN must contain only digits",
		},
		"PIN_LIMIT": {
			Code:    "AUTH_012",
			Message: "PIN must be exactly 6 digits",
		},
		"PIN_REDANDANT": {
			Code:    "AUTH_013",
			Message: "PIN cannot contain more than 4 redundant numbers",
		},
		"PIN_SEQ": {
			Code:    "AUTH_014",
			Message: "PIN cannot contain 4 or more sequential numbers",
		},
		"FAILD_TO_GEN_TOKEN": {
			Code:    "AUTH_015",
			Message: "Set your new password",
		},
		"FAILD_TO_RESET_PASS": {
			Code:    "AUTH_016",
			Message: "Failed to set your new password",
		},
		"FAILD_VALIDATION": {
			Code:    "AUTH_017",
			Message: "Validation failed",
		},
		"INVALID_BEARER": {
			Code:    "AUTH_018",
			Message: "Unauthorized: Invalid Bearer token format",
		},
		"USE_RIGHT_AUTH": {
			Code:    "AUTH_019",
			Message: "Use the Right Authentication",
		},
		"INVALID_CLAIM": {
			Code:    "AUTH_020",
			Message: "Invalid token claims",
		},
		"INVALID_TOKEN_DATA": {
			Code:    "AUTH_021",
			Message: "Invalid token data",
		},
		"UNABLE_TO_DYCRYPT_TOKEN": {
			Code:    "AUTH_022",
			Message: "Unable to decrypt token",
		},
		"FAILED_LOGIN": {
			Code:    "AUTH_023",
			Message: "Failed to login. you entered wrong password",
		},
		"ACCOUNT_LOCKEC": {
			Code:    "AUTH_024",
			Message: "account locked. please contact your addministrator",
		},
		"OTP_CHECK_FAILED": {
			Code:    "AUTH_025",
			Message: "OTP check failed",
		},
		"REG_ON_SUP": {
			Code:    "AUTH_026",
			Message: "please register on the Super APP and visit your nearest branch for Verification Code",
		},
		"VISIT_NEAREST": {
			Code:    "AUTH_027",
			Message: "please visit your nearest branch for Verification Code",
		},
	},
	Transaction: ErrorGroup{
		"TRANSACTION_NOT_FOUND": {
			Code:    "TXN_001",
			Message: "Transaction not found.",
		},
		"INSUFFICIENT_FUNDS": {
			Code:    "TXN_002",
			Message: "Insufficient balance.",
		},
		"COMMISSION_NOT_FOUND": {
			Code:    "COM_001",
			Message: "Commission not found",
		},
		"CUSTOMER_NOT_FOUND": {
			Code:    "CUS_001",
			Message: "Customer not found",
		},
	},
}
