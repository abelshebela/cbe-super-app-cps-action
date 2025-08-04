package common

type ErrorDefinition struct {
	Code    string
	Message string
	Details string
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
		"MISSING_REQUIRED_HEADERS": {
			Code:    "GEN_029",
			Message: "Missing required headers for device lookup",
		},
		"DEVICE_LOOKUP_FAILED": {
			Code:    "GEN_030",
			Message: "Device lookup operation failed",
		},
		"MISSING_REQUIRED_FIELDS": {
			Code:    "GEN_031",
			Message: "Missing required fields in request",
		},
		"MISSING_OTP": {
			Code:    "GEN_032",
			Message: "OTP code is required",
		},
		"INVALID_INPUT_PARAMETERS": {
			Code:    "GEN_033",
			Message: "Invalid input parameters provided",
		},
		"SAME_PIN": {
			Code:    "GEN_035",
			Message: "New PIN cannot be the same as current PIN",
		},
		"PIN_IN_HISTORY": {
			Code:    "GEN_036",
			Message: "PIN has been used recently and cannot be reused",
		},
		"ERROR_SETTING_PIN": {
			Code:    "GEN_037",
			Message: "Failed to set PIN due to system error",
		},
		"OTP_CREATION_FAILED": {
			Code:    "GEN_038",
			Message: "Failed to create OTP",
		},
		"DEVICE_NOT_FOUND": {
			Code:    "GEN_039",
			Message: "Device not found in system",
		},
		"DEVICE_FOUND": {
			Code:    "GEN_040",
			Message: "Device found and registered",
		},
		"INVALID_PIN": {
			Code:    "GEN_041",
			Message: "Invalid PIN provided.",
		},
		"PIN_RESET_SESSION_EXPIRED": {
			Code:    "GEN_042",
			Message: "PIN reset session has expired",
		},
		"PIN_RESET_SESSION_NOT_FOUND": {
			Code:    "GEN_043",
			Message: "PIN reset session not found",
		},
		"PIN_RESET_SESSION_ALREADY_VERIFIED": {
			Code:    "GEN_044",
			Message: "PIN reset session already verified",
		},
		"TOKEN_GENERATION_FAILED": {
			Code:    "GEN_045",
			Message: "Failed to generate token",
		},
		"OLDER_DEVICE": {
			Code:    "GEN_046",
			Message: "You have an older device please update",
		},
		"USER_DISABLED_BLOCKED": {
			Code: "GEN_047",
			Message: "User disabled or blocked",
		}
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
		"MISSING_REQUIRED_FIELDS": {
			Code:    "AUTH_027",
			Message: "Missing required fields in request",
		},
		"PIN_RESET_SESSION_EXPIRED": {
			Code:    "AUTH_028",
			Message: "PIN reset session has expired",
		},
		"PIN_RESET_SESSION_NOT_FOUND": {
			Code:    "AUTH_029",
			Message: "PIN reset session not found",
		},
		"PIN_RESET_SESSION_ALREADY_VERIFIED": {
			Code:    "AUTH_030",
			Message: "PIN reset session already verified",
		},
		"DEVICE_UUID_MISMATCH": {
			Code:    "AUTH_031",
			Message: "Device UUID mismatch",
		},
		"FAILED_TO_DELETE_PENDING_REG": {
			Code:    "AUTH_032",
			Message: "Failed to delete expired registration",
		},
		"PLEASE_REG_FIRST": {
			Code:    "AUTH_033",
			Message: "User Not Found Please register first",
		},
		"PROFILE_SET_ERROR": {
			Code:    "AUTH_033",
			Message: "Failed to set Profile",
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
}

// Registration errors
var RegistrationErrors = map[string]ErrorDefinition{
	"PHONE_ALREADY_EXISTS": {
		Code:    "REG001",
		Message: "A user with this phone number already exists",
		Details: "The provided phone number is already registered in the system",
	},
	"DEVICE_ALREADY_REGISTERED": {
		Code:    "REG002",
		Message: "This device is already registered",
		Details: "The provided device UUID is already associated with an account",
	},
	"INVALID_PHONE_NUMBER": {
		Code:    "REG003",
		Message: "Invalid phone number format",
		Details: "Please provide a valid phone number in the correct format",
	},
	"INVALID_PLATFORM": {
		Code:    "REG004",
		Message: "Invalid platform specified",
		Details: "Platform must be one of: android, ios, web",
	},
	"INVALID_DEVICE_UUID": {
		Code:    "REG005",
		Message: "Invalid device UUID format",
		Details: "The provided device UUID appears to be invalid",
	},
	"REGISTRATION_IN_PROGRESS": {
		Code:    "REG006",
		Message: "Registration already in progress",
		Details: "A registration process is already active for this phone number",
	},
	"REGISTRATION_RATE_LIMITED": {
		Code:    "REG007",
		Message: "Too many registration attempts",
		Details: "Please wait before attempting to register again",
	},
	"REGISTRATION_FAILED": {
		Code:    "REG008",
		Message: "Registration process failed",
		Details: "An error occurred during the registration process",
	},
	"DEVICE_UUID_MISMATCH": {
		Code:    "REG009",
		Message: "Device UUID mismatch",
		Details: "Device UUID in header does not match the one in request body",
	},
}
