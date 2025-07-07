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
	Department  ErrorGroup
	Bank        ErrorGroup
	Action      ErrorGroup
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
			Code:    "GEN_014",
			Message: "server error occurred. please try again later or contact support.",
		},
		"INVALID_INPUT": {
			Code:    "GEN_015",
			Message: "the provided input is invalid.",
		},
		"UNAUTHORIZED": {
			Code:    "GEN_016",
			Message: "missing authenticated user.",
		},
		"USER_REALM_NOT_FOUND": {
			Code:    "GEN_017",
			Message: "user is not permission to do this action.",
		},
		"ACTION_NOT_ALLOWED": {
			Code:    "GEN_018",
			Message: "user does not have permission to do this action.",
		},
		"NOT_FOUND": {
			Code:    "GEN_019",
			Message: "Resource not found.",
		},
		"EMPTY_ORG": {
			Code:    "GEN_020",
			Message: "Cant respond for empty organization id",
		},
		"INVALID_REQ": {
			Code:    "GEN_021",
			Message: "Invalid Request body",
		},
		"REG_FRST": {
			Code:    "GEN_022",
			Message: "please register on the Super APP and visit your nearest branch for Verification Code",
		},
		"WEAK_PIN": {
			Code:    "GEN_023",
			Message: "weak PIN used: please use a strong combination",
		},
		"WAIT_FOR_PREVIOUS_OTP_EXPIRATION": {
			Code:    "GEN_024",
			Message: "wait until the pervious otp expired",
		},
		"INCOMPLETE_USER_INFO": {
			Code:    "GEN_025",
			Message: "Incomplete user information",
		},
		"INVALID_ACTION_TYPE": {
			Code:    "GEN_026",
			Message: "Invalid action type",
		},
		"COULD_NOT_UNLINK_DEVICE": {
			Code:    "GEN_027",
			Message: "Could not unlink device",
		},
		"ACTION_NOT_FOUND": {
			Code:    "GEN_028",
			Message: "Action not found.",
		},
		"ERROR_CHANGING_PIN": {
			Code:    "GEN_029",
			Message: "Error changing PIN.",
		},
		"MISSING_REQUIRED_FIELDS": {
			Code:    "GEN_030",
			Message: "Missing required fields in request",
		},
		"MISSING_OTP": {
			Code:    "GEN_031",
			Message: "OTP code is required",
		},
		"INVALID_INPUT_PARAMETERS": {
			Code:    "GEN_032",
			Message: "Invalid input parameters provided",
		},
		"PENDING_REQUEST_EXISTS": {
			Code:    "GEN_033",
			Message: "You have a pending request for this action.",
		},
		"SAME_PIN": {
			Code:    "GEN_034",
			Message: "New PIN cannot be the same as current PIN",
		},
		"GENERAL_DB_QUERY_FAILED": {
			Code:    "GEN_035",
			Message: "Database query failed.",
		},
		"ERROR_SETTING_PIN": {
			Code:    "GEN_036",
			Message: "Failed to set PIN due to system error",
		},
		"OTP_CREATION_FAILED": {
			Code:    "GEN_037",
			Message: "Failed to create OTP",
		},
		"DEVICE_NOT_FOUND": {
			Code:    "GEN_038",
			Message: "Device not found in system",
		},
		"DEVICE_FOUND": {
			Code:    "GEN_039",
			Message: "Device found and registered",
		},
		"INVALID_PIN": {
			Code:    "GEN_040",
			Message: "Invalid PIN provided.",
		},
		"PIN_RESET_SESSION_EXPIRED": {
			Code:    "GEN_041",
			Message: "PIN reset session has expired",
		},
		"PIN_RESET_SESSION_NOT_FOUND": {
			Code:    "GEN_042",
			Message: "PIN reset session not found",
		},
		"PIN_RESET_SESSION_ALREADY_VERIFIED": {
			Code:    "GEN_043",
			Message: "PIN reset session already verified",
		},
		"FAILED_TO_GET_FEEDBACK_DATA": {
			Code:    "GEN_044",
			Message: "failed to get feedback data",
		},
		"MAXIMUM_AMOUNT_REQUIRED_FOR_OPEN_METHOD": {
			Code:    "GEN_045",
			Message: "Maximum amount is required for open method",
		},
		"FAILED_TO_GET_FEEDBACK_COUNTS": {
			Code:    "GEN_046",
			Message: "failed to get feedback total counts",
		},
		"EITHER_MINIMUM_OR_MAXIMUM_AMOUNT_REQUIRED_FOR_PIN_METHOD": {
			Code:    "GEN_047",
			Message: "Either minimum or maximum amount are required for PIN method",
		},
		"INVALID_OBJECT_ID_FORMAT": {
			Code:    "GEN_048",
			Message: "Invalid object ID format",
		},
		"AUTH_TIER_NOT_FOUND": {
			Code:    "GEN_049",
			Message: "Auth tier not found",
		},
		"AUTH_TIER_ALREADY_EXISTS": {
			Code:    "GEN_050",
			Message: "Auth tier already exists",
		},
		"AUTH_TIER_VALIDATION_FAILED": {
			Code:    "GEN_051",
			Message: "Auth tier validation failed",
		},
		"SERVICE_TYPE_CANNOT_BE_EMPTY": {
			Code:    "GEN_052",
			Message: "Service type cannot be empty",
		},
		"AUTH_TIER_INSERT_FAILED": {
			Code:    "GEN_053",
			Message: "Auth tier insert failed",
		},
		"AUTH_TIER_NOT_FOUND_FOR_ID": {
			Code:    "GEN_054",
			Message: "Auth tier not found for the provided ID",
		},
		"AUTH_TIER_FETCH_FAILED": {
			Code:    "GEN_055",
			Message: "Auth tier fetch failed",
		},
		"AUTH_TIER_APPROVE_FAILED": {
			Code:    "GEN_056",
			Message: "Auth tier approve failed",
		},
		"AUTH_TIER_REJECT_FAILED": {
			Code:    "GEN_057",
			Message: "Auth tier reject failed",
		},
		"AUTH_TIER_NOT_FOUND_FOR_ID_AND_DEPARTMENT": {
			Code:    "GEN_058",
			Message: "Auth tier not found for the provided ID and department",
		},
		"AUTH_TIER_FETCH_FAILED_FOR_ID_AND_DEPARTMENT": {
			Code:    "GEN_059",
			Message: "Auth tier fetch failed for the provided ID and department",
		},
		"AUTH_TIER_APPROVE_FAILED_FOR_ID_AND_DEPARTMENT": {
			Code:    "GEN_060",
			Message: "Auth tier approve failed for the provided ID and department",
		},
		"AUTH_TIER_REJECT_FAILED_FOR_ID_AND_DEPARTMENT": {
			Code:    "GEN_061",
			Message: "Auth tier reject failed for the provided ID and department",
		},
		"FAILED_TO_UPDATE_USER": {
			Code:    "GEN_062",
			Message: "failed to update user",
		},
		"AUTH_TIER_INSERT_FAILED_FOR_ID_AND_DEPARTMENT": {
			Code:    "GEN_063",
			Message: "Auth tier insert failed for the provided ID and department",
		},
		"ACTION_ID_IS_REQUIRED": {
			Code:    "GEN_064",
			Message: "action id is required",
		},
		"FAILED_TO_UPDATE_PIN_AUTH_TIER": {
			Code:    "GEN_065",
			Message: "Failed to update pin auth tier",
		},
		"FAILED_TO_UPDATE_OTP_AND_PIN_AUTH_TIER": {
			Code:    "GEN_066",
			Message: "Failed to update OTP and pin auth tier",
		},
		"FAILED_TO_UPDATE_OTP_AND_OPEN_AUTH_TIER": {
			Code:    "GEN_067",
			Message: "Failed to update OTP and open auth tier",
		},
		"FAILED_TO_UPDATE_OPEN_AUTH_TIER_FOR_ID": {
			Code:    "GEN_068",
			Message: "Failed to update open auth tier for the provided ID",
		},
		"MIN_AMOUNT_CANNOT_BE_GREATER_THAN_PIN_MIN_AMOUNT": {
			Code:    "GEN_069",
			Message: "Min amount cannot be greater than or equal to pin min amount",
		},
		"PIN_AUTHIER_NOT_FOUND": {
			Code:    "GEN_070",
			Message: "Pin authier not found",
		},
		"PIN_MIN_AMOUNT_CANNOT_BE_LESS_THAN_OPEN_MIN_AMOUNT": {
			Code:    "GEN_071",
			Message: "Pin min amount cannot be less than or equal to open min amount",
		},
		"MAX_PIN_AMOUNT_SHOULD_BE_GREATER_THAN_PIN_MIN_AMOUNT": {
			Code:    "GEN_072",
			Message: "Max pin amount should be greater than pin min amount",
		},
		"OPEN_AUTHIER_NOT_FOUND": {
			Code:    "GEN_073",
			Message: "Open authier not found",
		},
		"INVALID_ACTION_DATA": {
			Code:    "GEN_074",
			Message: "Invalid action data.",
		},
		"PENDING_CPS_ACTION_PRESENT": {
			Code:    "GEN_075",
			Message: "Pending cps action present",
		},
		"FAILED_TO_GET_FAYDA_ACCOUNT": {
			Code:    "GEN_076",
			Message: "Failed to get fayda account",
		},
		"FAILED_TO_GET_CUSTOMER_ACCOUNT": {
			Code:    "GEN_077",
			Message: "Failed to get customer account",
		},
		"FAILED_TO_CREATE_CPS_ACTION": {
			Code:    "GEN_078",
			Message: "Failed to create cps action",
		},
		"AT_LEAST_ONE_TIER_IS_REQUIRED": {
			Code:    "GEN_079",
			Message: "At least one tier is required",
		},
		"THE_FIRST_TIER_MINIMUM_MUST_START_FROM_0": {
			Code:    "GEN_080",
			Message: "The first tier's minimum must start from 0",
		},
		"TIER_MINIMUM_MUST_EQUAL_TO_PREVIOUS_TIER_MAXIMUM": {
			Code:    "GEN_081",
			Message: "Tier minimum must equal the previous tier's maximum",
		},
		"TIER_MAXIMUM_MUST_BE_GREATER_THAN_PREVIOUS_TIER_MAXIMUM": {
			Code:    "GEN_082",
			Message: "Tier maximum must be greater than the previous tier's maximum",
		},
		"ABOVE_AMOUNT_MUST_MATCH_LAST_TIER_MAXIMUM": {
			Code:    "GEN_083",
			Message: "Above amount must match the last tier's maximum",
		},
		"SERVICE_CODE_CANNOT_BE_EMPTY": {
			Code:    "GEN_084",
			Message: "Service code cannot be empty",
		},
		"SERVICE_NAME_CANNOT_BE_EMPTY": {
			Code:    "GEN_085",
			Message: "Service name cannot be empty",
		},
		"NO_LINKED_DEVICES": {
			Code:    "GEN_086",
			Message: "User has no linked devices",
		},
		"FAILED_TO_GET_SERVICE": {
			Code:    "GEN_087",
			Message: "Failed to get service",
		},
		"SERVICE_NOT_FOUND": {
			Code:    "GEN_088",
			Message: "Service not found",
		},
		"DEVICE_LOOKUP_FAILED": {
			Code:    "GEN_089",
			Message: "Device lookup operation failed",
		},
		"FAILED_TO_UPDATE_SERVICE_FEE": {
			Code:    "GEN_090",
			Message: "Failed to update service fee",
		},
		"FAILED_TO_REJECT_SERVICE_FEE_UPDATE": {
			Code:    "GEN_091",
			Message: "Failed to reject service fee update",
		},
		"AUTH_TIER_UPDATE_FAILED": {
			Code:    "GEN_092",
			Message: "Auth tier update failed",
		},
		"ACTION_IS_NOT_IN_PENDING_STATUS": {
			Code:    "GEN_093",
			Message: "Action is not in pending status",
		},
		"FAILED_TO_UPDATE_SERVICE_DETAILS": {
			Code:    "GEN_094",
			Message: "Failed to update service details",
		},
		"DEVICE_ID_REQUIRED": {
			Code:    "GEN_095",
			Message: "Device UUID is required",
		},
		"TOKEN_GENERATION_FAILED": {
			Code:    "GEN_096",
			Message: "Failed to generate token",
		},
		"FAILED_TO_GET_FEEDBACK": {
			Code:    "GEN_097",
			Message: "failed to get feedback",
		},
		"PIN_IN_HISTORY": {
			Code:    "GEN_098",
			Message: "PIN has been used recently and cannot be reused",
		},
		"MINIMUM_AMOUNT_REQUIRED_FOR_OTP_AND_PIN_METHOD": {
			Code:    "GEN_099",
			Message: "Minimum amount is required for OTP and PIN method",
		},
		"FAILED_TO_CONVERT_ID": {
			Code:    "GEN_100",
			Message: "failed to convert ID",
		},
		"FAIL_TO_FIND_USER": {
			Code:    "GEN_101",
			Message: "failed to find user",
		},
		"DUPLICATE_KEY": {
			Code:    "GEN_102",
			Message: "duplicate key error",
		},
		"INVALID_SERVICE_DATA_IN_ACTION": {
			Code:    "GEN_103",
			Message: "Invalid service data in action",
		},
		"FAIL_TO_FEATCH_CPS_USER": {
			Code:    "GEN_104",
			Message: "fail to featch cps user",
		},
		"AUTH_TIER_UPDATE_FAILED_FOR_ID_AND_DEPARTMENT": {
			Code:    "GEN_105",
			Message: "Auth tier update failed for the provided ID and department",
		},
		"FAILED_TO_UPDATE_CPS_ACTION": {
			Code:    "GEN_106",
			Message: "failed to update cps action",
		},
		"OPEN_TIER_MAX_AMOUNT_CANNOT_BE_GREATER_THAN_PIN_MAX_AMOUNT": {
			Code:    "GEN_107",
			Message: "Open tier max amount can not be greater than pin max amount",
		},
		"FAILED_TO_UPDATE_OPEN_AUTH_TIER": {
			Code:    "GEN_108",
			Message: "Failed to update open auth tier",
		},
		"MISSING_REQUIRED_HEADERS": {
			Code:    "GEN_109",
			Message: "Missing required headers for device lookup",
		},
		"USER_CODE_IS_REQUIRED": {
			Code:    "GEN_110",
			Message: "USER CODE IS REQUIRED",
		},
		"FAILED_TO_FIND_CPS_ACTION": {
			Code:    "GEN_111",
			Message: "Failed to find cps action",
		},
		"CPS_USER_NOT_FOUND": {
			Code:    "GEN_112",
			Message: "cps user not found",
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
			Code:    "AUTH_020",
			Message: "Old PIN does not match",
		},
		"PIN_IN_HISTORY": {
			Code:    "AUTH_021",
			Message: "New PIN cannot be one of the last 6 PINs used",
		},
		"INVALID_BEARER": {
			Code:    "AUTH_022",
			Message: "Unauthorized: Invalid Bearer token format",
		},
		"USE_RIGHT_AUTH": {
			Code:    "AUTH_023",
			Message: "Use the Right Authentication",
		},
		"INVALID_CLAIM": {
			Code:    "AUTH_024",
			Message: "Invalid token claims",
		},
		"INVALID_TOKEN_DATA": {
			Code:    "AUTH_025",
			Message: "Invalid token data",
		},
		"UNABLE_TO_DYCRYPT_TOKEN": {
			Code:    "AUTH_026",
			Message: "Unable to decrypt token",
		},
		"FAILED_LOGIN": {
			Code:    "AUTH_027",
			Message: "Failed to login. you entered wrong password",
		},
		"ACCOUNT_LOCKED": {
			Code:    "AUTH_028",
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
			Code:    "TXN_003",
			Message: "Commission not found",
		},
		"CUSTOMER_NOT_FOUND": {
			Code:    "TXN_004",
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
		// Additional missing error from the other definition
		"ACCOUNT_UPDATE_FAILED": {
			Code:    "ACC_014",
			Message: "Failed to update account.",
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
			Code:    "OTP_004",
			Message: "User already has an email address",
		},
		"WAIT_FOR_PREVIOUS_OTP_EXPIRATION": {
			Code:    "OTP_005",
			Message: "Please wait until the previous OTP expires",
		},
		"OTP_GENERATION_FAILED": {
			Code:    "OTP_006",
			Message: "Failed to generate OTP.",
		},
		"USER_KYC_LEVEL_ZERO": {
			Code:    "OTP_007",
			Message: "User KYC level is not sufficient for this operation",
		},
		"USER_KYC_LEVEL_WRONG": {
			Code:    "OTP_008",
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
		// Additional missing error from the other definition
		"MISSING_OR_INVALID_LOGO": {
			Code:    "FILE_006",
			Message: "Missing or invalid logo.",
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
		// Additional missing error from the other definition
		"BRANCHS_ARE_REQUIRED": {
			Code:    "BRN_007",
			Message: "Branches are required.",
		},
	},
	Department: ErrorGroup{
		"DEPARTMENT_NOT_FOUND": {
			Code:    "DEP_001",
			Message: "Department does not exist.",
		},
		"DEPARTMENT_ALREADY_EXISTS": {
			Code:    "DEP_002",
			Message: "Department already exists.",
		},
		"DEPARTMENT_CODE_REQUIRED": {
			Code:    "DEP_003",
			Message: "Missing department_code for update.",
		},
		"DEPARTMENT_NAME_REQUIRED": {
			Code:    "DEP_004",
			Message: "Department name is required.",
		},
		"PORTAL_CARDS_INVALID": {
			Code:    "DEP_005",
			Message: "Portal cards must be a list of strings.",
		},
	},
	Bank: ErrorGroup{
		"BANKS_NOT_FOUND": {
			Code:    "BNK_001",
			Message: "Banks data not found.",
		},
		"INVALID_BANK_NAME": {
			Code:    "BNK_002",
			Message: "Bank name must be between 3 and 10 alphabetic characters.",
		},
		"MISSING_BANK_NAME": {
			Code:    "BNK_003",
			Message: "Bank name is required.",
		},
		"MISSING_BANK_CODE": {
			Code:    "BNK_004",
			Message: "Bank code is required.",
		},
		"MISSING_BANK_BIC": {
			Code:    "BNK_005",
			Message: "Bank identifier code (BIC) is required.",
		},
	},
	Action: ErrorGroup{
		"PENDING_CPS_ACTION_EXISTS": {
			Code:    "ACT_001",
			Message: "A pending CPS action already exists.",
		},
		"INVALID_DECISION": {
			Code:    "ACT_002",
			Message: "Invalid decision value provided.",
		},
		"ACTION_REJECTION_FAILED": {
			Code:    "ACT_003",
			Message: "Failed to reject action.",
		},
		"ACTION_APPROVAL_FAILED": {
			Code:    "ACT_004",
			Message: "Failed to approve action.",
		},
		"UNLINK_ACTION_REQUEST_FAILED": {
			Code:    "ACT_005",
			Message: "Failed to request unlink action.",
		},
		"PENDING_ACTION_REJECTION_FAILED": {
			Code:    "ACT_006",
			Message: "Failed to reject pending action.",
		},
		"PENDING_ACTION_CHECK_FAILED": {
			Code:    "ACT_007",
			Message: "Failed to check pending actions.",
		},
	},
	User: ErrorGroup{
		"USER_STATUS_UPDATE_FAILED": {
			Code:    "USR_001",
			Message: "Failed to update user status.",
		},
	},
}

func (e ErrorDefinition) Error() string {
	return e.Message
}
