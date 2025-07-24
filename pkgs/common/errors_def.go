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
	Region      ErrorGroup
	Department  ErrorGroup
	Bank        ErrorGroup
	Action      ErrorGroup
	Wallet      ErrorGroup
	AD          ErrorGroup
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
		"REQUIRED_FIELDS_MISSING": {
			Code:    "GEN_113",
			Message: "One or more required fields are missing.",
		},
		"VALIDATION_RULE_ID_MISMATCH": {
			Code:    "GEN_114",
			Message: "Validation rule ID does not match the expected value.",
		},
		"FAILED_TO_UPDATE_VALIDATION_RULE": {
			Code:    "GEN_115",
			Message: "Failed to update the validation rule.",
		},
		"FAILED_TO_UNMARSHAL_CURRENT_ACTION": {
			Code:    "GEN_116",
			Message: "Failed to parse the current action data.",
		},
		"CURRENT_ACTION_INVALID_TYPE": {
			Code:    "GEN_117",
			Message: "Current action has an invalid data type.",
		},
		"CURRENT_ACTION_NIL": {
			Code:    "GEN_118",
			Message: "Current action is missing or nil.",
		},
		"ACTION_NOT_PENDING": {
			Code:    "GEN_119",
			Message: "Action is not in a pending state.",
		},
		"CHECKER_ID_EMPTY": {
			Code:    "GEN_120",
			Message: "Checker ID is required but was not provided.",
		},
		"ACTION_ID_EMPTY": {
			Code:    "GEN_121",
			Message: "Action ID is required but was not provided.",
		},
		"FAILED_TO_MARSHAL_CURRENT_ACTION": {
			Code:    "GEN_122",
			Message: "Failed to serialize the current action data.",
		},
		"FAILED_TO_MARSHAL_PREVIOUS_ACTION": {
			Code:    "GEN_123",
			Message: "Failed to serialize the previous action data.",
		},
		"PENDING_ACTION_EXISTS": {
			Code:    "GEN_124",
			Message: "A pending action already exists for this validation rule.",
		},
		"FAILED_TO_FETCH_PENDING_ACTIONS": {
			Code:    "GEN_125",
			Message: "Failed to fetch pending actions for this validation rule.",
		},
		"FAILED_TO_GET_AD": {
			Code:    "GEN_213",
			Message: "Failed to get ad.",
		},
		"ADVERT_NOT_FOUND": {
			Code:    "GEN_214",
			Message: "advert not found",
		},
		"FAILED_TO_CAST_ACTION_DATA": {
			Code:    "GEN_215",
			Message: "failed to cast action data to ad request",
		},
		"FAILED_TO_CHECK_AD_BUCKET": {
			Code:    "GEN_216",
			Message: "failed to check ad bucket",
		},
		"FAILED_TO_CREATE_AD_BUCKET": {
			Code:    "GEN_217",
			Message: "failed to create ad bucket",
		},
		"FAILED_TO_OPEN_FILE": {
			Code:    "GEN_218",
			Message: "failed to create ad bucket",
		},
		"FAILED_TO_SAVE_OBJECT_TO_MINIO": {
			Code:    "GEN_219",
			Message: "failed to save object to MinIO",
		},
		"TIERS_REQUIRED": {
			Code:    "GEN_126",
			Message: "At least one tier is required.",
		},
		"TIERS_FIRST_MIN_ZERO": {
			Code:    "GEN_127",
			Message: "The first tier's minimum must start from 0.",
		},
		"TIERS_MIN_MUST_EQUAL_PREV_MAX": {
			Code:    "GEN_127",
			Message: "Tier minimum must equal the previous tier's maximum.",
		},
		"TIERS_MAX_MUST_INCREASE": {
			Code:    "GEN_129",
			Message: "Tier maximum must be greater than the previous tier's maximum.",
		},
		"TIERS_ABOVE_AMOUNT_MISMATCH": {
			Code:    "GEN_130",
			Message: "Above amount must match the last tier's maximum value.",
		},
		"NOT_IMPLEMENTED": {
			Code:    "GEN_131",
			Message: "Not implemented.",
		},
		"INPUT_TOO_LONG": {
			Code:    "GEN_132",
			Message: "Input exceeds maximum allowed length.",
		},
		"INPUT_INVALID_CHARACTERS": {
			Code:    "GEN_133",
			Message: "Input contains invalid characters.",
		},
		"PAGE_NOT_FOUND": {
			Code:    "GEN_134",
			Message: "Page not found",
		},
		"FAILED_TO_GET_AVATAR": {
			Code:    "GEN_135",
			Message: "failed to get avatar",
		},
		"PENDING_REQUEST_CHECK_FAILED_FOR_CREATE_PERMISSION": {
			Code:    "GEN_136",
			Message: "Permission group already exists with name",
		},
		"PERMISSION_GROUP_ALREADY_EXIXTS": {
			Code:    "GEN_137",
			Message: "Permission group already exists with name",
		},
		"FAILED_TO_UPDATE_USER_DATA": {
			Code:    "GEN_132",
			Message: "failed to update user data",
		},
		"FAILED_TO_INSERT_CPS_ACTION": {
			Code:    "GEN_133",
			Message: "failed to insert cps action",
		},
		"ACCOUNT_ALREADY_DISABLED": {
			Code:    "GEN_134",
			Message: "Account is already disabled",
		},
		"NO_PENDING_ACTION_FOUND": {
			Code:    "GEN_135",
			Message: "No pending action found",
		},
		"OpenMinGEOpenMax": {
			Code:    "GEN_136",
			Message: "Open min amount cannot be greater than or equal to open max amount.",
		},
		"OpenMaxGEPinMax": {
			Code:    "GEN_137",
			Message: "Open max amount cannot be greater than or equal to pin max amount.",
		},
		"PinMinGEPinMax": {
			Code:    "GEN_138",
			Message: "Pin min amount cannot be greater than or equal to pin max amount.",
		},
		"PinMinLEOpenMin": {
			Code:    "GEN_139",
			Message: "Pin min amount cannot be less than or equal to open min amount.",
		},
		"PinMaxLEPinMin": {
			Code:    "GEN_140",
			Message: "Pin max amount must be greater than pin min amount.",
		},
		"OTPMinGEPinMin": {
			Code:    "GEN_141",
			Message: "OTP min amount cannot be greater than or equal to pin min amount.",
		},
		"TierAuthNotFound": {
			Code:    "GEN_142",
			Message: "Authentication tier not found.",
		},
		"FailedToGetAuthTier": {
			Code:    "GEN_143",
			Message: "Failed to retrieve authentication tier.",
		},
		"REQUIRED_TITLE": {
			Code:    "GEN_144",
			Message: "title is required",
		},
		"REQUIRED_DESCRIPTION": {
			Code:    "GEN_145",
			Message: "description is required",
		},
		"TITLE_TOO_LONG": {
			Code:    "GEN_146",
			Message: "title length is between 3 and 10 characters",
		},
		"DESCRIPTION_TOO_LONG": {
			Code:    "GEN_147",
			Message: "description length is between 3 and 100 characters",
		},
		"MISSING_CONFIG_KEY": {
			Code:    "GEN_148",
			Message: "Configuration key is missing.",
		},
		"ACTION_CODE_REQUIRED": {
			Code:    "GEN_149",
			Message: "ActionCode is required",
		},
		"INVALID_ACTION_STATUS": {
			Code:    "GEN_150",
			Message: "invalid Action Status",
		},
		"INVALID_REQUEST_ACTION": {
			Code:    "GEN_151",
			Message: "invalid RequestAction",
		},
		"UNSUPPORTED_REQUEST_ACTION": {
			Code:    "GEN_152",
			Message: "unsupported request action",
		},

		"FAILED_TO_GET_DEPARTMENT": {
			Code:    "GEN_152",
			Message: "failed to get department",
		},
		"KYC_LEVEL_REQUIRED": {
			Code:    "GEN_153",
			Message: "kyc_level is required",
		},
		"INVALID_KYC_LEVEL": {
			Code:    "GEN_154",
			Message: "kyc_level must be a number between 0 and 2",
		},
		"NO_DATA_PROVIDED_FOR_UPDATE": {
			Code:    "GEN_155",
			Message: "Update request must include at least one field to modify.",
		},
		"CONTENT_TYPE_MUST_BE_FORM": {
			Code:    "GEN_156",
			Message: "Content-Type must be multipart/form-data",
		},
		"CONTENT_TYPE_MUST_BE_JSON": {
			Code:    "GEN_157",
			Message: "Content-Type must be application/json",
		},
		"RESOURCE_ALREADY_ENABLED": {
			Code:    "GEN_158",
			Message: "Resource is already enabled",
		},
		"RESOURCE_ALREADY_DISABLED": {
			Code:    "GEN_159",
			Message: "Resource is already disabled",
		},
		"INVALID_PAYLOAD": {
			Code:    "GEN_160",
			Message: "The request payload is invalid or malformed. Please check the input structure and types.",
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
		"ACTION_NOT_FOUND_AFTER_UPDATE": {
			Code:    "ACT_404_UPD",
			Message: "action not found after update",
		},
		"ACTION_ALREADY_APPROVED": {
			Code:    "GEN_005",
			Message: "This action is already approved.",
		},

		"ACCOUNT_LOCKED": {
			Code:    "AUTH_028",
			Message: "account locked. please contact your addministrator",
		},
		"ACCESS_TOKEN_REQUIRED": {
			Code:    "AUTH_030",
			Message: "access token required",
		},
		"DEPARTMEN_REQUIRED": {
			Code:    "AUTH_031",
			Message: "department is required in context",
		},
		"INCOMPLETE_USER_INFO": {
			Code:    "AUTH_032",
			Message: "Incomplete user information",
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
		"BLOCKED_ACTION_USER_ALREADY_EXIST": {
			Code:    "AUTH_029",
			Message: "A pending block action already exists for this user",
		},
		"BRANCH_DISABLE_ACTION_ALREADY_EXISTS": {
			Code:    "BRN_008",
			Message: "A pending disable action already exists for this branch",
		},
		"BRANCH_DISABLE_MULTI_ACTION_ALREADY_EXISTS": {
			Code:    "BRN_010",
			Message: "A pending disable action already exists for one or more of the selected branches",
		},
		"BRANCH_BULK_ACTION_ALREADY_PROCESSED": {
			Code:    "BRN_011",
			Message: "This bulk action has already been processed for one or more branches",
		},
		"BRANCH_ID_REQUIRED": {
			Code:    "BRN_006",
			Message: "Branch ID is required.",
		},
		"BRANCH_REGION_AND_DISTRICT_REQUIRED": {
			Code:    "BRN_012",
			Message: "Region and district are required",
		},
		"BRANCH_REGION_AND_DISTRICT_MIN_LENGTH": {
			Code:    "BRN_014",
			Message: "Region and district must be at least 3 characters",
		},
		"FAILED_TO_FETCH_BRANCHES": {
			Code:    "BRN_013",
			Message: "Failed to fetch branches",
		},
		"INVALID_JSON_PAYLOAD": {
			Code:    "GEN_001",
			Message: "Invalid JSON body",
		},
		"BRANCH_NOT_FOUND": {
			Code:    "BRN_001",
			Message: "No branch found for the given region and district",
		},
		"UNHANDLED_SERVER_ERROR": {
			Code:    "GEN_004",
			Message: "Internal server error",
		},
		"BLOCK_DISTRICT_ALREADY_PROCESSED": {
			Code:    "DST_001",
			Message: "This district action has already been processed",
		},
		// Additional missing error from the other definition
		"ACCOUNT_UPDATE_FAILED": {
			Code:    "ACC_014",
			Message: "Failed to update account.",
		},
		"ACCOUNT_BLOCKED": {
			Code:    "ACC_014",
			Message: "The account is blocked.",
		},
		"REGION_CODE_AND_NAME_REQUIRED": {
			Code:    "ACC_015",
			Message: "Region code  is required",
		},
		"BLOCK_REGION_ALREADY_PROCESSED": {
			Code:    "DST_016",
			Message: "This region already blocked",
		},
		"BANK_ID_REQUIRED": {
			Code:    "DST_017",
			Message: "Bank id is required",
		},
		"FAILED_TO_GET_DISTRICT": {
			Code:    "DST_018",
			Message: "district not found",
		},
		"BRANCH_ALREADY_DISABLE": {
			Code:    "DST_019",
			Message: "branch code %v already disabled",
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
		"MISSING_OR_INVALID_IMAGE": {
			Code:    "FILE_006",
			Message: "Missing or invalid image. Could not parse form or file too large.",
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
		"CITY_NOT_FOUND": {
			Code:    "BRN_008",
			Message: "City not found",
		},
		"DISTRICT_ALREADY_BLOCKED": {
			Code:    "BRN_009",
			Message: "The district already blocked",
		},
		"DISTRICT_NOT_FOUND": {
			Code:    "BRN_010",
			Message: "District Not Found",
		},
		"REGION_NOT_FOUND": {
			Code:    "BRN_011",
			Message: "Region Not Found",
		},
		"BRANCH_ALREADY_BLOCKED": {
			Code:    "BRN_010",
			Message: "The branch already blocked",
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
			Message: "Bank data not found.",
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
		"BANK_ALREADY_EXISTS": {
			Code:    "BNK_006",
			Message: "Bank with this name already exists.",
		},
		"BANK_ALREADY_ENABLE": {
			Code:    "BNK_007",
			Message: "Bank  already  enabled",
		},
		"BANK_ALREADY_DISABLED": {
			Code:    "BNK_008",
			Message: "Bank  already  disabled",
		},
		"NAME_OF_BANK_ALREADY_EXIST": {
			Code:    "BNK_009",
			Message: "Bank with this name already exist",
		},
		"BIC_CODE_ALREADY_EXIST": {
			Code:    "BNK_010",
			Message: "Bank BIC CODE already exist",
		},
		"BANK_BIC_CODE_ALREADY_EXIST": {
			Code:    "BNK_011",
			Message: "Bank BIC CODE already exist",
		},
		"BANK_ALREADY_CREATED_WITH_THIS_PARAMETER": {
			Code:    "BNK_012",
			Message: "Bank already exist with this parameter",
		},
		"BANK_FETCH_FAILED": {
			Code:    "BNK_013",
			Message: "Bank fetch failed",
		},
		"BANK_NAME_ALREADY_EXIST": {
			Code:    "BNK_014",
			Message: "Bank Name already exist",
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
		"FAILED_TO_CREATE_ACTION": {
			Code:    "ACT_008",
			Message: "Unable to create action.",
		},
		"ACTION_NOT_PENDING": {
			Code:    "ACT_009",
			Message: "Action is not in pending status.",
		},
		"FAILED_TO_UPDATE_ACTION": {
			Code:    "ACT_010",
			Message: "Failed to update action.",
		}, "FAILED_TO_FETCH_ACTION": {
			Code:    "ACT_011",
			Message: "Failed to fetch action.",
		},
		"FAILED_TO_UPDATE_SERVICE": {
			Code:    "ACT_012",
			Message: "Failed to update service details.",
		},
		"FAILED_TO_UPDATE_CAP_MIN": {
			Code:    "ACT_013",
			Message: "Failed to update cap minimum amount.",
		},
		"MISSING_REJECT_REASON": {
			Code:    "ACT_014",
			Message: "Missing reason for rejection.",
		},
		"REJECT_REASON_TOO_SHORT": {
			Code:    "ACT_015",
			Message: "Rejection reason must be at between 30 to 100 characters long.",
		},
		"PENDING_ACTION_ALREADY_EXIST": {
			Code:    "ACT_016",
			Message: "A pending block action already exists for this city",
		},
		"CITY_ALREADY_BLOCKED": {
			Code:    "ACT_017",
			Message: "This city already blocked",
		},
		"PENDING_DISTRICT_ACTION_ALREADY_EXIST": {
			Code:    "ACT_017",
			Message: "A pending block action already exists for this district",
		},
	},
	User: ErrorGroup{
		"USER_STATUS_UPDATE_FAILED": {
			Code:    "USR_001",
			Message: "Failed to update user status.",
		},
	},
	Wallet: ErrorGroup{
		"WALLET_NOT_FOUND": {
			Code:    "WAL_001",
			Message: "Wallet not found.",
		},
		"WALLET_CREATION_FAILED": {
			Code:    "WAL_002",
			Message: "Failed to create wallet.",
		},
		"WALLET_UPDATE_FAILED": {
			Code:    "WAL_003",
			Message: "Failed to update wallet.",
		},
		"WALLET_DELETION_FAILED": {
			Code:    "WAL_004",
			Message: "Failed to delete wallet.",
		},
		"WALLET_BALANCE_INSUFFICIENT": {
			Code:    "WAL_005",
			Message: "Insufficient wallet balance.",
		},
		"WALLET_TRANSACTION_FAILED": {
			Code:    "WAL_006",
			Message: "Failed to process wallet transaction.",
		},
		"WALLET_NOT_ACTIVE": {
			Code:    "WAL_007",
			Message: "Wallet is not active.",
		},
		"WALLET_ALREADY_EXISTS": {
			Code:    "WAL_008",
			Message: "Wallet already exists.",
		},
		"WALLET_TYPE_NOT_SUPPORTED": {
			Code:    "WAL_009",
			Message: "Wallet type is not supported.",
		},
		"WALLET_LIMIT_EXCEEDED": {
			Code:    "WAL_010",
			Message: "Wallet limit exceeded.",
		},
		"WALLET_TRANSACTION_NOT_FOUND": {
			Code:    "WAL_011",
			Message: "Wallet transaction not found.",
		},
		"WALLET_TRANSACTION_ALREADY_EXISTS": {
			Code:    "WAL_012",
			Message: "Wallet transaction already exists.",
		},
		"WALLET_INFORMATION_ALREADY_EXISTST": {
			Code:    "WAL_014",
			Message: "One or more fields already exists",
		},
	},
	AD: ErrorGroup{
		"AD_NOT_FOUND": {
			Code:    "AD_001",
			Message: "Advert not found.",
		},
		"AD_CREATION_FAILED": {
			Code:    "AD_002",
			Message: "Failed to create advert.",
		},
		"AD_UPDATE_FAILED": {
			Code:    "AD_003",
			Message: "Failed to update advert.",
		},
		"AD_ALREADY_EXISTS": {
			Code:    "AD_004",
			Message: "Advert with this ID already exists.",
		},
		"AD_DELETION_FAILED": {
			Code:    "AD_005",
			Message: "Failed to delete advert.",
		},
		"START_DATE_REQUIRED": {
			Code:    "AD_006",
			Message: "started date is required",
		},
		"EXPIRE_DATE_REQUIRED": {
			Code:    "AD_006",
			Message: "end date is required",
		},
	},
}

func (e ErrorDefinition) Error() string {
	return e.Message
}
