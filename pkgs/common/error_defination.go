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
	Account     ErrorGroup
	OTP         ErrorGroup
	File        ErrorGroup
	Branch      ErrorGroup
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
		"ERROR_CHANGING_PIN": {
			Code:    "GEN_030",
			Message: "Error changing PIN.",
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
			Code:    "GEN_017",
			Message: "missing authenticated user.",
		},
		"USER_REALM_NOT_FOUND": {
			Code:    "GEN_018",
			Message: "user is not permission to do this action.",
		},
		"ACTION_NOT_ALLOWED": {
			Code:    "GEN_019",
			Message: "user is not permission to do this action.",
		},
		"WAIT_FOR_PREVIOUS_OTP_EXPIRATION": {
			Code:    "GEN_025",
			Message: "wait until the pervious otp expired",
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
		"NOT_FOUND": {
			Code:    "GEN_020",
			Message: "Resource not found.",
		},
		"DEVICE_ID_REQUIRED": {
			Code:    "GEN_026",
			Message: "Device UUID is required",
		},
		"NO_LINKED_DEVICES": {
			Code:    "GEN_027",
			Message: "User has no linked devices",
		},
		"COULD_NOT_UNLINK_DEVICE": {
			Code:    "GEN_028",
			Message: "Could not unlink device",
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
			Code:    "AUTH_003",
			Message: "Invalid password.",
		},
		"AUTH_USER_HAS_NO_PASSWORD": {
			Code:    "AUTH_004",
			Message: "Please reset your password. to login",
		},
		"AUTH_TOO_MANY_ATTEMPTS": {
			Code:    "AUTH_005",
			Message: `Too many attempts. Try again after ${waiting_time} minutes`,
		},
		"AUTH_INVALID_OTP": {
			Code:    "AUTH_006",
			Message: "The verification Code you entered is incorrect. Please check the Code and try again.",
		},
		"AUTH_EXPIRED_OTP": {
			Code:    "AUTH_007",
			Message: "The verification Code has expired. Please request a new Code to continue.",
		},
		"AUTH_USER_RESET_PASSWORD_REQUIRED": {
			Code:    "AUTH_008",
			Message: "Too many login attempts. Please reset your password. Your account is locked until then.",
		},
		"AUTH_USER_ALREADY_EXISTS": {
			Code:    "AUTH_009",
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
			Code:    "AUTH_012",
			Message: "New password can't be the same as old password.",
		},
		// ========================
		"PIN_OLY_DIG": {
			Code:    "AUTH_013",
			Message: "PIN must contain only digits",
		},
		"PIN_LIMIT": {
			Code:    "AUTH_014",
			Message: "PIN must be exactly 6 digits",
		},
		"PIN_REDANDANT": {
			Code:    "AUTH_015",
			Message: "PIN cannot contain more than 4 redundant numbers",
		},
		"PIN_SEQ": {
			Code:    "AUTH_016",
			Message: "PIN cannot contain 4 or more sequential numbers",
		},
		"FAILD_TO_GEN_TOKEN": {
			Code:    "AUTH_017",
			Message: "Set your new password",
		},
		"FAILD_TO_RESET_PASS": {
			Code:    "AUTH_018",
			Message: "Failed to set your new password",
		},
		"FAILD_VALIDATION": {
			Code:    "AUTH_019",
			Message: "Validation failed",
		},
		"OLD_PIN_MISMATCH": {
			Code:    "AUTH_029",
			Message: "Old PIN does not match",
		},
		"SAME_PIN": {
			Code:    "AUTH_027",
			Message: "New PIN cannot be the same as old PIN",
		},
		"PIN_IN_HISTORY": {
			Code:    "AUTH_028",
			Message: "New PIN cannot be one of the last 6 PINs used",
		},
		"INVALID_BEARER": {
			Code:    "AUTH_020",
			Message: "Unauthorized: Invalid Bearer token format",
		},
		"USE_RIGHT_AUTH": {
			Code:    "AUTH_021",
			Message: "Use the Right Authentication",
		},
		"INVALID_CLAIM": {
			Code:    "AUTH_022",
			Message: "Invalid token claims",
		},
		"INVALID_TOKEN_DATA": {
			Code:    "AUTH_023",
			Message: "Invalid token data",
		},
		"UNABLE_TO_DYCRYPT_TOKEN": {
			Code:    "AUTH_024",
			Message: "Unable to decrypt token",
		},
		"FAILED_LOGIN": {
			Code:    "AUTH_025",
			Message: "Failed to login. you entered wrong password",
		},
		"ACCOUNT_LOCKED": {
			Code:    "AUTH_026",
			Message: "account locked. please contact your addministrator",
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
	Account: ErrorGroup{
		"PHONE_LOOKUP_FAILED": {
			Code:    "ACC_001",
			Message: "Failed to lookup phone number.",
		},
		"API_REQUEST_FAILED": {
			Code:    "ACC_002",
			Message: "API request failed.",
		},
		"ACCOUNT_NOT_FOUND": {
			Code:    "ACC_003",
			Message: "Account not found.",
		},
		"USER_PHONE_EXISTS": {
			Code:    "ACC_004",
			Message: "User phone number already exists.",
		},
		"GENERAL_SIF_GENERATION_FAILED": {
			Code:    "ACC_005",
			Message: "Failed to generate SIF and account number.",
		},
		"GENERAL_DB_UPDATE_FAILED": {
			Code:    "ACC_006",
			Message: "Failed to update database record.",
		},
		"GENERAL_DB_INSERT_FAILED": {
			Code:    "ACC_007",
			Message: "Failed to insert database record.",
		},
		"CORE_ACCOUNT_NOT_FOUND": {
			Code:    "ACC_008",
			Message: "Core account not found.",
		},
		"CORE_ACCOUNT_TYPE_ERROR": {
			Code:    "ACC_009",
			Message: "Invalid core account detail type.",
		},
		"ACCOUNT_LOOKUP_FAILED": {
			Code:    "ACC_010",
			Message: "Account lookup failed.",
		},
		"ACCOUNT_DETAIL_FAILED": {
			Code:    "ACC_011",
			Message: "Account detail lookup failed.",
		},
		"INVALID_CORE_RESPONSE": {
			Code:    "ACC_012",
			Message: "Invalid core response.",
		},
		"ACCOUNT_ALREADY_LINKED": {
			Code:    "ACC_013",
			Message: "Account is already linked to another user.",
		},
	},
	OTP: ErrorGroup{
		"INVALID_OTP": {
			Code:    "OTP_001",
			Message: "Invalid OTP provided.",
		},
		"EXPIRED_OTP": {
			Code:    "OTP_002",
			Message: "OTP has expired.",
		},
		"EMAIL_IN_USE": {
			Code:    "OTP_003",
			Message: "Email is already in use.",
		},
		"USER_ALREADY_HAS_EMAIL": {
			Code:    "USER_002",
			Message: "User already has an email address",
		},
		"WAIT_FOR_PREVIOUS_OTP_EXPIRATION": {
			Code:    "USER_003",
			Message: "Please wait until the previous OTP expires",
		},
		"OTP_GENERATION_FAILED": {
			Code:    "OTP_004",
			Message: "Failed to generate OTP.",
		},
		"USER_KYC_LEVEL_ZERO": {
			Code:    "KYC_005",
			Message: "User KYC level is not sufficient for this operation",
		},
		"USER_KYC_LEVEL_WRONG": {
			Code:    "KYC_006",
			Message: "User KYC level is not sufficient for this operation",
		},
	},
	File: ErrorGroup{
		"INVALID_FORM": {
			Code:    "FILE_001",
			Message: "Could not parse form or file too large.",
		},
		"NO_FILE": {
			Code:    "FILE_002",
			Message: "File is required.",
		},
		"FILE_TOO_LARGE": {
			Code:    "FILE_003",
			Message: "File exceeds size limit.",
		},
		"INVALID_FILE_TYPE": {
			Code:    "FILE_004",
			Message: "Invalid file type.",
		},
		"UPLOAD_FAILED": {
			Code:    "FILE_005",
			Message: "Failed to upload file.",
		},
	},
	Branch: ErrorGroup{
		"BRANCH_NOT_FOUND": {
			Code:    "BRN_001",
			Message: "Branch not found.",
		},
		"BRANCH_DISABLED": {
			Code:    "BRN_002",
			Message: "Branch is disabled.",
		},
		"FAILED_TO_FETCH_BRANCHES": {
			Code:    "BRN_003",
			Message: "Failed to fetch branches.",
		},
		"INVALID_BRANCH_ID": {
			Code:    "BRN_004",
			Message: "Invalid branch ID provided.",
		},
		"INVALID_LOCATION_FILTER": {
			Code:    "BRN_005",
			Message: "Invalid location filter parameters.",
		},
		"BRANCH_ID_REQUIRED": {
			Code:    "BRN_006",
			Message: "Branch ID is required.",
		},
	},
}

func (e ErrorDefinition) Error() string {
	return e.Message
}
