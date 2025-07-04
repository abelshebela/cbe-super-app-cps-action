package common

// type ErrorDefinition struct {
// 	Code       string
// 	Message    string
// 	Details    string
// 	StatusCode int
// }

// type ErrorGroup map[string]ErrorDefinition

// type ErrorDefinitions struct {
// 	General     ErrorGroup
// 	Auth        ErrorGroup
// 	User        ErrorGroup
// 	Transaction ErrorGroup
// 	Account     ErrorGroup
// 	OTP         ErrorGroup
// 	File        ErrorGroup
// }

// var DefineError = ErrorDefinitions{
// 	General: ErrorGroup{
// 		"CONFLICT_KEY": {
// 			Code:       "GEN_001",
// 			Message:    "Duplicate key error: The specified field already exists.",
// 			StatusCode: 409,
// 		},
// 		"INVALID_ID": {
// 			Code:       "GEN_002",
// 			Message:    "Invalid ObjectId provided.",
// 			StatusCode: 400,
// 		},
// 		"INVALID_JSON_PAYLOAD": {
// 			Code:       "GEN_003",
// 			Message:    "Invalid JSON payload.",
// 			StatusCode: 400,
// 		},
// 		"UNHANDLED_SERVER_ERROR": {
// 			Code:       "GEN_004",
// 			Message:    "An unexpected server error occurred.",
// 			StatusCode: 500,
// 		},
// 		"INVALID_TOKEN_FORMAT": {
// 			Code:       "GEN_005",
// 			Message:    "Invalid token format.",
// 			StatusCode: 401,
// 		},
// 		"INVALID_TOKEN": {
// 			Code:       "GEN_006",
// 			Message:    "Invalid token.",
// 			StatusCode: 401,
// 		},
// 		"EXPIRED_TOKEN": {
// 			Code:       "GEN_007",
// 			Message:    "Token has expired.",
// 			StatusCode: 401,
// 		},
// 		"ERROR_CHANGING_PIN": {
// 			Code:       "GEN_030",
// 			Message:    "Error changing PIN.",
// 			StatusCode: 500,
// 		},
// 		"ENCRYPTED_PAYLOAD_REQUIRED": {
// 			Code:       "GEN_008",
// 			Message:    "Encrypted payload is required.",
// 			StatusCode: 400,
// 		},
// 		"INVALID_SERVER_CONFIGURATION": {
// 			Code:       "GEN_009",
// 			Message:    "Invalid server configuration",
// 			StatusCode: 500,
// 		},
// 		"AUTH_HEADER_MISSING": {
// 			Code:       "GEN_010",
// 			Message:    "Authorization header is missing. please provide a valid bearer token.",
// 			StatusCode: 401,
// 		},
// 		"INVALID_KEY_CONFIGURATION": {
// 			Code:       "GEN_011",
// 			Message:    "Invalid key configuration.",
// 			StatusCode: 500,
// 		},
// 		"INVALID_ENCRYPTED_PAYLOAD": {
// 			Code:       "GEN_012",
// 			Message:    "The provided encrypted payload is invalid.",
// 			StatusCode: 400,
// 		},
// 		"DECRYPTION_ERROR": {
// 			Code:       "GEN_013",
// 			Message:    "Decryption error occurred. please check your key configuration.",
// 			StatusCode: 500,
// 		},
// 		"SERVER_KEYS_NOT_CONFIGURED": {
// 			Code:       "GEN_015",
// 			Message:    "server error occurred. please try again later or contact support.",
// 			StatusCode: 500,
// 		},
// 		"INVALID_INPUT": {
// 			Code:       "GEN_016",
// 			Message:    "the provided input is invalid.",
// 			StatusCode: 400,
// 		},
// 		"UNAUTHORIZED": {
// 			Code:       "GEN_017",
// 			Message:    "missing authenticated user.",
// 			StatusCode: 401,
// 		},
// 		"USER_REALM_NOT_FOUND": {
// 			Code:       "GEN_018",
// 			Message:    "user is not permission to do this action.",
// 			StatusCode: 403,
// 		},
// 		"ACTION_NOT_ALLOWED": {
// 			Code:       "GEN_019",
// 			Message:    "user is not permission to do this action.",
// 			StatusCode: 403,
// 		},
// 		"WAIT_FOR_PREVIOUS_OTP_EXPIRATION": {
// 			Code:       "GEN_025",
// 			Message:    "wait until the pervious otp expired",
// 			StatusCode: 429,
// 		},
// 		"EMPTY_ORG": {
// 			Code:       "GEN_021",
// 			Message:    "Cant respond for empty organization id",
// 			StatusCode: 400,
// 		},
// 		"INVALID_REQ": {
// 			Code:       "GEN_022",
// 			Message:    "Invalid Request body",
// 			StatusCode: 400,
// 		},
// 		"REG_FRST": {
// 			Code:       "GEN_023",
// 			Message:    "please register on the Super APP and visit your nearest branch for Verification Code",
// 			StatusCode: 403,
// 		},
// 		"WEAK_PIN": {
// 			Code:       "GEN_024",
// 			Message:    "weak PIN used: please use a strong combination",
// 			StatusCode: 400,
// 		},
// 		"NOT_FOUND": {
// 			Code:       "GEN_020",
// 			Message:    "Resource not found.",
// 			StatusCode: 404,
// 		},
// 		"DEVICE_ID_REQUIRED": {
// 			Code:       "GEN_026",
// 			Message:    "Device UUID is required",
// 			StatusCode: 400,
// 		},
// 		"NO_LINKED_DEVICES": {
// 			Code:       "GEN_027",
// 			Message:    "User has no linked devices",
// 			StatusCode: 404,
// 		},
// 		"COULD_NOT_UNLINK_DEVICE": {
// 			Code:       "GEN_028",
// 			Message:    "Could not unlink device",
// 			StatusCode: 500,
// 		},
// 		"MISSING_REQUIRED_HEADERS": {
// 			Code:       "GEN_029",
// 			Message:    "Missing required headers for device lookup",
// 			StatusCode: 400,
// 		},
// 		"DEVICE_LOOKUP_FAILED": {
// 			Code:       "GEN_030",
// 			Message:    "Device lookup operation failed",
// 			StatusCode: 500,
// 		},
// 		"MISSING_REQUIRED_FIELDS": {
// 			Code:       "GEN_031",
// 			Message:    "Missing required fields in request",
// 			StatusCode: 400,
// 		},
// 		"MISSING_OTP": {
// 			Code:       "GEN_032",
// 			Message:    "OTP code is required",
// 			StatusCode: 400,
// 		},
// 		"INVALID_INPUT_PARAMETERS": {
// 			Code:       "GEN_033",
// 			Message:    "Invalid input parameters provided",
// 			StatusCode: 400,
// 		},
// 		"SAME_PIN": {
// 			Code:       "GEN_035",
// 			Message:    "New PIN cannot be the same as current PIN",
// 			StatusCode: 400,
// 		},
// 		"PIN_IN_HISTORY": {
// 			Code:       "GEN_036",
// 			Message:    "PIN has been used recently and cannot be reused",
// 			StatusCode: 400,
// 		},
// 		"ERROR_SETTING_PIN": {
// 			Code:       "GEN_037",
// 			Message:    "Failed to set PIN due to system error",
// 			StatusCode: 500,
// 		},
// 		"OTP_CREATION_FAILED": {
// 			Code:       "GEN_038",
// 			Message:    "Failed to create OTP",
// 			StatusCode: 500,
// 		},
// 		"DEVICE_NOT_FOUND": {
// 			Code:       "GEN_039",
// 			Message:    "Device not found in system",
// 			StatusCode: 404,
// 		},
// 		"DEVICE_FOUND": {
// 			Code:       "GEN_040",
// 			Message:    "Device found and registered",
// 			StatusCode: 200,
// 		},
// 		"INVALID_PIN": {
// 			Code:       "GEN_041",
// 			Message:    "Invalid PIN provided.",
// 			StatusCode: 400,
// 		},
// 		"PIN_RESET_SESSION_EXPIRED": {
// 			Code:       "GEN_042",
// 			Message:    "PIN reset session has expired",
// 			StatusCode: 410,
// 		},
// 		"PIN_RESET_SESSION_NOT_FOUND": {
// 			Code:       "GEN_043",
// 			Message:    "PIN reset session not found",
// 			StatusCode: 404,
// 		},
// 		"PIN_RESET_SESSION_ALREADY_VERIFIED": {
// 			Code:       "GEN_044",
// 			Message:    "PIN reset session already verified",
// 			StatusCode: 409,
// 		},
// 		"TOKEN_GENERATION_FAILED": {
// 			Code:       "GEN_045",
// 			Message:    "Failed to generate token",
// 			StatusCode: 500,
// 		},
// 		"MAXIMUM_AMOUNT_REQUIRED_FOR_OPEN_METHOD": {
// 			Code:       "GEN_046",
// 			Message:    "Maximum amount is required for open method",
// 			StatusCode: 400,
// 		},
// 		"MINIMUM_AMOUNT_REQUIRED_FOR_OTP_AND_PIN_METHOD": {
// 			Code:       "GEN_047",
// 			Message:    "Minimum amount is required for OTP and PIN method",
// 			StatusCode: 400,
// 		},
// 		"EITHER_MINIMUM_OR_MAXIMUM_AMOUNT_REQUIRED_FOR_PIN_METHOD": {
// 			Code:       "GEN_048",
// 			Message:    "Either minimum or maximum amount are required for PIN method",
// 			StatusCode: 400,
// 		},
// 		"INVALID_OBJECT_ID_FORMAT": {
// 			Code:       "GEN_049",
// 			Message:    "Invalid object ID format",
// 			StatusCode: 400,
// 		},
// 		"AUTH_TIER_NOT_FOUND": {
// 			Code:       "GEN_050",
// 			Message:    "Auth tier not found",
// 			StatusCode: 404,
// 		},
// 		"AUTH_TIER_ALREADY_EXISTS": {
// 			Code:       "GEN_051",
// 			Message:    "Auth tier already exists",
// 			StatusCode: 409,
// 		},
// 		"AUTH_TIER_VALIDATION_FAILED": {
// 			Code:       "GEN_052",
// 			Message:    "Auth tier validation failed",
// 			StatusCode: 400,
// 		},
// 		"AUTH_TIER_UPDATE_FAILED": {
// 			Code:       "GEN_053",
// 			Message:    "Auth tier update failed",
// 			StatusCode: 500,
// 		},
// 		"AUTH_TIER_INSERT_FAILED": {
// 			Code:       "GEN_054",
// 			Message:    "Auth tier insert failed",
// 			StatusCode: 500,
// 		},
// 		"AUTH_TIER_NOT_FOUND_FOR_ID": {
// 			Code:       "GEN_055",
// 			Message:    "Auth tier not found for the provided ID",
// 			StatusCode: 404,
// 		},
// 		"AUTH_TIER_FETCH_FAILED": {
// 			Code:       "GEN_056",
// 			Message:    "Auth tier fetch failed",
// 			StatusCode: 500,
// 		},
// 		"AUTH_TIER_APPROVE_FAILED": {
// 			Code:       "GEN_057",
// 			Message:    "Auth tier approve failed",
// 			StatusCode: 500,
// 		},
// 		"AUTH_TIER_REJECT_FAILED": {
// 			Code:       "GEN_058",
// 			Message:    "Auth tier reject failed",
// 			StatusCode: 500,
// 		},
// 		"AUTH_TIER_NOT_FOUND_FOR_ID_AND_DEPARTMENT": {
// 			Code:       "GEN_059",
// 			Message:    "Auth tier not found for the provided ID and department",
// 			StatusCode: 404,
// 		},
// 		"AUTH_TIER_FETCH_FAILED_FOR_ID_AND_DEPARTMENT": {
// 			Code:       "GEN_060",
// 			Message:    "Auth tier fetch failed for the provided ID and department",
// 			StatusCode: 500,
// 		},
// 		"AUTH_TIER_APPROVE_FAILED_FOR_ID_AND_DEPARTMENT": {
// 			Code:       "GEN_061",
// 			Message:    "Auth tier approve failed for the provided ID and department",
// 			StatusCode: 500,
// 		},
// 		"AUTH_TIER_REJECT_FAILED_FOR_ID_AND_DEPARTMENT": {
// 			Code:       "GEN_062",
// 			Message:    "Auth tier reject failed for the provided ID and department",
// 			StatusCode: 500,
// 		},
// 		"AUTH_TIER_UPDATE_FAILED_FOR_ID_AND_DEPARTMENT": {
// 			Code:       "GEN_063",
// 			Message:    "Auth tier update failed for the provided ID and department",
// 			StatusCode: 500,
// 		},
// 		"AUTH_TIER_INSERT_FAILED_FOR_ID_AND_DEPARTMENT": {
// 			Code:       "GEN_064",
// 			Message:    "Auth tier insert failed for the provided ID and department",
// 			StatusCode: 500,
// 		},
// 		"FAILED_TO_UPDATE_OPEN_AUTH_TIER": {
// 			Code:       "GEN_065",
// 			Message:    "Failed to update open auth tier",
// 			StatusCode: 500,
// 		},
// 		"FAILED_TO_UPDATE_PIN_AUTH_TIER": {
// 			Code:       "GEN_066",
// 			Message:    "Failed to update pin auth tier",
// 			StatusCode: 500,
// 		},
// 		"FAILED_TO_UPDATE_OTP_AND_PIN_AUTH_TIER": {
// 			Code:       "GEN_067",
// 			Message:    "Failed to update OTP and pin auth tier",
// 			StatusCode: 500,
// 		},
// 		"FAILED_TO_UPDATE_OTP_AND_OPEN_AUTH_TIER": {
// 			Code:       "GEN_068",
// 			Message:    "Failed to update OTP and open auth tier",
// 			StatusCode: 500,
// 		},
// 		"FAILED_TO_UPDATE_OPEN_AUTH_TIER_FOR_ID": {
// 			Code:       "GEN_069",
// 			Message:    "Failed to update open auth tier for the provided ID",
// 			StatusCode: 500,
// 		},
// 	},
// 	Auth: ErrorGroup{
// 		"AUTH_USER_NOT_FOUND": {
// 			Code:       "AUTH_001",
// 			Message:    "User not found.",
// 			StatusCode: 404,
// 		},
// 		"AUTH_USER_DISABLED": {
// 			Code:       "AUTH_002",
// 			Message:    "User is not allowed to login, please contact your admin!",
// 			StatusCode: 403,
// 		},
// 		"AUTH_INVALID_PASSWORD": {
// 			Code:       "AUTH_003",
// 			Message:    "Invalid password.",
// 			StatusCode: 401,
// 		},
// 		"AUTH_USER_HAS_NO_PASSWORD": {
// 			Code:       "AUTH_004",
// 			Message:    "Please reset your password. to login",
// 			StatusCode: 403,
// 		},
// 		"AUTH_TOO_MANY_ATTEMPTS": {
// 			Code:       "AUTH_005",
// 			Message:    `Too many attempts. Try again after ${waiting_time} minutes`,
// 			StatusCode: 429,
// 		},
// 		"AUTH_INVALID_OTP": {
// 			Code:       "AUTH_006",
// 			Message:    "The verification Code you entered is incorrect. Please check the Code and try again.",
// 			StatusCode: 400,
// 		},
// 		"AUTH_EXPIRED_OTP": {
// 			Code:       "AUTH_007",
// 			Message:    "The verification Code has expired. Please request a new Code to continue.",
// 			StatusCode: 410,
// 		},
// 		"AUTH_USER_RESET_PASSWORD_REQUIRED": {
// 			Code:       "AUTH_008",
// 			Message:    "Too many login attempts. Please reset your password. Your account is locked until then.",
// 			StatusCode: 423,
// 		},
// 		"AUTH_USER_ALREADY_EXISTS": {
// 			Code:       "AUTH_009",
// 			Message:    `user with this ${name} already exist`,
// 			StatusCode: 409,
// 		},
// 		"LOGIN_PROHIBITED_FOR_15_MIN": {
// 			Code:       "AUTH_010",
// 			Message:    "Incorrect credentials. Login is prohibited for 15 minutes.",
// 			StatusCode: 429,
// 		},
// 		"INCORRECT_PASSWORD": {
// 			Code:       "AUTH_011",
// 			Message:    `Incorrect password. ${triesLeft} login attempt(s) left.`,
// 			StatusCode: 401,
// 		},
// 		"OLD_PASSWORD_SAME_AS_NEW": {
// 			Code:       "AUTH_012",
// 			Message:    "New password can't be the same as old password.",
// 			StatusCode: 400,
// 		},
// 		"PIN_OLY_DIG": {
// 			Code:       "AUTH_013",
// 			Message:    "PIN must contain only digits",
// 			StatusCode: 400,
// 		},
// 		"PIN_LIMIT": {
// 			Code:       "AUTH_014",
// 			Message:    "PIN must be exactly 6 digits",
// 			StatusCode: 400,
// 		},
// 		"PIN_REDANDANT": {
// 			Code:       "AUTH_015",
// 			Message:    "PIN cannot contain more than 4 redundant numbers",
// 			StatusCode: 400,
// 		},
// 		"PIN_SEQ": {
// 			Code:       "AUTH_016",
// 			Message:    "PIN cannot contain 4 or more sequential numbers",
// 			StatusCode: 400,
// 		},
// 		"FAILD_TO_GEN_TOKEN": {
// 			Code:       "AUTH_017",
// 			Message:    "Set your new password",
// 			StatusCode: 500,
// 		},
// 		"FAILD_TO_RESET_PASS": {
// 			Code:       "AUTH_018",
// 			Message:    "Failed to set your new password",
// 			StatusCode: 500,
// 		},
// 		"FAILD_VALIDATION": {
// 			Code:       "AUTH_019",
// 			Message:    "Validation failed",
// 			StatusCode: 400,
// 		},
// 		"OLD_PIN_MISMATCH": {
// 			Code:       "AUTH_029",
// 			Message:    "Old PIN does not match",
// 			StatusCode: 400,
// 		},
// 		"SAME_PIN": {
// 			Code:       "AUTH_027",
// 			Message:    "New PIN cannot be the same as old PIN",
// 			StatusCode: 400,
// 		},
// 		"PIN_IN_HISTORY": {
// 			Code:       "AUTH_028",
// 			Message:    "New PIN cannot be one of the last 6 PINs used",
// 			StatusCode: 400,
// 		},
// 		"INVALID_BEARER": {
// 			Code:       "AUTH_020",
// 			Message:    "Unauthorized: Invalid Bearer token format",
// 			StatusCode: 401,
// 		},
// 		"USE_RIGHT_AUTH": {
// 			Code:       "AUTH_021",
// 			Message:    "Use the Right Authentication",
// 			StatusCode: 401,
// 		},
// 		"INVALID_CLAIM": {
// 			Code:       "AUTH_022",
// 			Message:    "Invalid token claims",
// 			StatusCode: 401,
// 		},
// 		"INVALID_TOKEN_DATA": {
// 			Code:       "AUTH_023",
// 			Message:    "Invalid token data",
// 			StatusCode: 401,
// 		},
// 		"UNABLE_TO_DYCRYPT_TOKEN": {
// 			Code:       "AUTH_024",
// 			Message:    "Unable to decrypt token",
// 			StatusCode: 500,
// 		},
// 		"FAILED_LOGIN": {
// 			Code:       "AUTH_025",
// 			Message:    "Failed to login. you entered wrong password",
// 			StatusCode: 401,
// 		},
// 		"ACCOUNT_LOCKED": {
// 			Code:       "AUTH_026",
// 			Message:    "account locked. please contact your addministrator",
// 			StatusCode: 423,
// 		},
// 		"MISSING_REQUIRED_FIELDS": {
// 			Code:       "AUTH_027",
// 			Message:    "Missing required fields in request",
// 			StatusCode: 400,
// 		},
// 		"PIN_RESET_SESSION_EXPIRED": {
// 			Code:       "AUTH_028",
// 			Message:    "PIN reset session has expired",
// 			StatusCode: 410,
// 		},
// 		"PIN_RESET_SESSION_NOT_FOUND": {
// 			Code:       "AUTH_029",
// 			Message:    "PIN reset session not found",
// 			StatusCode: 404,
// 		},
// 		"PIN_RESET_SESSION_ALREADY_VERIFIED": {
// 			Code:       "AUTH_030",
// 			Message:    "PIN reset session already verified",
// 			StatusCode: 409,
// 		},
// 	},
// 	Transaction: ErrorGroup{
// 		"TRANSACTION_NOT_FOUND": {
// 			Code:       "TXN_001",
// 			Message:    "Transaction not found.",
// 			StatusCode: 404,
// 		},
// 		"INSUFFICIENT_FUNDS": {
// 			Code:       "TXN_002",
// 			Message:    "Insufficient balance.",
// 			StatusCode: 402,
// 		},
// 		"COMMISSION_NOT_FOUND": {
// 			Code:       "COM_001",
// 			Message:    "Commission not found",
// 			StatusCode: 404,
// 		},
// 		"CUSTOMER_NOT_FOUND": {
// 			Code:       "CUS_001",
// 			Message:    "Customer not found",
// 			StatusCode: 404,
// 		},
// 	},
// 	Account: ErrorGroup{
// 		"PHONE_LOOKUP_FAILED": {
// 			Code:       "ACC_001",
// 			Message:    "Failed to lookup phone number.",
// 			StatusCode: 500,
// 		},
// 		"API_REQUEST_FAILED": {
// 			Code:       "ACC_002",
// 			Message:    "API request failed.",
// 			StatusCode: 500,
// 		},
// 		"ACCOUNT_NOT_FOUND": {
// 			Code:       "ACC_003",
// 			Message:    "Account not found.",
// 			StatusCode: 404,
// 		},
// 	},
// 	OTP: ErrorGroup{
// 		"INVALID_OTP": {
// 			Code:       "OTP_001",
// 			Message:    "Invalid OTP provided.",
// 			StatusCode: 400,
// 		},
// 		"EXPIRED_OTP": {
// 			Code:       "OTP_002",
// 			Message:    "OTP has expired.",
// 			StatusCode: 410,
// 		},
// 		"EMAIL_IN_USE": {
// 			Code:       "OTP_003",
// 			Message:    "Email is already in use.",
// 			StatusCode: 409,
// 		},
// 		"USER_ALREADY_HAS_EMAIL": {
// 			Code:       "USER_002",
// 			Message:    "User already has an email address",
// 			StatusCode: 409,
// 		},
// 		"WAIT_FOR_PREVIOUS_OTP_EXPIRATION": {
// 			Code:       "USER_003",
// 			Message:    "Please wait until the previous OTP expires",
// 			StatusCode: 429,
// 		},
// 		"OTP_GENERATION_FAILED": {
// 			Code:       "OTP_004",
// 			Message:    "Failed to generate OTP.",
// 			StatusCode: 500,
// 		},
// 	},
// 	File: ErrorGroup{
// 		"INVALID_FORM": {
// 			Code:       "FILE_001",
// 			Message:    "Could not parse form or file too large.",
// 			StatusCode: 400,
// 		},
// 		"NO_FILE": {
// 			Code:       "FILE_002",
// 			Message:    "File is required.",
// 			StatusCode: 400,
// 		},
// 		"FILE_TOO_LARGE": {
// 			Code:       "FILE_003",
// 			Message:    "File exceeds size limit.",
// 			StatusCode: 413,
// 		},
// 		"INVALID_FILE_TYPE": {
// 			Code:       "FILE_004",
// 			Message:    "Invalid file type.",
// 			StatusCode: 415,
// 		},
// 		"UPLOAD_FAILED": {
// 			Code:       "FILE_005",
// 			Message:    "Failed to upload file.",
// 			StatusCode: 500,
// 		},
// 	},
// }

// // Registration errors
// var RegistrationErrors = map[string]ErrorDefinition{
// 	"PHONE_ALREADY_EXISTS": {
// 		Code:       "REG001",
// 		Message:    "A user with this phone number already exists",
// 		Details:    "The provided phone number is already registered in the system",
// 		StatusCode: 409,
// 	},
// 	"DEVICE_ALREADY_REGISTERED": {
// 		Code:       "REG002",
// 		Message:    "This device is already registered",
// 		Details:    "The provided device UUID is already associated with an account",
// 		StatusCode: 409,
// 	},
// 	"INVALID_PHONE_NUMBER": {
// 		Code:       "REG003",
// 		Message:    "Invalid phone number format",
// 		Details:    "Please provide a valid phone number in the correct format",
// 		StatusCode: 400,
// 	},
// 	"INVALID_PLATFORM": {
// 		Code:       "REG004",
// 		Message:    "Invalid platform specified",
// 		Details:    "Platform must be one of: android, ios, web",
// 		StatusCode: 400,
// 	},
// 	"INVALID_DEVICE_UUID": {
// 		Code:       "REG005",
// 		Message:    "Invalid device UUID format",
// 		Details:    "The provided device UUID appears to be invalid",
// 		StatusCode: 400,
// 	},
// 	"REGISTRATION_IN_PROGRESS": {
// 		Code:       "REG006",
// 		Message:    "Registration already in progress",
// 		Details:    "A registration process is already active for this phone number",
// 		StatusCode: 409,
// 	},
// 	"REGISTRATION_RATE_LIMITED": {
// 		Code:       "REG007",
// 		Message:    "Too many registration attempts",
// 		Details:    "Please wait before attempting to register again",
// 		StatusCode: 429,
// 	},
// 	"REGISTRATION_FAILED": {
// 		Code:       "REG008",
// 		Message:    "Registration process failed",
// 		Details:    "An error occurred during the registration process",
// 		StatusCode: 500,
// 	},
// 	"DEVICE_UUID_MISMATCH": {
// 		Code:       "REG009",
// 		Message:    "Device UUID mismatch",
// 		Details:    "Device UUID in header does not match the one in request body",
// 		StatusCode: 400,
// 	},
// }

// func GetErrorByCode(code string) (ErrorDefinition, bool) {
// 	// Search in all error groups
// 	groups := []ErrorGroup{
// 		DefineError.General,
// 		DefineError.Auth,
// 		DefineError.User,
// 		DefineError.Transaction,
// 		DefineError.Account,
// 		DefineError.OTP,
// 		DefineError.File,
// 	}
// 	for _, group := range groups {
// 		for _, errDef := range group {
// 			if errDef.Code == code {
// 				return errDef, true
// 			}
// 		}
// 	}
// 	// Search in registration errors
// 	for _, errDef := range RegistrationErrors {
// 		if errDef.Code == code {
// 			return errDef, true
// 		}
// 	}
// 	return ErrorDefinition{}, false
// }
