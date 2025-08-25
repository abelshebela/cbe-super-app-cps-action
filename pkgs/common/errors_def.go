package common

import "net/http"

// HTTP Status Code Constants
const (
	// 2xx Success
	StatusOK        = http.StatusOK        // 200
	StatusCreated   = http.StatusCreated   // 201
	StatusAccepted  = http.StatusAccepted  // 202
	StatusNoContent = http.StatusNoContent // 204

	// 4xx Client Errors
	StatusBadRequest           = http.StatusBadRequest           // 400
	StatusUnauthorized         = http.StatusUnauthorized         // 401
	StatusForbidden            = http.StatusForbidden            // 403
	StatusNotFound             = http.StatusNotFound             // 404
	StatusMethodNotAllowed     = http.StatusMethodNotAllowed     // 405
	StatusConflict             = http.StatusConflict             // 409
	StatusUnprocessableEntity  = http.StatusUnprocessableEntity  // 422
	StatusTooManyRequests      = http.StatusTooManyRequests      // 429
	StatusUnsupportedMediaType = http.StatusUnsupportedMediaType // 415

	// 5xx Server Errors
	StatusInternalServerError = http.StatusInternalServerError // 500
	StatusNotImplemented      = http.StatusNotImplemented      // 501
	StatusBadGateway          = http.StatusBadGateway          // 502
	StatusServiceUnavailable  = http.StatusServiceUnavailable  // 503
)

type ErrorDefinition struct {
	Code    string
	Message string
	Status  int
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
	District    ErrorGroup
	City        ErrorGroup
	Department  ErrorGroup
	Bank        ErrorGroup
	Action      ErrorGroup
	Wallet      ErrorGroup
	AD          ErrorGroup
	BulkService ErrorGroup
	Permission  ErrorGroup
	MiniApp     ErrorGroup
	Event       ErrorGroup
	Service     ErrorGroup
	Donation    ErrorGroup
}

var DefineError = ErrorDefinitions{
	General: ErrorGroup{
		"CONFLICT_KEY": {
			Code:    "GEN_001",
			Status:  StatusConflict,
			Message: "Duplicate key error: The specified field already exists.",
		},
		"INVALID_ID": {
			Code:    "GEN_002",
			Status:  StatusBadRequest,
			Message: "Invalid ID provided.",
		},
		"INVALID_ID_FORMAT": {
			Code:    "GEN_018",
			Status:  StatusBadRequest,
			Message: "Invalid ID format provided.",
		},
		"INVALID_JSON_PAYLOAD": {
			Code:    "GEN_003",
			Status:  StatusBadRequest,
			Message: "Invalid JSON payload.",
		},
		"UNHANDLED_SERVER_ERROR": {
			Code:    "GEN_004",
			Status:  StatusInternalServerError,
			Message: "An unexpected server error occurred.",
		},
		"INVALID_TOKEN_FORMAT": {
			Code:    "GEN_005",
			Status:  StatusBadRequest,
			Message: "Invalid token format.",
		},
		"INVALID_TOKEN": {
			Code:    "GEN_006",
			Status:  StatusUnauthorized,
			Message: "Invalid token.",
		},
		"EXPIRED_TOKEN": {
			Code:    "GEN_007",
			Status:  StatusUnauthorized,
			Message: "Token has expired.",
		},
		"ENCRYPTED_PAYLOAD_REQUIRED": {
			Code:    "GEN_008",
			Status:  StatusBadRequest,
			Message: "Encrypted payload is required.",
		},
		"INVALID_SERVER_CONFIGURATION": {
			Code:    "GEN_009",
			Status:  StatusInternalServerError,
			Message: "Invalid server configuration",
		},
		"AUTH_HEADER_MISSING": {
			Code:    "GEN_010",
			Status:  StatusUnauthorized,
			Message: "Authorization header is missing. please provide a valid bearer token.",
		},
		"INVALID_KEY_CONFIGURATION": {
			Code:    "GEN_011",
			Status:  StatusInternalServerError,
			Message: "Invalid key configuration.",
		},
		"INVALID_ENCRYPTED_PAYLOAD": {
			Code:    "GEN_012",
			Status:  StatusBadRequest,
			Message: "The provided encrypted payload is invalid.",
		},
		"DECRYPTION_ERROR": {
			Code:    "GEN_013",
			Status:  StatusBadRequest,
			Message: "Decryption error occurred. please check your key configuration.",
		},
		"SERVER_KEYS_NOT_CONFIGURED": {
			Code:    "GEN_014",
			Status:  StatusInternalServerError,
			Message: "server error occurred. please try again later or contact support.",
		},
		"INVALID_INPUT": {
			Code:    "GEN_015",
			Status:  StatusBadRequest,
			Message: "the provided input is invalid.",
		},
		"UNAUTHORIZED": {
			Code:    "GEN_016",
			Status:  StatusUnauthorized,
			Message: "missing authenticated user.",
		},
		"USER_REALM_NOT_FOUND": {
			Code:    "GEN_017",
			Status:  StatusNotFound,
			Message: "user is not permission to do this action.",
		},
		"ACTION_NOT_ALLOWED": {
			Code:    "GEN_018",
			Status:  StatusForbidden,
			Message: "user does not have permission to do this action.",
		},
		"NOT_FOUND": {
			Code:    "GEN_019",
			Status:  StatusNotFound,
			Message: "Resource not found.",
		},
		"EMPTY_ORG": {
			Code:    "GEN_020",
			Status:  StatusBadRequest,
			Message: "Cant respond for empty organization id",
		},
		"INVALID_REQ": {
			Code:    "GEN_021",
			Status:  StatusBadRequest,
			Message: "Invalid Request body",
		},
		"REG_FRST": {
			Code:    "GEN_022",
			Status:  StatusForbidden,
			Message: "please register on the Super APP and visit your nearest branch for Verification Code",
		},
		"WEAK_PIN": {
			Code:    "GEN_023",
			Status:  StatusBadRequest,
			Message: "weak PIN used: please use a strong combination",
		},
		"WAIT_FOR_PREVIOUS_OTP_EXPIRATION": {
			Code:    "GEN_024",
			Status:  StatusTooManyRequests,
			Message: "wait until the pervious otp expired",
		},
		"INCOMPLETE_USER_INFO": {
			Code:    "GEN_025",
			Status:  StatusBadRequest,
			Message: "Incomplete user information",
		},
		"INVALID_ACTION_TYPE": {
			Code:    "GEN_026",
			Status:  StatusBadRequest,
			Message: "Invalid action type",
		},
		"COULD_NOT_UNLINK_DEVICE": {
			Code:    "GEN_027",
			Status:  StatusInternalServerError,
			Message: "Could not unlink device",
		},
		"ACTION_NOT_FOUND": {
			Code:    "GEN_028",
			Status:  StatusNotFound,
			Message: "Action not found.",
		},
		"ERROR_CHANGING_PIN": {
			Code:    "GEN_029",
			Status:  StatusBadRequest,
			Message: "Error changing PIN.",
		},
		"MISSING_REQUIRED_FIELDS": {
			Code:    "GEN_030",
			Status:  StatusBadRequest,
			Message: "Missing required fields in request",
		},
		"MISSING_OTP": {
			Code:    "GEN_031",
			Status:  StatusBadRequest,
			Message: "OTP code is required",
		},
		"INVALID_INPUT_PARAMETERS": {
			Code:    "GEN_032",
			Status:  StatusBadRequest,
			Message: "Invalid input parameters provided",
		},
		"PENDING_REQUEST_EXISTS": {
			Code:    "GEN_033",
			Status:  StatusConflict,
			Message: "You have a pending request for this action.",
		},
		"SAME_PIN": {
			Code:    "GEN_034",
			Status:  StatusBadRequest,
			Message: "New PIN cannot be the same as current PIN",
		},
		"GENERAL_DB_QUERY_FAILED": {
			Code:    "GEN_035",
			Status:  StatusInternalServerError,
			Message: "Database query failed.",
		},
		"ERROR_SETTING_PIN": {
			Code:    "GEN_036",
			Status:  StatusInternalServerError,
			Message: "Failed to set PIN due to system error",
		},
		"OTP_CREATION_FAILED": {
			Code:    "GEN_037",
			Status:  StatusInternalServerError,
			Message: "Failed to create OTP",
		},
		"DEVICE_NOT_FOUND": {
			Code:    "GEN_038",
			Status:  StatusNotFound,
			Message: "Device not found in system",
		},
		"DEVICE_FOUND": {
			Code:    "GEN_039",
			Status:  StatusOK,
			Message: "Device found and registered",
		},
		"INVALID_PIN": {
			Code:    "GEN_040",
			Status:  StatusBadRequest,
			Message: "Invalid PIN provided.",
		},
		"PIN_RESET_SESSION_EXPIRED": {
			Code:    "GEN_041",
			Status:  StatusBadRequest,
			Message: "PIN reset session has expired",
		},
		"PIN_RESET_SESSION_NOT_FOUND": {
			Code:    "GEN_042",
			Status:  StatusNotFound,
			Message: "PIN reset session not found",
		},
		"PIN_RESET_SESSION_ALREADY_VERIFIED": {
			Code:    "GEN_043",
			Status:  StatusBadRequest,
			Message: "PIN reset session already verified",
		},
		"FAILED_TO_GET_FEEDBACK_DATA": {
			Code:    "GEN_044",
			Status:  StatusInternalServerError,
			Message: "failed to get feedback data",
		},
		"MAXIMUM_AMOUNT_REQUIRED_FOR_OPEN_METHOD": {
			Code:    "GEN_045",
			Status:  StatusBadRequest,
			Message: "Maximum amount is required for open method",
		},
		"FAILED_TO_GET_FEEDBACK_COUNTS": {
			Code:    "GEN_046",
			Status:  StatusInternalServerError,
			Message: "failed to get feedback total counts",
		},
		"EITHER_MINIMUM_OR_MAXIMUM_AMOUNT_REQUIRED_FOR_PIN_METHOD": {
			Code:    "GEN_047",
			Status:  StatusBadRequest,
			Message: "Either minimum or maximum amount are required for PIN method",
		},
		"INVALID_OBJECT_ID_FORMAT": {
			Code:    "GEN_048",
			Status:  StatusBadRequest,
			Message: "Invalid object ID format",
		},
		"AUTH_TIER_NOT_FOUND": {
			Code:    "GEN_049",
			Status:  StatusNotFound,
			Message: "Auth tier not found",
		},
		"AUTH_TIER_ALREADY_EXISTS": {
			Code:    "GEN_050",
			Status:  StatusBadRequest,
			Message: "Auth tier already exists",
		},
		"AUTH_TIER_VALIDATION_FAILED": {
			Code:    "GEN_051",
			Status:  StatusBadRequest,
			Message: "Auth tier validation failed",
		},
		"SERVICE_TYPE_CANNOT_BE_EMPTY": {
			Code:    "GEN_052",
			Status:  StatusBadRequest,
			Message: "Service type cannot be empty",
		},
		"AUTH_TIER_INSERT_FAILED": {
			Code:    "GEN_053",
			Status:  StatusInternalServerError,
			Message: "Auth tier insert failed",
		},
		"AUTH_TIER_NOT_FOUND_FOR_ID": {
			Code:    "GEN_054",
			Status:  StatusNotFound,
			Message: "Auth tier not found for the provided ID",
		},
		"AUTH_TIER_FETCH_FAILED": {
			Code:    "GEN_055",
			Status:  StatusInternalServerError,
			Message: "Auth tier fetch failed",
		},
		"AUTH_TIER_APPROVE_FAILED": {
			Code:    "GEN_056",
			Status:  StatusInternalServerError,
			Message: "Auth tier approve failed",
		},
		"AUTH_TIER_REJECT_FAILED": {
			Code:    "GEN_057",
			Status:  StatusInternalServerError,
			Message: "Auth tier reject failed",
		},
		"AUTH_TIER_NOT_FOUND_FOR_ID_AND_DEPARTMENT": {
			Code:    "GEN_058",
			Status:  StatusNotFound,
			Message: "Auth tier not found for the provided ID and department",
		},
		"AUTH_TIER_FETCH_FAILED_FOR_ID_AND_DEPARTMENT": {
			Code:    "GEN_059",
			Status:  StatusInternalServerError,
			Message: "Auth tier fetch failed for the provided ID and department",
		},
		"AUTH_TIER_APPROVE_FAILED_FOR_ID_AND_DEPARTMENT": {
			Code:    "GEN_060",
			Status:  StatusInternalServerError,
			Message: "Auth tier approve failed for the provided ID and department",
		},
		"AUTH_TIER_REJECT_FAILED_FOR_ID_AND_DEPARTMENT": {
			Code:    "GEN_061",
			Status:  StatusInternalServerError,
			Message: "Auth tier reject failed for the provided ID and department",
		},
		"FAILED_TO_UPDATE_USER": {
			Code:    "GEN_062",
			Status:  StatusInternalServerError,
			Message: "failed to update user",
		},
		"AUTH_TIER_INSERT_FAILED_FOR_ID_AND_DEPARTMENT": {
			Code:    "GEN_063",
			Status:  StatusInternalServerError,
			Message: "Auth tier insert failed for the provided ID and department",
		},
		"ACTION_ID_IS_REQUIRED": {
			Code:    "GEN_064",
			Status:  StatusBadRequest,
			Message: "action id is required",
		},
		"FAILED_TO_UPDATE_PIN_AUTH_TIER": {
			Code:    "GEN_065",
			Status:  StatusInternalServerError,
			Message: "Failed to update pin auth tier",
		},
		"FAILED_TO_UPDATE_OTP_AND_PIN_AUTH_TIER": {
			Code:    "GEN_066",
			Status:  StatusInternalServerError,
			Message: "Failed to update OTP and pin auth tier",
		},
		"FAILED_TO_UPDATE_OTP_AND_OPEN_AUTH_TIER": {
			Code:    "GEN_067",
			Status:  StatusInternalServerError,
			Message: "Failed to update OTP and open auth tier",
		},
		"FAILED_TO_UPDATE_OPEN_AUTH_TIER_FOR_ID": {
			Code:    "GEN_068",
			Status:  StatusInternalServerError,
			Message: "Failed to update open auth tier for the provided ID",
		},
		"MIN_AMOUNT_CANNOT_BE_GREATER_THAN_PIN_MIN_AMOUNT": {
			Code:    "GEN_069",
			Status:  StatusBadRequest,
			Message: "Min amount cannot be greater than or equal to pin min amount",
		},
		"PIN_AUTHIER_NOT_FOUND": {
			Code:    "GEN_070",
			Status:  StatusNotFound,
			Message: "Pin authier not found",
		},
		"PIN_MIN_AMOUNT_CANNOT_BE_LESS_THAN_OPEN_MIN_AMOUNT": {
			Code:    "GEN_071",
			Status:  StatusBadRequest,
			Message: "Pin min amount cannot be less than or equal to open min amount",
		},
		"MAX_PIN_AMOUNT_SHOULD_BE_GREATER_THAN_PIN_MIN_AMOUNT": {
			Code:    "GEN_072",
			Status:  StatusBadRequest,
			Message: "Max pin amount should be greater than pin min amount",
		},
		"OPEN_AUTHIER_NOT_FOUND": {
			Code:    "GEN_073",
			Status:  StatusNotFound,
			Message: "Open authier not found",
		},
		"INVALID_ACTION_DATA": {
			Code:    "GEN_074",
			Status:  StatusBadRequest,
			Message: "Invalid action data.",
		},
		"PENDING_CPS_ACTION_PRESENT": {
			Code:    "GEN_075",
			Status:  StatusConflict,
			Message: "Pending cps action present",
		},
		"FAILED_TO_GET_FAYDA_ACCOUNT": {
			Code:    "GEN_076",
			Status:  StatusInternalServerError,
			Message: "Failed to get fayda account",
		},
		"FAILED_TO_GET_CUSTOMER_ACCOUNT": {
			Code:    "GEN_077",
			Status:  StatusInternalServerError,
			Message: "Failed to get customer account",
		},
		"FAILED_TO_CREATE_CPS_ACTION": {
			Code:    "GEN_078",
			Status:  StatusInternalServerError,
			Message: "Failed to create cps action",
		},
		"AT_LEAST_ONE_TIER_IS_REQUIRED": {
			Code:    "GEN_079",
			Status:  StatusBadRequest,
			Message: "At least one tier is required",
		},
		"THE_FIRST_TIER_MINIMUM_MUST_START_FROM_0": {
			Code:    "GEN_080",
			Status:  StatusBadRequest,
			Message: "The first tier's minimum must start from 0",
		},
		"TIER_MINIMUM_MUST_EQUAL_TO_PREVIOUS_TIER_MAXIMUM": {
			Code:    "GEN_081",
			Status:  StatusBadRequest,
			Message: "Tier minimum must equal the previous tier's maximum",
		},
		"TIER_MAXIMUM_MUST_BE_GREATER_THAN_PREVIOUS_TIER_MAXIMUM": {
			Code:    "GEN_082",
			Status:  StatusBadRequest,
			Message: "Tier maximum must be greater than the previous tier's maximum",
		},
		"ABOVE_AMOUNT_MUST_MATCH_LAST_TIER_MAXIMUM": {
			Code:    "GEN_083",
			Status:  StatusBadRequest,
			Message: "Above amount must match the last tier's maximum",
		},
		"SERVICE_CODE_CANNOT_BE_EMPTY": {
			Code:    "GEN_084",
			Status:  StatusBadRequest,
			Message: "Service code cannot be empty",
		},
		"SERVICE_NAME_CANNOT_BE_EMPTY": {
			Code:    "GEN_085",
			Status:  StatusBadRequest,
			Message: "Service name cannot be empty",
		},
		"NO_LINKED_DEVICES": {
			Code:    "GEN_086",
			Status:  StatusNotFound,
			Message: "User has no linked devices",
		},
		"FAILED_TO_GET_SERVICE": {
			Code:    "GEN_087",
			Status:  StatusInternalServerError,
			Message: "Failed to get service",
		},
		"SERVICE_NOT_FOUND": {
			Code:    "GEN_088",
			Status:  StatusNotFound,
			Message: "Service not found",
		},
		"DEVICE_LOOKUP_FAILED": {
			Code:    "GEN_089",
			Status:  StatusInternalServerError,
			Message: "Device lookup operation failed",
		},
		"FAILED_TO_UPDATE_SERVICE_FEE": {
			Code:    "GEN_090",
			Status:  StatusInternalServerError,
			Message: "Failed to update service fee",
		},
		"FAILED_TO_REJECT_SERVICE_FEE_UPDATE": {
			Code:    "GEN_091",
			Status:  StatusInternalServerError,
			Message: "Failed to reject service fee update",
		},
		"AUTH_TIER_UPDATE_FAILED": {
			Code:    "GEN_092",
			Status:  StatusInternalServerError,
			Message: "Auth tier update failed",
		},
		"ACTION_IS_NOT_IN_PENDING_STATUS": {
			Code:    "GEN_093",
			Status:  StatusBadRequest,
			Message: "Action is not in pending status",
		},
		"FAILED_TO_UPDATE_SERVICE_DETAILS": {
			Code:    "GEN_094",
			Status:  StatusInternalServerError,
			Message: "Failed to update service details",
		},
		"DEVICE_ID_REQUIRED": {
			Code:    "GEN_095",
			Status:  StatusBadRequest,
			Message: "Device UUID is required",
		},
		"TOKEN_GENERATION_FAILED": {
			Code:    "GEN_096",
			Status:  StatusInternalServerError,
			Message: "Failed to generate token",
		},
		"FAILED_TO_GET_FEEDBACK": {
			Code:    "GEN_097",
			Status:  StatusInternalServerError,
			Message: "failed to get feedback",
		},
		"PIN_IN_HISTORY": {
			Code:    "GEN_098",
			Status:  StatusBadRequest,
			Message: "PIN has been used recently and cannot be reused",
		},
		"MINIMUM_AMOUNT_REQUIRED_FOR_OTP_AND_PIN_METHOD": {
			Code:    "GEN_099",
			Status:  StatusBadRequest,
			Message: "Minimum amount is required for OTP and PIN method",
		},
		"FAILED_TO_CONVERT_ID": {
			Code:    "GEN_100",
			Status:  StatusInternalServerError,
			Message: "failed to convert ID",
		},
		"FAIL_TO_FIND_USER": {
			Code:    "GEN_101",
			Status:  StatusInternalServerError,
			Message: "failed to find user",
		},
		"DUPLICATE_KEY": {
			Code:    "GEN_102",
			Status:  StatusConflict,
			Message: "duplicate key error",
		},
		"INVALID_SERVICE_DATA_IN_ACTION": {
			Code:    "GEN_103",
			Status:  StatusBadRequest,
			Message: "Invalid service data in action",
		},
		"FAIL_TO_FEATCH_CPS_USER": {
			Code:    "GEN_104",
			Status:  StatusInternalServerError,
			Message: "fail to featch cps user",
		},
		"AUTH_TIER_UPDATE_FAILED_FOR_ID_AND_DEPARTMENT": {
			Code:    "GEN_105",
			Status:  StatusInternalServerError,
			Message: "Auth tier update failed for the provided ID and department",
		},
		"FAILED_TO_UPDATE_CPS_ACTION": {
			Code:    "GEN_106",
			Status:  StatusInternalServerError,
			Message: "failed to update cps action",
		},
		"OPEN_TIER_MAX_AMOUNT_CANNOT_BE_GREATER_THAN_PIN_MAX_AMOUNT": {
			Code:    "GEN_107",
			Status:  StatusBadRequest,
			Message: "Open tier max amount can not be greater than pin max amount",
		},
		"FAILED_TO_UPDATE_OPEN_AUTH_TIER": {
			Code:    "GEN_108",
			Status:  StatusInternalServerError,
			Message: "Failed to update open auth tier",
		},
		"MISSING_REQUIRED_HEADERS": {
			Code:    "GEN_109",
			Status:  StatusBadRequest,
			Message: "Missing required headers for device lookup",
		},
		"USER_CODE_IS_REQUIRED": {
			Code:    "GEN_110",
			Status:  StatusBadRequest,
			Message: "USER CODE IS REQUIRED",
		},
		"FAILED_TO_FIND_CPS_ACTION": {
			Code:    "GEN_111",
			Status:  StatusInternalServerError,
			Message: "Failed to find cps action",
		},
		"CPS_USER_NOT_FOUND": {
			Code:    "GEN_112",
			Status:  StatusNotFound,
			Message: "cps user not found",
		},
		"REQUIRED_FIELDS_MISSING": {
			Code:    "GEN_113",
			Status:  StatusBadRequest,
			Message: "One or more required fields are missing.",
		},
		"VALIDATION_RULE_ID_MISMATCH": {
			Code:    "GEN_114",
			Status:  StatusBadRequest,
			Message: "Validation rule ID does not match the expected value.",
		},
		"FAILED_TO_UPDATE_VALIDATION_RULE": {
			Code:    "GEN_115",
			Status:  StatusInternalServerError,
			Message: "Failed to update the validation rule.",
		},
		"FAILED_TO_UNMARSHAL_CURRENT_ACTION": {
			Code:    "GEN_116",
			Status:  StatusInternalServerError,
			Message: "Failed to parse the current action data.",
		},
		"CURRENT_ACTION_INVALID_TYPE": {
			Code:    "GEN_117",
			Status:  StatusBadRequest,
			Message: "Current action has an invalid data type.",
		},
		"CURRENT_ACTION_NIL": {
			Code:    "GEN_118",
			Status:  StatusBadRequest,
			Message: "Current action is missing or nil.",
		},
		"ACTION_NOT_PENDING": {
			Code:    "GEN_119",
			Status:  StatusBadRequest,
			Message: "Action is not in a pending state.",
		},
		"CHECKER_ID_EMPTY": {
			Code:    "GEN_120",
			Status:  StatusBadRequest,
			Message: "Checker ID is required but was not provided.",
		},
		"ACTION_ID_EMPTY": {
			Code:    "GEN_121",
			Status:  StatusBadRequest,
			Message: "Action ID is required but was not provided.",
		},
		"FAILED_TO_MARSHAL_CURRENT_ACTION": {
			Code:    "GEN_122",
			Status:  StatusInternalServerError,
			Message: "Failed to serialize the current action data.",
		},
		"FAILED_TO_MARSHAL_PREVIOUS_ACTION": {
			Code:    "GEN_123",
			Status:  StatusInternalServerError,
			Message: "Failed to serialize the previous action data.",
		},
		"PENDING_ACTION_EXISTS": {
			Code:    "GEN_124",
			Status:  StatusConflict,
			Message: "A pending action already exist.",
		},
		"FAILED_TO_FETCH_PENDING_ACTIONS": {
			Code:    "GEN_125",
			Status:  StatusInternalServerError,
			Message: "Failed to fetch pending actions for this validation rule.",
		},
		"FAILED_TO_GET_AD": {
			Code:    "GEN_213",
			Status:  StatusInternalServerError,
			Message: "Failed to get ad.",
		},
		"ADVERT_NOT_FOUND": {
			Code:    "GEN_214",
			Status:  StatusNotFound,
			Message: "advert not found",
		},
		"FAILED_TO_CAST_ACTION_DATA": {
			Code:    "GEN_215",
			Status:  StatusInternalServerError,
			Message: "failed to cast action data to ad request",
		},
		"FAILED_TO_CHECK_AD_BUCKET": {
			Code:    "GEN_216",
			Status:  StatusInternalServerError,
			Message: "failed to check ad bucket",
		},
		"FAILED_TO_CREATE_AD_BUCKET": {
			Code:    "GEN_217",
			Status:  StatusInternalServerError,
			Message: "failed to create ad bucket",
		},
		"FAILED_TO_OPEN_FILE": {
			Code:    "GEN_218",
			Status:  StatusInternalServerError,
			Message: "failed to create ad bucket",
		},
		"FAILED_TO_SAVE_OBJECT_TO_MINIO": {
			Code:    "GEN_219",
			Status:  StatusInternalServerError,
			Message: "failed to save object to MinIO",
		},
		"TIERS_REQUIRED": {
			Code:    "GEN_126",
			Status:  StatusBadRequest,
			Message: "At least one tier is required.",
		},
		"TIERS_FIRST_MIN_ZERO": {
			Code:    "GEN_127",
			Status:  StatusBadRequest,
			Message: "The first tier's minimum must start from 0.",
		},
		"TIERS_MIN_MUST_EQUAL_PREV_MAX": {
			Code:    "GEN_127",
			Status:  StatusBadRequest,
			Message: "Tier minimum must equal the previous tier's maximum.",
		},
		"TIERS_MAX_MUST_INCREASE": {
			Code:    "GEN_129",
			Status:  StatusBadRequest,
			Message: "Tier maximum must be greater than the previous tier's maximum.",
		},
		"TIERS_ABOVE_AMOUNT_MISMATCH": {
			Code:    "GEN_130",
			Status:  StatusBadRequest,
			Message: "Above amount must match the last tier's maximum value.",
		},
		"NOT_IMPLEMENTED": {
			Code:    "GEN_131",
			Status:  StatusNotImplemented,
			Message: "Not implemented.",
		},
		"INPUT_TOO_LONG": {
			Code:    "GEN_132",
			Status:  StatusBadRequest,
			Message: "Input exceeds maximum allowed length.",
		},
		"INPUT_INVALID_CHARACTERS": {
			Code:    "GEN_133",
			Status:  StatusBadRequest,
			Message: "Input contains invalid characters.",
		},
		"PAGE_NOT_FOUND": {
			Code:    "GEN_134",
			Status:  StatusNotFound,
			Message: "Page not found",
		},
		"FAILED_TO_GET_AVATAR": {
			Code:    "GEN_135",
			Status:  StatusNotFound,
			Message: "failed to get avatar",
		},
		"PENDING_REQUEST_CHECK_FAILED_FOR_CREATE_PERMISSION": {
			Code:    "GEN_136",
			Status:  StatusConflict,
			Message: "Permission group already exists with name",
		},
		"PERMISSION_GROUP_ALREADY_EXIXTS": {
			Code:    "GEN_137",
			Status:  StatusConflict,
			Message: "Permission group already exists with name",
		},
		"FAILED_TO_UPDATE_USER_DATA": {
			Code:    "GEN_132",
			Status:  StatusBadRequest,
			Message: "failed to update user data",
		},
		"FAILED_TO_INSERT_CPS_ACTION": {
			Code:    "GEN_133",
			Status:  StatusBadRequest,
			Message: "failed to insert cps action",
		},
		"ACCOUNT_ALREADY_DISABLED": {
			Code:    "GEN_134",
			Status:  StatusConflict,
			Message: "Account is already disabled",
		},
		"NO_PENDING_ACTION_FOUND": {
			Code:    "GEN_135",
			Status:  StatusNotFound,
			Message: "No pending action found",
		},
		"OpenMinGEOpenMax": {
			Code:    "GEN_136",
			Status:  StatusBadRequest,
			Message: "Open min amount cannot be greater than or equal to open max amount.",
		},
		"OpenMaxGEPinMax": {
			Code:    "GEN_137",
			Status:  StatusBadRequest,
			Message: "Open max amount cannot be greater than or equal to pin max amount.",
		},
		"PinMinGEPinMax": {
			Code:    "GEN_138",
			Status:  StatusBadRequest,
			Message: "Pin min amount cannot be greater than or equal to pin max amount.",
		},
		"PinMinLEOpenMin": {
			Code:    "GEN_139",
			Status:  StatusBadRequest,
			Message: "Pin min amount cannot be less than or equal to open min amount.",
		},
		"PinMaxLEPinMin": {
			Code:    "GEN_140",
			Status:  StatusBadRequest,
			Message: "Pin max amount must be greater than pin min amount.",
		},
		"OTPMinGEPinMin": {
			Code:    "GEN_141",
			Status:  StatusBadRequest,
			Message: "OTP min amount cannot be greater than or equal to pin min amount.",
		},
		"TierAuthNotFound": {
			Code:    "GEN_142",
			Status:  StatusNotFound,
			Message: "Authentication tier not found.",
		},
		"FailedToGetAuthTier": {
			Code:    "GEN_143",
			Status:  StatusInternalServerError,
			Message: "Failed to retrieve authentication tier.",
		},
		"REQUIRED_TITLE": {
			Code:    "GEN_144",
			Status:  StatusBadRequest,
			Message: "title is required",
		},
		"REQUIRED_DESCRIPTION": {
			Code:    "GEN_145",
			Status:  StatusBadRequest,
			Message: "description is required",
		},
		"TITLE_TOO_LONG": {
			Code:    "GEN_146",
			Status:  StatusBadRequest,
			Message: "title length is between 3 and 10 characters",
		},
		"DESCRIPTION_TOO_LONG": {
			Code:    "GEN_147",
			Status:  StatusBadRequest,
			Message: "description length is between 3 and 100 characters",
		},
		"MISSING_CONFIG_KEY": {
			Code:    "GEN_148",
			Status:  StatusBadRequest,
			Message: "Configuration key is missing.",
		},
		"ACTION_CODE_REQUIRED": {
			Code:    "GEN_149",
			Status:  StatusBadRequest,
			Message: "ActionCode is required",
		},
		"INVALID_ACTION_STATUS": {
			Code:    "GEN_150",
			Status:  StatusBadRequest,
			Message: "invalid Action Status",
		},
		"INVALID_REQUEST_ACTION": {
			Code:    "GEN_151",
			Status:  StatusBadRequest,
			Message: "invalid RequestAction",
		},
		"UNSUPPORTED_REQUEST_ACTION": {
			Code:    "GEN_152",
			Status:  StatusBadRequest,
			Message: "unsupported request action",
		},

		"FAILED_TO_GET_DEPARTMENT": {
			Code:    "GEN_152",
			Status:  StatusBadRequest,
			Message: "failed to get department",
		},
		"KYC_LEVEL_REQUIRED": {
			Code:    "GEN_153",
			Status:  StatusBadRequest,
			Message: "kyc_level is required",
		},
		"INVALID_KYC_LEVEL": {
			Code:    "GEN_154",
			Status:  StatusBadRequest,
			Message: "kyc_level must be a number between 0 and 2",
		},

		"UNEXPECTED_DATABASE_ERROR": {
			Code:    "GEN_155",
			Status:  StatusInternalServerError,
			Message: "An unexpected database error occurred. Please try again later.",
		},
		"USER_CODE_ALREADY_EXIST": {
			Code:    "GEN_156",
			Status:  StatusConflict,
			Message: "User code already exists. Please choose a different user code.",
		},
		"PHONE_NUMBER_EXISTS": {
			Code:    "GEN_157",
			Status:  StatusConflict,
			Message: "Phone number already exists. Please use a different phone number.",
		},
		"EMAIL_ALREADY_EXISTS": {
			Code:    "GEN_158",
			Status:  StatusConflict,
			Message: "Email already exists. Please use a different email address.",
		},
		"USERNAME_ALREADY_EXISTS": {
			Code:    "GEN_159",
			Status:  StatusConflict,
			Message: "Username already exists. Please choose a different username.",
		},
		"MAKER_OR_CHECKER": {
			Code:    "GEN_160",
			Status:  StatusBadRequest,
			Message: "User role can either be 'maker' or 'checker'",
		},
		"COLOR_ALREADY_EXISTED": {
			Code:    "GEN_160",
			Status:  StatusConflict,
			Message: "Color already Existed",
		},
		"NO_DATA_PROVIDED_FOR_UPDATE": {
			Code:    "GEN_161",
			Status:  StatusBadRequest,
			Message: "Update request must include at least one field to modify",
		},
		"CONTENT_TYPE_MUST_BE_FORM": {
			Code:    "GEN_162",
			Status:  StatusUnsupportedMediaType,
			Message: "Content-Type must be multipart/form-data",
		},
		"CONTENT_TYPE_MUST_BE_JSON": {
			Code:    "GEN_163",
			Status:  StatusUnsupportedMediaType,
			Message: "Content-Type must be application/json",
		},
		"RESOURCE_ALREADY_ENABLED": {
			Code:    "GEN_164",
			Status:  StatusConflict,
			Message: "Resource is already enabled",
		},
		"RESOURCE_ALREADY_DISABLED": {
			Code:    "GEN_165",
			Status:  StatusConflict,
			Message: "Resource is already disabled",
		},
		"INVALID_ACTION_FORMAT": {
			Code:    "GEN_166",
			Status:  StatusBadRequest,
			Message: "invalid action data format",
		},
		"MISSING_ICON": {
			Code:    "GEN_167",
			Status:  StatusBadRequest,
			Message: "missing icon_id in action data",
		},
		"FAILED_TO_AUTHORIZE": {
			Code:    "GEN_168",
			Status:  StatusInternalServerError,
			Message: "failed to authorize icon creation",
		},
		"MISSING_COLOR_ID": {
			Code:    "GEN_169",
			Status:  StatusBadRequest,
			Message: "missing color id in action data",
		},
		"FAILED_COLOR_UPDATE": {
			Code:    "GEN_170",
			Status:  StatusBadRequest,
			Message: "color update failed",
		},
		"MISSING_ACTION_DATA": {
			Code:    "GEN_171",
			Status:  StatusBadRequest,
			Message: "missing action data",
		},
		"CURENT_MAX_CAN_NOT_BE_GRETER_THAN_NEXT": {
			Code:    "GEN_166",
			Status:  StatusBadRequest,
			Message: "current max can not be greter than",
		},
		"INVALID_IMG_FORMAT": {
			Code:    "GEN_167",
			Status:  StatusBadRequest,
			Message: "Only image file formats are allowed",
		},
		"INVALID_OBJECT_ID": {
			Code:    "GEN_168",
			Status:  StatusInternalServerError,
			Message: "invalid object id",
		},
		"TIER_CAN_BE_UPDATED": {
			Code:    "GEN_169",
			Status:  StatusBadRequest,
			Message: "tier can't be updated",
		},
		"MAX_NOT_BE_LESS": {
			Code:    "GEN_168",
			Status:  StatusInternalServerError,
			Message: "max length not be less than min length",
		},
		"METHOD_NOT_ALLOWED": {
			Code:    "GEN_169",
			Status:  StatusBadRequest,
			Message: "Method not allowed",
		},
		"RESOURCE_INFORMATION_ALREADY_EXISTS": {
			Code:    "GEN_170",
			Status:  StatusConflict,
			Message: "One or more fields already exists",
		},
		"MERCHANT_NOT_FOUND": {
			Code:    "GEN_171",
			Status:  StatusBadRequest,
			Message: "merchant not found",
		},
		"UNSUPPORTED_PHONE_NUMBER_FORMAT": {
			Code:    "GEN_172",
			Status:  StatusBadRequest,
			Message: "Only Ethiopian numbers in local or international format are acceptable",
		},
		"NO_DOC_FOUND": {
			Code:    "GEN_171	",
			Message: "mongo: no documents in result",
		},
		"INVALID_TOTAL_CAP_VALUE": {
			Code:    "GEN_172	",
			Message: "invalid total cap value",
		},
		"INDIVIDUAL_SINGLE_CAP_EXCEEDS_TOTAL_CAP": {
			Code:    "GEN_173	",
			Message: "individual single cap exceed total cap",
		},
		"INDIVIDUAL_DAILY_CAP_EXCEEDS_TOTAL_CAP": {
			Code:    "GEN_174	",
			Message: "individual daily cap exceed total cap",
		},
		"CORPORATE_SINGLE_CAP_EXCEEDS_TOTAL_CAP": {
			Code:    "GEN_175	",
			Message: "corporate single cap exceed total cap",
		},
		"CORPORATE_DAILY_CAP_EXCEEDS_TOTAL_CAP": {
			Code:    "GEN_176	",
			Message: "corporate daily cap exceed total cap",
		},
		"DATABASE_ERROR_CHECKING_PENDING_ACTION": {
			Code:    "GEN_177",
			Message: "An error occured while checking pending action",
		},
		"SINGLE_MAX_TRANSFER_CAN_NOT_LESS_OR_EQUAL": {
			Code:    "GEN_177",
			Message: "Single max transfers can not be less or equal to min_amount",
		},
		"ACCOUNT_NUMBER_CAN_NOT_BE_EMPTY": {
			Code:    "GEN_178",
			Message: "account number can not be empty",
		},
		"USERCODE_CANT_BE_EMPTY": {
			Code:    "GEN_179",
			Message: "user code can not be empty",
		},
		"FAILED_TO_FIND_USER": {
			Code:    "GEN_177",
			Status:  http.StatusNotFound,
			Message: "failed to find user",
		},
		"FAYDA_USER_ALREADY_ENABLED": {
			Code:    "GEN_178",
			Status:  http.StatusConflict,
			Message: "fayda user already enabled",
		},
		"FAYDA_USER_ALREADY_DISABLE": {
			Code:    "GEN_179",
			Status:  http.StatusConflict,
			Message: "fayda user already disable",
		},
		"PLEASE_ADD_VALID_PAGE_OR_PERPAGE": {
			Code:    "GEN_180",
			Status:  http.StatusConflict,
			Message: "please add valid page or per page",
		},
		"FAILED_TO_FETCH_ARCHIVED_USERS": {
			Code:    "GEN_181",
			Status:  http.StatusConflict,
			Message: "Failed to fetch archived users",
		},
		"DUPLICATE_ACTION": {
			Code:    "GEN_182",
			Status:  http.StatusConflict,
			Message: "You are requesting a duplicate action",
		},
		"NO_TOTAL_CAP_FOUND": {
			Code:    "GEN_183",
			Status:  http.StatusNotFound,
			Message: "no total cap data found",
		},
	},
	Auth: ErrorGroup{
		"AUTH_USER_NOT_FOUND": {
			Code:    "AUTH_001",
			Status:  StatusUnauthorized,
			Message: "User not found.",
		},
		"AUTH_USER_DISABLED": {
			Code:    "AUTH_002",
			Status:  StatusUnauthorized,
			Message: "User is not allowed to login, please contact your admin!",
		},
		"AUTH_USER_HAS_NO_PASSWORD": {
			Code:    "AUTH_004",
			Status:  StatusUnauthorized,
			Message: "Please reset your password. to login",
		},
		"AUTH_TOO_MANY_ATTEMPTS": {
			Code:    "AUTH_005",
			Status:  StatusUnauthorized,
			Message: `Too many attempts. Try again after ${waiting_time} minutes`,
		},
		"AUTH_INVALID_OTP": {
			Code:    "AUTH_006",
			Status:  StatusUnauthorized,
			Message: "The verification Code you entered is incorrect. Please check the Code and try again.",
		},
		"AUTH_EXPIRED_OTP": {
			Code:    "AUTH_007",
			Status:  StatusUnauthorized,
			Message: "The verification Code has expired. Please request a new Code to continue.",
		},
		"AUTH_USER_RESET_PASSWORD_REQUIRED": {
			Code:    "AUTH_008",
			Status:  StatusUnauthorized,
			Message: "Too many login attempts. Please reset your password. Your account is locked until then.",
		},
		"AUTH_USER_ALREADY_EXISTS": {
			Code:    "AUTH_009",
			Status:  StatusConflict,
			Message: `user already exist`,
		},
		"LOGIN_PROHIBITED_FOR_15_MIN": {
			Code:    "AUTH_010",
			Status:  StatusUnauthorized,
			Message: "Incorrect credentials. Login is prohibited for 15 minutes.",
		},
		"INCORRECT_PASSWORD": {
			Code:    "AUTH_011",
			Status:  StatusUnauthorized,
			Message: `Incorrect password. ${triesLeft} login attempt(s) left.`,
		},
		"OLD_PASSWORD_SAME_AS_NEW": {
			Code:    "AUTH_012",
			Status:  StatusUnauthorized,
			Message: "New password can't be the same as old password.",
		},
		"PIN_OLY_DIG": {
			Code:    "AUTH_013",
			Status:  StatusUnauthorized,
			Message: "PIN must contain only digits",
		},
		"PIN_LIMIT": {
			Code:    "AUTH_014",
			Status:  StatusUnauthorized,
			Message: "PIN must be exactly 6 digits",
		},
		"PIN_REDANDANT": {
			Code:    "AUTH_015",
			Status:  StatusUnauthorized,
			Message: "PIN cannot contain more than 4 redundant numbers",
		},
		"PIN_SEQ": {
			Code:    "AUTH_016",
			Status:  StatusUnauthorized,
			Message: "PIN cannot contain 4 or more sequential numbers",
		},
		"FAILD_TO_GEN_TOKEN": {
			Code:    "AUTH_017",
			Status:  StatusUnauthorized,
			Message: "Set your new password",
		},
		"FAILD_TO_RESET_PASS": {
			Code:    "AUTH_018",
			Status:  StatusUnauthorized,
			Message: "Failed to set your new password",
		},
		"FAILD_VALIDATION": {
			Code:    "AUTH_019",
			Status:  StatusUnauthorized,
			Message: "Validation failed",
		},
		"OLD_PIN_MISMATCH": {
			Code:    "AUTH_020",
			Status:  StatusUnauthorized,
			Message: "Old PIN does not match",
		},
		"PIN_IN_HISTORY": {
			Code:    "AUTH_021",
			Status:  StatusUnauthorized,
			Message: "New PIN cannot be one of the last 6 PINs used",
		},
		"INVALID_BEARER": {
			Code:    "AUTH_022",
			Status:  StatusUnauthorized,
			Message: "Unauthorized: Invalid Bearer token format",
		},
		"USE_RIGHT_AUTH": {
			Code:    "AUTH_023",
			Status:  StatusUnauthorized,
			Message: "Use the Right Authentication",
		},
		"INVALID_CLAIM": {
			Code:    "AUTH_024",
			Status:  StatusUnauthorized,
			Message: "Invalid token claims",
		},
		"INVALID_TOKEN_DATA": {
			Code:    "AUTH_025",
			Status:  StatusUnauthorized,
			Message: "Invalid token data",
		},
		"UNABLE_TO_DYCRYPT_TOKEN": {
			Code:    "AUTH_026",
			Status:  StatusUnauthorized,
			Message: "Unable to decrypt token",
		},
		"FAILED_LOGIN": {
			Code:    "AUTH_027",
			Status:  StatusUnauthorized,
			Message: "Failed to login. you entered wrong password",
		},
		"ACTION_NOT_FOUND_AFTER_UPDATE": {
			Code:    "ACT_404_UPD",
			Status:  StatusNotFound,
			Message: "action not found after update",
		},
		"ACTION_ALREADY_APPROVED": {
			Code:    "GEN_005",
			Status:  StatusConflict,
			Message: "This action is already approved.",
		},

		"ACCOUNT_LOCKED": {
			Code:    "AUTH_028",
			Status:  StatusUnauthorized,
			Message: "account locked. please contact your addministrator",
		},
		"ACCESS_TOKEN_REQUIRED": {
			Code:    "AUTH_030",
			Status:  StatusUnauthorized,
			Message: "access token required",
		},
		"DEPARTMEN_REQUIRED": {
			Code:    "AUTH_031",
			Status:  StatusUnauthorized,
			Message: "department is required in context",
		},
		"INCOMPLETE_USER_INFO": {
			Code:    "AUTH_032",
			Status:  StatusUnauthorized,
			Message: "Incomplete user information",
		},
		"APP_NAME_EXIST": {
			Code:    "AUTH_033",
			Status:  StatusUnauthorized,
			Message: "The App name already exists",
		},
	},
	Transaction: ErrorGroup{
		"TRANSACTION_NOT_FOUND": {
			Code:    "TXN_001",
			Status:  StatusBadRequest,
			Message: "Transaction not found.",
		},
		"INSUFFICIENT_FUNDS": {
			Code:    "TXN_002",
			Status:  StatusBadRequest,
			Message: "Insufficient balance.",
		},
		"COMMISSION_NOT_FOUND": {
			Code:    "TXN_003",
			Status:  StatusBadRequest,
			Message: "Commission not found",
		},
		"CUSTOMER_NOT_FOUND": {
			Code:    "TXN_004",
			Status:  StatusBadRequest,
			Message: "Customer not found",
		},
	},
	Account: ErrorGroup{
		"PHONE_LOOKUP_FAILED": {
			Code:    "ACC_001",
			Status:  StatusBadRequest,
			Message: "Failed to lookup phone number.",
		},
		"API_REQUEST_FAILED": {
			Code:    "ACC_002",
			Status:  StatusBadRequest,
			Message: "API request failed.",
		},
		"ACCOUNT_NOT_FOUND": {
			Code:    "ACC_003",
			Status:  StatusBadRequest,
			Message: "Account not found.",
		},
		"USER_PHONE_EXISTS": {
			Code:    "ACC_004",
			Status:  StatusConflict,
			Message: "User phone number already exists.",
		},
		"GENERAL_SIF_GENERATION_FAILED": {
			Code:    "ACC_005",
			Status:  StatusBadRequest,
			Message: "Failed to generate SIF and account number.",
		},
		"GENERAL_DB_UPDATE_FAILED": {
			Code:    "ACC_006",
			Status:  StatusBadRequest,
			Message: "Failed to update database record.",
		},
		"GENERAL_DB_INSERT_FAILED": {
			Code:    "ACC_007",
			Status:  StatusBadRequest,
			Message: "Failed to insert database record.",
		},
		"CORE_ACCOUNT_NOT_FOUND": {
			Code:    "ACC_008",
			Status:  StatusBadRequest,
			Message: "Core account not found.",
		},
		"CORE_ACCOUNT_TYPE_ERROR": {
			Code:    "ACC_009",
			Status:  StatusBadRequest,
			Message: "Invalid core account detail type.",
		},
		"ACCOUNT_LOOKUP_FAILED": {
			Code:    "ACC_010",
			Status:  StatusBadRequest,
			Message: "Account lookup failed.",
		},
		"ACCOUNT_DETAIL_FAILED": {
			Code:    "ACC_011",
			Status:  StatusBadRequest,
			Message: "Account detail lookup failed.",
		},
		"INVALID_CORE_RESPONSE": {
			Code:    "ACC_012",
			Status:  StatusBadRequest,
			Message: "Invalid core response.",
		},
		"ACCOUNT_ALREADY_LINKED": {
			Code:    "ACC_013",
			Status:  StatusConflict,
			Message: "Account is already linked to another user.",
		},
		"BLOCKED_ACTION_USER_ALREADY_EXIST": {
			Code:    "AUTH_029",
			Status:  StatusConflict,
			Message: "A pending block action already exists for this user",
		},
		"BRANCH_DISABLE_ACTION_ALREADY_EXISTS": {
			Code:    "BRN_008",
			Status:  StatusBadRequest,
			Message: "A pending disable action already exists for this branch",
		},
		"BRANCH_DISABLE_MULTI_ACTION_ALREADY_EXISTS": {
			Code:    "BRN_010",
			Status:  StatusBadRequest,
			Message: "A pending disable action already exists for one or more of the selected branches",
		},
		"BRANCH_BULK_ACTION_ALREADY_PROCESSED": {
			Code:    "BRN_011",
			Status:  StatusBadRequest,
			Message: "This bulk action has already been processed for one or more branches",
		},
		"BRANCH_ID_REQUIRED": {
			Code:    "BRN_006",
			Status:  StatusBadRequest,
			Message: "Branch ID is required.",
		},
		"BRANCH_REGION_AND_DISTRICT_REQUIRED": {
			Code:    "BRN_012",
			Status:  StatusBadRequest,
			Message: "Region and district are required",
		},
		"BRANCH_REGION_AND_DISTRICT_MIN_LENGTH": {
			Code:    "BRN_014",
			Status:  StatusBadRequest,
			Message: "Region and district must be at least 3 characters",
		},
		"FAILED_TO_FETCH_BRANCHES": {
			Code:    "BRN_013",
			Status:  StatusBadRequest,
			Message: "Failed to fetch branches",
		},
		"INVALID_JSON_PAYLOAD": {
			Code:    "GEN_001",
			Status:  StatusConflict,
			Message: "Invalid JSON body",
		},
		"BRANCH_NOT_FOUND": {
			Code:    "BRN_001",
			Status:  StatusBadRequest,
			Message: "No branch found",
		},
		"UNHANDLED_SERVER_ERROR": {
			Code:    "GEN_004",
			Status:  StatusInternalServerError,
			Message: "Internal server error",
		},
		"BLOCK_DISTRICT_ALREADY_PROCESSED": {
			Code:    "DST_001",
			Status:  StatusConflict,
			Message: "This district action has already been processed",
		},
		// Additional missing error from the other definition
		"ACCOUNT_UPDATE_FAILED": {
			Code:    "ACC_014",
			Status:  StatusBadRequest,
			Message: "Failed to update account.",
		},
		"ACCOUNT_BLOCKED": {
			Code:    "ACC_014",
			Status:  StatusBadRequest,
			Message: "The account is blocked.",
		},
		"REGION_CODE_AND_NAME_REQUIRED": {
			Code:    "ACC_015",
			Status:  StatusBadRequest,
			Message: "Region code  is required",
		},
		"BLOCK_REGION_ALREADY_PROCESSED": {
			Code:    "DST_016",
			Status:  StatusConflict,
			Message: "This region already blocked",
		},
		"BANK_ID_REQUIRED": {
			Code:    "DST_017",
			Status:  StatusBadRequest,
			Message: "Bank id is required",
		},
		"FAILED_TO_GET_DISTRICT": {
			Code:    "DST_018",
			Status:  StatusNotFound,
			Message: "district not found",
		},
		"BRANCH_ALREADY_DISABLE": {
			Code:    "DST_019",
			Status:  StatusConflict,
			Message: "branch code %v already disabled",
		},
	},
	OTP: ErrorGroup{
		"INVALID_OTP": {
			Code:    "OTP_001",
			Status:  StatusBadRequest,
			Message: "Invalid OTP provided.",
		},
		"EXPIRED_OTP": {
			Code:    "OTP_002",
			Status:  StatusBadRequest,
			Message: "OTP has expired.",
		},
		"EMAIL_IN_USE": {
			Code:    "OTP_003",
			Status:  StatusConflict,
			Message: "Email is already in use.",
		},
		"USER_ALREADY_HAS_EMAIL": {
			Code:    "OTP_004",
			Status:  StatusConflict,
			Message: "User already has an email address",
		},
		"WAIT_FOR_PREVIOUS_OTP_EXPIRATION": {
			Code:    "OTP_005",
			Status:  StatusBadRequest,
			Message: "Please wait until the previous OTP expires",
		},
		"OTP_GENERATION_FAILED": {
			Code:    "OTP_006",
			Status:  StatusBadRequest,
			Message: "Failed to generate OTP.",
		},
		"USER_KYC_LEVEL_ZERO": {
			Code:    "OTP_007",
			Status:  StatusBadRequest,
			Message: "User KYC level is not sufficient for this operation",
		},
		"USER_KYC_LEVEL_WRONG": {
			Code:    "OTP_008",
			Status:  StatusBadRequest,
			Message: "User KYC level is not sufficient for this operation",
		},
	},
	File: ErrorGroup{
		"INVALID_FORM": {
			Code:    "FILE_001",
			Status:  StatusBadRequest,
			Message: "Could not parse form or file too large.",
		},
		"NO_FILE": {
			Code:    "FILE_002",
			Status:  StatusBadRequest,
			Message: "File is required.",
		},
		"FILE_TOO_LARGE": {
			Code:    "FILE_003",
			Status:  StatusBadRequest,
			Message: "File exceeds size limit.",
		},
		"INVALID_FILE_TYPE": {
			Code:    "FILE_004",
			Status:  StatusBadRequest,
			Message: "Invalid file type.",
		},
		"UPLOAD_FAILED": {
			Code:    "FILE_005",
			Status:  StatusBadRequest,
			Message: "Failed to upload file.",
		},
		"MISSING_OR_INVALID_IMAGE": {
			Code:    "FILE_006",
			Status:  StatusBadRequest,
			Message: "Missing or invalid image. Could not parse form or file too large.",
		},
	},
	Branch: ErrorGroup{
		"BRANCH_NOT_FOUND": {
			Code:    "BRN_001",
			Status:  StatusNotFound,
			Message: "Branch not found.",
		},
		"BRANCH_DISABLED": {
			Code:    "BRN_002",
			Status:  StatusBadRequest,
			Message: "Branch is disabled.",
		},
		"FAILED_TO_FETCH_BRANCHES": {
			Code:    "BRN_003",
			Status:  StatusBadRequest,
			Message: "Failed to fetch branches.",
		},
		"INVALID_BRANCH_ID": {
			Code:    "BRN_004",
			Status:  StatusBadRequest,
			Message: "Invalid branch ID provided.",
		},
		"INVALID_LOCATION_FILTER": {
			Code:    "BRN_005",
			Status:  StatusBadRequest,
			Message: "Invalid location filter parameters.",
		},
		"BRANCH_ID_REQUIRED": {
			Code:    "BRN_006",
			Status:  StatusBadRequest,
			Message: "Branch ID is required.",
		},
		// Additional missing error from the other definition
		"BRANCHS_ARE_REQUIRED": {
			Code:    "BRN_007",
			Status:  StatusBadRequest,
			Message: "Branches are required.",
		},
		"CITY_NOT_FOUND": {
			Code:    "BRN_008",
			Status:  StatusNotFound,
			Message: "City not found",
		},
		"DISTRICT_ALREADY_BLOCKED": {
			Code:    "BRN_009",
			Status:  StatusConflict,
			Message: "The district already blocked",
		},
		"DISTRICT_NOT_FOUND": {
			Code:    "BRN_010",
			Status:  StatusNotFound,
			Message: "District Not Found",
		},
		"REGION_NOT_FOUND": {
			Code:    "BRN_011",
			Status:  StatusNotFound,
			Message: "Region Not Found",
		},
		"BRANCH_ALREADY_BLOCKED": {
			Code:    "BRN_010",
			Status:  StatusConflict,
			Message: "The branch already blocked",
		},
		"BRANCH_ENABLE_ACTION_ALREADY_EXISTS": {
			Code:    "BRN_013",
			Message: "branch enable action already exists",
		},
		"BRANCH_ALREADY_ENABLED": {
			Code:    "BRN_014",
			Message: "Branch already enabled",
		},
		"MULTIPLE_BRANCH_ENABLE_HAVE_ALREADY_ENABLED_BRANCH": {
			Code:    "BRN_015",
			Message: "multiple branch enable list have  already enabled branch.",
		},
		"BRANCH_CODE_IS_REQUIRED": {
			Code:    "BRN_016",
			Status:  StatusBadRequest,
			Message: "One or more branch code is required",
		},
		"USER_ALREADY_ENABLED": {
			Code:    "BRN_017",
			Status:  StatusConflict,
			Message: "This user is already enabled",
		},
		"USER_ALREADY_DISABLED": {
			Code:    "BRN_018",
			Status:  StatusConflict,
			Message: "This user is already disabled",
		},
		"FAILED_TO_PARSE_FILTERS": {
			Code:    "BRN_019",
			Status:  StatusInternalServerError,
			Message: "Failed to parse filter params",
		},
		"BRANCH_CODE_REQUIRED": {
			Code:    "020",
			Status:  StatusBadRequest,
			Message: "Branch code is required",
		},
	},
	Region: ErrorGroup{
		"REGION_CODE_IS_REQUIRED": {
			Code:    "REG_001",
			Status:  StatusBadRequest,
			Message: "One or more region code is required",
		},
	},
	District: ErrorGroup{
		"DISTRICT_CODE_IS_REQUIRED": {
			Code:    "DIST_001",
			Status:  StatusBadRequest,
			Message: "One or more district code is required",
		},
	},
	City: ErrorGroup{
		"CITY_CODE_IS_REQUIRED": {
			Code:    "DIST_001",
			Status:  StatusBadRequest,
			Message: "One or more city code is required",
		},
	},
	Department: ErrorGroup{
		"DEPARTMENT_NOT_FOUND": {
			Code:    "DEP_001",
			Status:  StatusBadRequest,
			Message: "Department does not exist.",
		},
		"DEPARTMENT_ALREADY_EXISTS": {
			Code:    "DEP_002",
			Status:  StatusConflict,
			Message: "Department already exists.",
		},
		"DEPARTMENT_CODE_REQUIRED": {
			Code:    "DEP_003",
			Status:  StatusBadRequest,
			Message: "Missing department_code for update.",
		},
		"DEPARTMENT_NAME_REQUIRED": {
			Code:    "DEP_004",
			Status:  StatusBadRequest,
			Message: "Department name is required.",
		},
		"PORTAL_CARDS_INVALID": {
			Code:    "DEP_005",
			Status:  StatusBadRequest,
			Message: "Portal cards must be a list of strings.",
		},

		"DEPARTMENT_ALREADY_ENABLED": {
			Code:    "DEP_006",
			Status:  StatusConflict,
			Message: "Department is already enabled.",
		},
		"DEPARTMENT_ALREADY_DISABLED": {
			Code:    "DEP_007",
			Status:  StatusConflict,
			Message: "Department is alreaPORTALdy disabled.",
		},
		"PORTAL_CARD_ARRAY_EMPTY": {
			Code:    "DEP_008",
			Status:  StatusBadRequest,
			Message: "Portal card array cannot be empty.",
		},
		"PORTAL_CARD_NOT_FOUND": {
			Code:    "DEP_009",
			Status:  StatusBadRequest,
			Message: "invalid portal card used.",
		},
	},
	Bank: ErrorGroup{
		"BANKS_NOT_FOUND": {
			Code:    "BNK_001",
			Status:  StatusBadRequest,
			Message: "Bank data not found.",
		},
		"INVALID_BANK_NAME": {
			Code:    "BNK_002",
			Status:  StatusBadRequest,
			Message: "Bank name must be between 3 and 10 alphabetic characters.",
		},
		"MISSING_BANK_NAME": {
			Code:    "BNK_003",
			Status:  StatusBadRequest,
			Message: "Bank name is required.",
		},
		"MISSING_BANK_CODE": {
			Code:    "BNK_004",
			Status:  StatusBadRequest,
			Message: "Bank code is required.",
		},
		"MISSING_BANK_BIC": {
			Code:    "BNK_005",
			Status:  StatusBadRequest,
			Message: "Bank identifier code (BIC) is required.",
		},
		"BANK_ALREADY_EXISTS": {
			Code:    "BNK_006",
			Status:  StatusConflict,
			Message: "Bank with this name already exists.",
		},
		"BANK_ALREADY_ENABLE": {
			Code:    "BNK_007",
			Status:  StatusConflict,
			Message: "Bank  already  enabled",
		},
		"BANK_ALREADY_DISABLED": {
			Code:    "BNK_008",
			Status:  StatusConflict,
			Message: "Bank  already  disabled",
		},
		"NAME_OF_BANK_ALREADY_EXIST": {
			Code:    "BNK_009",
			Status:  StatusConflict,
			Message: "Bank with this name already exist",
		},
		"BIC_CODE_ALREADY_EXIST": {
			Code:    "BNK_010",
			Status:  StatusConflict,
			Message: "Bank BIC CODE already exist",
		},
		"BANK_BIC_CODE_ALREADY_EXIST": {
			Code:    "BNK_011",
			Status:  StatusConflict,
			Message: "Bank BIC CODE already exist",
		},
		"BANK_ALREADY_CREATED_WITH_THIS_PARAMETER": {
			Code:    "BNK_012",
			Status:  StatusConflict,
			Message: "Bank already exist with this parameter",
		},
		"BANK_FETCH_FAILED": {
			Code:    "BNK_013",
			Status:  StatusBadRequest,
			Message: "Bank fetch failed",
		},
		"BANK_NAME_ALREADY_EXIST": {
			Code:    "BNK_014",
			Status:  StatusConflict,
			Message: "Bank Name already exist",
		},
	},
	Action: ErrorGroup{
		"PENDING_CPS_ACTION_EXISTS": {
			Code:    "ACT_001",
			Status:  StatusBadRequest,
			Message: "A pending CPS action already exists.",
		},
		"INVALID_DECISION": {
			Code:    "ACT_002",
			Status:  StatusBadRequest,
			Message: "Invalid decision value provided.",
		},
		"ACTION_REJECTION_FAILED": {
			Code:    "ACT_003",
			Status:  StatusBadRequest,
			Message: "Failed to reject action.",
		},
		"ACTION_APPROVAL_FAILED": {
			Code:    "ACT_004",
			Status:  StatusBadRequest,
			Message: "Failed to approve action.",
		},
		"UNLINK_ACTION_REQUEST_FAILED": {
			Code:    "ACT_005",
			Status:  StatusBadRequest,
			Message: "Failed to request unlink action.",
		},
		"PENDING_ACTION_REJECTION_FAILED": {
			Code:    "ACT_006",
			Status:  StatusBadRequest,
			Message: "Failed to reject pending action.",
		},
		"PENDING_ACTION_CHECK_FAILED": {
			Code:    "ACT_007",
			Status:  StatusBadRequest,
			Message: "Failed to check pending actions.",
		},
		"FAILED_TO_CREATE_ACTION": {
			Code:    "ACT_008",
			Status:  StatusBadRequest,
			Message: "Unable to create action.",
		},
		"ACTION_NOT_PENDING": {
			Code:    "ACT_009",
			Status:  StatusBadRequest,
			Message: "Action is not in pending status.",
		},
		"FAILED_TO_UPDATE_ACTION": {
			Code:    "ACT_010",
			Status:  StatusBadRequest,
			Message: "Failed to update action.",
		}, "FAILED_TO_FETCH_ACTION": {
			Code:    "ACT_011",
			Status:  StatusBadRequest,
			Message: "Failed to fetch action.",
		},
		"FAILED_TO_UPDATE_SERVICE": {
			Code:    "ACT_012",
			Status:  StatusBadRequest,
			Message: "Failed to update service details.",
		},
		"FAILED_TO_UPDATE_CAP_MIN": {
			Code:    "ACT_013",
			Status:  StatusBadRequest,
			Message: "Failed to update cap minimum amount.",
		},
		"MISSING_REJECT_REASON": {
			Code:    "ACT_014",
			Status:  StatusBadRequest,
			Message: "Missing reason for rejection.",
		},
		"REJECT_REASON_TOO_SHORT": {
			Code:    "ACT_015",
			Status:  StatusBadRequest,
			Message: "Rejection reason must be at between 30 to 100 characters long.",
		},
		"PENDING_ACTION_ALREADY_EXIST": {
			Code:    "ACT_016",
			Status:  StatusBadRequest,
			Message: "A pending block action already exists for this city",
		},
		"CITY_ALREADY_BLOCKED": {
			Code:    "ACT_017",
			Status:  StatusConflict,
			Message: "This city already blocked",
		},
		"PENDING_DISTRICT_ACTION_ALREADY_EXIST": {
			Code:    "ACT_017",
			Status:  StatusBadRequest,
			Message: "A pending block action already exists for this district",
		},
		"REGION_ALREADY_ENABLED": {
			Code:    "ACT_018",
			Message: "A region already enbled",
		},
		"ENABLE_REGION_ACTION_ALREADY_EXISTS": {
			Code:    "ACT_019",
			Message: "enabel region action already exists",
		},
		"ACTION_CODE_REQUIRED": {
			Code:    "ACT_020",
			Message: "action code required",
		},
		"USER_IS_NOT_FAYDA_USER": {
			Code:    "ACT_021",
			Message: "user in not fayda user",
		},
		"FAYDA_USER_ALREADY_DISABLED": {
			Code:    "ACT_022",
			Message: "fayida user already disabled",
		},
	},
	User: ErrorGroup{
		"USER_STATUS_UPDATE_FAILED": {
			Code:    "USR_001",
			Status:  StatusBadRequest,
			Message: "Failed to update user status.",
		},
		"FAILED_TO_MARSHAL_INCOMING_USER": {
			Code:    "USR_002",
			Status:  StatusInternalServerError,
			Message: "Failed to marshal incoming user data",
		},
		"FAILED_TO_UNMARSHAL_INCOMING_USER": {
			Code:    "USR_003",
			Status:  StatusInternalServerError,
			Message: "Failed to unmarshal incomming user data",
		},
		"FAILED_TO_MARSHAL_EXISTING_USER": {
			Code:    "USR_004",
			Status:  StatusInternalServerError,
			Message: "Failed to marshal existing user data",
		},
		"FAILED_TO_UNMARSHAL_EXISTING_USER": {
			Code:    "USR_005",
			Status:  StatusInternalServerError,
			Message: "Failed to unmarshal existing user data",
		},
	},
	Wallet: ErrorGroup{
		"WALLET_NOT_FOUND": {
			Code:    "WAL_001",
			Status:  StatusBadRequest,
			Message: "Wallet not found.",
		},
		"WALLET_CREATION_FAILED": {
			Code:    "WAL_002",
			Status:  StatusBadRequest,
			Message: "Failed to create wallet.",
		},
		"WALLET_UPDATE_FAILED": {
			Code:    "WAL_003",
			Status:  StatusBadRequest,
			Message: "Failed to update wallet.",
		},
		"WALLET_DELETION_FAILED": {
			Code:    "WAL_004",
			Status:  StatusBadRequest,
			Message: "Failed to delete wallet.",
		},
		"WALLET_BALANCE_INSUFFICIENT": {
			Code:    "WAL_005",
			Status:  StatusBadRequest,
			Message: "Insufficient wallet balance.",
		},
		"WALLET_TRANSACTION_FAILED": {
			Code:    "WAL_006",
			Status:  StatusBadRequest,
			Message: "Failed to process wallet transaction.",
		},
		"WALLET_NOT_ACTIVE": {
			Code:    "WAL_007",
			Status:  StatusBadRequest,
			Message: "Wallet is not active.",
		},
		"WALLET_ALREADY_EXISTS": {
			Code:    "WAL_008",
			Status:  StatusConflict,
			Message: "Wallet already exists.",
		},
		"WALLET_TYPE_NOT_SUPPORTED": {
			Code:    "WAL_009",
			Status:  StatusBadRequest,
			Message: "Wallet type is not supported.",
		},
		"WALLET_LIMIT_EXCEEDED": {
			Code:    "WAL_010",
			Status:  StatusBadRequest,
			Message: "Wallet limit exceeded.",
		},
		"WALLET_TRANSACTION_NOT_FOUND": {
			Code:    "WAL_011",
			Status:  StatusBadRequest,
			Message: "Wallet transaction not found.",
		},
		"WALLET_TRANSACTION_ALREADY_EXISTS": {
			Code:    "WAL_012",
			Status:  StatusConflict,
			Message: "Wallet transaction already exists.",
		},
	},
	AD: ErrorGroup{
		"AD_NOT_FOUND": {
			Code:    "AD_001",
			Status:  StatusBadRequest,
			Message: "Advert not found.",
		},
		"AD_CREATION_FAILED": {
			Code:    "AD_002",
			Status:  StatusBadRequest,
			Message: "Failed to create advert.",
		},
		"AD_UPDATE_FAILED": {
			Code:    "AD_003",
			Status:  StatusBadRequest,
			Message: "Failed to update advert.",
		},
		"AD_ALREADY_EXISTS": {
			Code:    "AD_004",
			Status:  StatusConflict,
			Message: "Advert with this ID already exists.",
		},
		"AD_DELETION_FAILED": {
			Code:    "AD_005",
			Status:  StatusBadRequest,
			Message: "Failed to delete advert.",
		},
		"START_DATE_REQUIRED": {
			Code:    "AD_006",
			Status:  StatusBadRequest,
			Message: "started date is required",
		},
		"EXPIRE_DATE_REQUIRED": {
			Code:    "AD_006",
			Status:  StatusBadRequest,
			Message: "end date is required",
		},
	},
	BulkService: ErrorGroup{
		"FAILED_TO_UPDATE_PARENT": {
			Code:    "BULK_001",
			Status:  StatusInternalServerError,
			Message: "Failed to update parent access list",
		},
		"FAILED_TO_UPDATE_CHILD": {
			Code:    "BULK_002",
			Status:  StatusInternalServerError,
			Message: "Failed to update sub access list",
		},
		"INVALID_KEY_FORMAT": {
			Code:    "BULK_003",
			Status:  StatusBadRequest,
			Message: "The keys you entered are not valid",
		},
		"INVALID_CURRENT_ACTION": {
			Code:    "BULK_004",
			Status:  StatusBadRequest,
			Message: "Current action format is not valid",
		},
		"SURVICE_NOT_FOUND": {
			Code:    "BULK_005",
			Status:  StatusNotFound,
			Message: "The bulk service you requested is not found",
		},
		"NO_RESOURCE_FOUND": {
			Code:    "BULK_006",
			Status:  StatusNotFound,
			Message: "No document/resource found",
		},
		"BULK_SERVICE_CODE_IS_REQUIRED": {
			Code:    "BULK_007",
			Status:  StatusBadRequest,
			Message: "One or more bulk service code is required",
		},
	},
	Permission: ErrorGroup{
		"NO_PERMISSION_CATEGORY_FOUND": {
			Code:    "PERM_001",
			Status:  StatusBadRequest,
			Message: "No permission category found.",
		},
		"INVALID_PERMISSION_CATEGORY_ID": {
			Code:    "PERM_002",
			Status:  StatusBadRequest,
			Message: "Invalid permission category ID provided.",
		},
		"ONE_OR_MORE_PERMISSION_CATEGORIES_NOT_FOUND": {
			Code:    "PERM_003",
			Status:  StatusBadRequest,
			Message: "One or more permission categories not found.",
		},
		"INVALID_PERMISSION_GROUP_ID": {
			Code:    "PERM_004",
			Status:  StatusBadRequest,
			Message: "Invalid permission group ID provided.",
		},
		"NO_PERMISSION_GROUP_FOUND": {
			Code:    "PERM_005",
			Status:  StatusBadRequest,
			Message: "No permission group found.",
		},
		"ONE_OR_MORE_PERMISSION_GROUPS_NOT_FOUND": {
			Code:    "PERM_006",
			Status:  StatusBadRequest,
			Message: "One or more permission groups not found.",
		},
	},
	MiniApp: ErrorGroup{
		"APP_NAME_REQUIRED": {
			Code:    "MINIAPP_001",
			Status:  StatusBadRequest,
			Message: "app_name is required.",
		},
		"MERCHANT_ID_REQUIRED": {
			Code:    "MINIAPP_002",
			Status:  StatusBadRequest,
			Message: "merchant_id is required.",
		},
		"APP_ICON_REQUIRED": {
			Code:    "MINIAPP_003",
			Status:  StatusBadRequest,
			Message: "app_icon is required.",
		},
		"APP_VIEW_TYPE_REQUIRED": {
			Code:    "MINIAPP_004",
			Status:  StatusBadRequest,
			Message: "app_view_type is required.",
		},
		"APP_TYPE_REQUIRED": {
			Code:    "MINIAPP_005",
			Status:  StatusBadRequest,
			Message: "app_type is required.",
		},
		"URL_REQUIRED": {
			Code:    "MINIAPP_006",
			Status:  StatusBadRequest,
			Message: "url is required.",
		},
		"INVALID_URL": {
			Code:    "MINIAPP_007",
			Status:  StatusBadRequest,
			Message: "invalid url.",
		},
		"MPAAS_ID_REQUIRED": {
			Code:    "MINIAPP_008",
			Status:  StatusBadRequest,
			Message: "mpaas_id is required.",
		},
		"INCOMPLETE_BRANCH_PRODUCT_CODES": {
			Code:    "MINIAPP_009",
			Status:  StatusBadRequest,
			Message: "all fields for branch_type must be provided if one is set.",
		},
		"NO_PRODUCT_CODES_PROVIDED": {
			Code:    "MINIAPP_010",
			Status:  StatusBadRequest,
			Message: "at least one complete set of product codes (IFB or CB) must be provided.",
		},
		"BOTH_PRODUCT_CODES_REQUIRED": {
			Code:    "MINIAPP_011",
			Status:  StatusBadRequest,
			Message: "both IFB and CB product codes must be provided for BOTH app_view_type.",
		},
		"APP_VIEW_TYPE_INVALID_OR_MISSING": {
			Code:    "MINIAPP_012",
			Status:  StatusBadRequest,
			Message: "app_view_type is required and must be one of: BOTH, CB, IFB.",
		},
		"INVALID_APP_VIEW_TYPE": {
			Code:    "MINIAPP_013",
			Status:  StatusBadRequest,
			Message: "invalid app_view_type; must be one of BOTH, CB, IFB.",
		},
		"APP_TYPE_MISSING": {
			Code:    "MINIAPP_014",
			Status:  StatusBadRequest,
			Message: "one app type must be provided.",
		},
		"EXCLUSIVE_APP_FLAGS": {
			Code:    "MINIAPP_015",
			Status:  StatusBadRequest,
			Message: "only one of is_event_mini_app or is_three_click can be true.",
		},
	},

	Event: ErrorGroup{
		"EVENT_NAME_ALREADY_EXISTS": {
			Code:    "EVE_001",
			Message: "Event name already exists",
			Status:  StatusBadRequest,
		},
	},
	Service: ErrorGroup{
		"SERVICE_FEE_VALIDATION_FAILED": {
			Code:    "SRV_001",
			Message: "Service fee validation failed",
			Status:  StatusBadRequest,
		},
		"SERVICE_FEE_UPDATE_FAILED": {
			Code:    "SRV_002",
			Message: "Failed to update service fee",
			Status:  StatusInternalServerError,
		},
		"SERVICE_FEE_NOT_FOUND": {
			Code:    "SRV_003",
			Message: "Service fee not found",
			Status:  StatusNotFound,
		},
		"INVALID_TIER_STRUCTURE": {
			Code:    "SRV_004",
			Message: "Invalid tier structure provided",
			Status:  StatusBadRequest,
		},
		"TIER_CONTINUITY_VIOLATION": {
			Code:    "SRV_005",
			Message: "Tier continuity violation: each tier's max must equal next tier's min",
			Status:  StatusBadRequest,
		},
		"INVALID_TRANSFER_CAP": {
			Code:    "SRV_006",
			Message: "Invalid transfer cap configuration",
			Status:  StatusBadRequest,
		},
		"TRANSFER_CAP_EXCEEDS_TOTAL": {
			Code:    "SRV_007",
			Message: "Transfer cap cannot exceed total cap",
			Status:  StatusBadRequest,
		},
		"GL_ENTRY_VALIDATION_FAILED": {
			Code:    "SRV_008",
			Message: "GL entry validation failed",
			Status:  StatusBadRequest,
		},
		"PENDING_SERVICE_FEE_ACTION": {
			Code:    "SRV_009",
			Message: "A pending service fee action already exists",
			Status:  StatusConflict,
		},
		"TOTAL_CAP_VALIDATION_FAILED":{
			Code:    "SRV_010",
			Message: "Total cap validation failed",
			Status:  StatusBadRequest,
		},
	},
	Donation: ErrorGroup{
		"CATEGORY_NAME_ALREADY_EXISTS": {
			Code:    "DON_001",
			Message: "Category name already exists",
			Status:  StatusBadRequest,
		},
		"FAILED_TO_UPLOAD_ICON": {
			Code:    "DON_002",
			Message: "Failed to upload icon",
			Status:  StatusBadRequest,
		},
		"ICON_IS_REQUIRED": {
			Code:    "DON_003",
			Message: "Icon is required",
			Status:  StatusBadRequest,
		},
		"MINIO_TIME_SYNC_ERROR": {
			Code:    "DON_004",
			Message: "MinIO time synchronization error",
			Status:  StatusInternalServerError,
		},
		"MINIO_BUCKET_ERROR": {
			Code:    "DON_005",
			Message: "MinIO bucket operation failed",
			Status:  StatusInternalServerError,
		},
		"MINIO_CLIENT_NOT_CONFIGURED": {
			Code:    "DON_006",
			Message: "MinIO client is not configured",
			Status:  StatusInternalServerError,
		},
		"CONFIGURATION_NOT_LOADED": {
			Code:    "DON_007",
			Message: "Configuration is not loaded",
			Status:  StatusInternalServerError,
		},
		"UNABLE_TO_CHECK_ACCOUNT": {
			Code:    "DON_008",
			Message: "Unable to check account",
			Status:  StatusInternalServerError,
		},
		"TIME_OUT_ERROR": {
			Code:    "DON_009",
			Message: "Operation timed out",
			Status:  StatusInternalServerError,
		},
		"INVALID_DONATION_AMOUNT": {
			Code:    "DON_010",
			Message: "Invalid donation amount",
			Status:  StatusBadRequest,
		},
		"INVALID_END_DATE_FORMAT": {
			Code:    "DON_011",
			Message: "Invalid end date format",
			Status:  StatusBadRequest,
		},
		"INVALID_START_DATE_FORMAT": {
			Code:    "DON_012",
			Message: "Invalid start date format",
			Status:  StatusBadRequest,
		},
		"DONATION_IMAGES_REQUIRED": {
			Code:    "DON_013",
			Message: "Donation images are required",
			Status:  StatusBadRequest,
		},
		"COMPANY_NAME_ALREADY_EXISTS": {
			Code:    "DON_014",
			Message: "Company name already exists",
			Status:  StatusBadRequest,
		},
		"ACCOUNT_NUMBER_ALREADY_EXISTS": {
			Code:    "DON_015",
			Message: "Account number already exists",
			Status:  StatusBadRequest,
		},
		"LOGO_IS_REQUIRED": {
			Code:    "DON_016",
			Message: "Logo is required",
			Status:  StatusBadRequest,
		},
		"FAILED_TO_UPLOAD_LOGO": {
			Code:    "DON_017",
			Message: "Failed to upload logo",
			Status:  StatusBadRequest,
		},
		"DONATION_TITLE_ALREADY_EXISTS": {
			Code:    "DON_018",
			Message: "Donation title already exists",
			Status:  StatusBadRequest,
		},
		"COMPANY_NOT_FOUND": {
			Code:    "DON_019",
			Message: "Company not found",
			Status:  StatusBadRequest,
		},
		"CATEGORY_NOT_FOUND": {
			Code:    "DON_020",
			Message: "Category not found",
			Status:  StatusBadRequest,
		},
		"AT_LEAST_ONE_IMAGE_REQUIRED": {
			Code:    "DON_021",
			Message: "At least one image is required",
			Status:  StatusBadRequest,
		},
		"FAILED_TO_UPLOAD_IMAGES": {
			Code:    "DON_022",
			Message: "Failed to upload images",
			Status:  StatusBadRequest,
		},
		"ICON_URL_MISSING": {
			Code:    "DON_023",
			Message: "Icon URL is missing in CPS request",
			Status:  StatusBadRequest,
		},
		"PREVIOUS_ACTION_REQUIRED": {
			Code:    "DON_024",
			Message: "Previous action is required for update",
			Status:  StatusBadRequest,
		},
		"INVALID_PREVIOUS_ACTION_FORMAT": {
			Code:    "DON_025",
			Message: "Invalid previous action format",
			Status:  StatusBadRequest,
		},
		"ID_NOT_FOUND": {
			Code:    "DON_026",
			Message: "ID not found in previous action",
			Status:  StatusBadRequest,
		},
		"INVALID_ID_FORMAT": {
			Code:    "DON_027",
			Message: "Invalid ID format",
			Status:  StatusBadRequest,
		},
		"LOGO_URL_MISSING": {
			Code:    "DON_028",
			Message: "Logo URL is missing in CPS request",
			Status:  StatusBadRequest,
		},
		"DONATION_IMAGES_MISSING": {
			Code:    "DON_029",
			Message: "Donation images are missing in CPS request",
			Status:  StatusBadRequest,
		},
		"ACCOUNT_NUMBER_IS_REQUIRED": {
			Code:    "DON_030",
			Message: "Account number is required",
			Status:  StatusBadRequest,
		},
		"ACCOUNT_NOT_FOUND": {
			Code:    "DON_031",
			Message: "Account not found",
			Status:  StatusBadRequest,
		},
		"COMPANY_LOOKUP_FAILED": {
			Code:    "DON_032",
			Message: "Company lookup failed",
			Status:  StatusInternalServerError,
		},
		"CATEGORY_LOOKUP_FAILED": {
			Code:    "DON_033",
			Message: "Category lookup failed",
			Status:  StatusInternalServerError,
		},
		"FAILED_TO_UPLOAD_IMAGE_1": {
			Code:    "DON_034",
			Message: "Failed to upload image 1",
			Status:  StatusBadRequest,
		},
		"FAILED_TO_UPLOAD_IMAGE_2": {
			Code:    "DON_035",
			Message: "Failed to upload image 2",
			Status:  StatusBadRequest,
		},
		"FAILED_TO_UPLOAD_IMAGE_3": {
			Code:    "DON_036",
			Message: "Failed to upload image 3",
			Status:  StatusBadRequest,
		},
		"FAILED_TO_UPLOAD_IMAGE_4": {
			Code:    "DON_037",
			Message: "Failed to upload image 4",
			Status:  StatusBadRequest,
		},
		"FAILED_TO_UPLOAD_IMAGE_5": {
			Code:    "DON_038",
			Message: "Failed to upload image 5",
			Status:  StatusBadRequest,
		},
		"FAILED_TO_UPDATE_DONATION": {
			Code:    "DON_039",
			Message: "Failed to update donation",
			Status:  StatusInternalServerError,
		},
		"FAILED_TO_UPDATE_DONATION_CATEGORY": {
			Code:    "DON_040",
			Message: "Failed to update donation category",
			Status:  StatusInternalServerError,
		},
		"FAILED_TO_UPDATE_DONATION_COMPANY": {
			Code:    "DON_041",
			Message: "Failed to update donation company",
			Status:  StatusInternalServerError,
		},
		"FAILED_TO_CREATE_DONATION": {
			Code:    "DON_042",
			Message: "Failed to create donation",
			Status:  StatusInternalServerError,
		},
		"FAILED_TO_CREATE_DONATION_CATEGORY": {
			Code:    "DON_043",
			Message: "Failed to create donation category",
			Status:  StatusInternalServerError,
		},
		"FAILED_TO_CREATE_DONATION_COMPANY": {
			Code:    "DON_044",
			Message: "Failed to create donation company",
			Status:  StatusInternalServerError,
		},
		"FAILED_TO_CHECK_DONATION_COMPANY_NAME": {
			Code:    "DON_045",
			Message: "Failed to check donation company name",
			Status:  StatusInternalServerError,
		},
		"FAILED_TO_CHECK_DONATION_COMPANY_ACCOUNT": {
			Code:    "DON_046",
			Message: "Failed to check donation company account",
			Status:  StatusInternalServerError,
		},
		"FAILED_TO_FETCH_DONATION_COMPANIES": {
			Code:    "DON_047",
			Message: "Failed to fetch donation companies",
			Status:  StatusInternalServerError,
		},
		"FAILED_TO_DECODE_DONATION_COMPANIES": {
			Code:    "DON_048",
			Message: "Failed to decode donation companies",
			Status:  StatusInternalServerError,
		},
		"FAILED_TO_COUNT_DONATION_COMPANIES": {
			Code:    "DON_049",
			Message: "Failed to count donation companies",
			Status:  StatusInternalServerError,
		},
		"FAILED_TO_FETCH_DONATION_CATEGORIES": {
			Code:    "DON_050",
			Message: "Failed to fetch donation categories",
			Status:  StatusInternalServerError,
		},
		"FAILED_TO_DECODE_DONATION_CATEGORIES": {
			Code:    "DON_051",
			Message: "Failed to decode donation categories",
			Status:  StatusInternalServerError,
		},
		"FAILED_TO_COUNT_DONATION_CATEGORIES": {
			Code:    "DON_052",
			Message: "Failed to count donation categories",
			Status:  StatusInternalServerError,
		},
		"FAILED_TO_FETCH_DONATIONS": {
			Code:    "DON_053",
			Message: "Failed to fetch donations",
			Status:  StatusInternalServerError,
		},
		"FAILED_TO_DECODE_DONATIONS": {
			Code:    "DON_054",
			Message: "Failed to decode donations",
			Status:  StatusInternalServerError,
		},
		"FAILED_TO_COUNT_DONATIONS": {
			Code:    "DON_055",
			Message: "Failed to count donations",
			Status:  StatusInternalServerError,
		},
		"FAILED_TO_CONVERT_STRING_TO_OBJECT_ID": {
			Code:    "DON_056",
			Message: "Failed to convert string to Object ID",
			Status:  StatusBadRequest,
		},
		"FAILED_TO_SAFE_CONVERT_TO_OBJECT_ID": {
			Code:    "DON_057",
			Message: "Failed to safely convert to Object ID",
			Status:  StatusBadRequest,
		},
		"DONATION_LOOKUP_FAILED": {
			Code:    "DON_058",
			Message: "Donation lookup failed",
			Status:  StatusInternalServerError,
		},
		"INVALID_COMPANY_ID_FORMAT": {
			Code:    "DON_060",
			Message: "Invalid company ID format",
			Status:  StatusBadRequest,
		},
		"INVALID_CATEGORY_ID_FORMAT": {
			Code:    "DON_061",
			Message: "Invalid category ID format",
			Status:  StatusBadRequest,
		},
	},
}

func (e ErrorDefinition) Error() string {
	return e.Message
}
