package utils

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/common"
)

type APIErrorResponse struct {
	Message       string      `json:"message"`
	Code          string      `json:"code,omitempty"`
	AccountLocked bool        `json:"accountLocked,omitempty"`
	Auth          bool        `json:"auth,omitempty"`
	Errors        interface{} `json:"errors,omitempty"`
}

var errorKeyToStatus = map[string]int{
	// General
	"CONFLICT_KEY":                                   409,
	"INVALID_ID":                                     400,
	"INVALID_JSON_PAYLOAD":                           400,
	"UNHANDLED_SERVER_ERROR":                         500,
	"INVALID_TOKEN_FORMAT":                           400,
	"INVALID_TOKEN":                                  401,
	"EXPIRED_TOKEN":                                  401,
	"ERROR_CHANGING_PIN":                             400,
	"ENCRYPTED_PAYLOAD_REQUIRED":                     400,
	"INVALID_SERVER_CONFIGURATION":                   500,
	"AUTH_HEADER_MISSING":                            401,
	"INVALID_KEY_CONFIGURATION":                      500,
	"INVALID_ENCRYPTED_PAYLOAD":                      400,
	"DECRYPTION_ERROR":                               400,
	"SERVER_KEYS_NOT_CONFIGURED":                     500,
	"INVALID_INPUT":                                  400,
	"UNAUTHORIZED":                                   401,
	"USER_REALM_NOT_FOUND":                           403,
	"ACTION_NOT_ALLOWED":                             403,
	"WAIT_FOR_PREVIOUS_OTP_EXPIRATION":               429,
	"EMPTY_ORG":                                      400,
	"INVALID_REQ":                                    400,
	"REG_FRST":                                       403,
	"WEAK_PIN":                                       400,
	"NOT_FOUND":                                      404,
	"DEVICE_ID_REQUIRED":                             400,
	"NO_LINKED_DEVICES":                              404,
	"COULD_NOT_UNLINK_DEVICE":                        500,
	"INCOMPLETE_USER_INFO":                           400,
	"INVALID_ACTION_TYPE":                            400,
	"ACTION_NOT_FOUND":                               400,
	"PENDING_REQUEST_EXISTS":                         409,
	"MISSING_REQUIRED_HEADERS":                       400,
	"DEVICE_LOOKUP_FAILED":                           500,
	"MISSING_REQUIRED_FIELDS":                        400,
	"MISSING_OTP":                                    400,
	"INVALID_INPUT_PARAMETERS":                       400,
	"SAME_PIN":                                       400,
	"PIN_IN_HISTORY":                                 400,
	"ERROR_SETTING_PIN":                              500,
	"OTP_CREATION_FAILED":                            500,
	"DEVICE_NOT_FOUND":                               404,
	"DEVICE_FOUND":                                   200,
	"INVALID_PIN":                                    400,
	"PIN_RESET_SESSION_EXPIRED":                      400,
	"PIN_RESET_SESSION_NOT_FOUND":                    404,
	"PIN_RESET_SESSION_ALREADY_VERIFIED":             400,
	"TOKEN_GENERATION_FAILED":                        500,
	"MAXIMUM_AMOUNT_REQUIRED_FOR_OPEN_METHOD":        400,
	"MINIMUM_AMOUNT_REQUIRED_FOR_OTP_AND_PIN_METHOD": 400,
	"EITHER_MINIMUM_OR_MAXIMUM_AMOUNT_REQUIRED_FOR_PIN_METHOD":   400,
	"INVALID_OBJECT_ID_FORMAT":                                   400,
	"AUTH_TIER_NOT_FOUND":                                        404,
	"AUTH_TIER_ALREADY_EXISTS":                                   409,
	"AUTH_TIER_VALIDATION_FAILED":                                400,
	"AUTH_TIER_UPDATE_FAILED":                                    500,
	"AUTH_TIER_INSERT_FAILED":                                    500,
	"AUTH_TIER_NOT_FOUND_FOR_ID":                                 404,
	"AUTH_TIER_FETCH_FAILED":                                     500,
	"AUTH_TIER_APPROVE_FAILED":                                   500,
	"AUTH_TIER_REJECT_FAILED":                                    500,
	"AUTH_TIER_NOT_FOUND_FOR_ID_AND_DEPARTMENT":                  404,
	"AUTH_TIER_FETCH_FAILED_FOR_ID_AND_DEPARTMENT":               500,
	"AUTH_TIER_APPROVE_FAILED_FOR_ID_AND_DEPARTMENT":             500,
	"AUTH_TIER_REJECT_FAILED_FOR_ID_AND_DEPARTMENT":              500,
	"AUTH_TIER_UPDATE_FAILED_FOR_ID_AND_DEPARTMENT":              500,
	"AUTH_TIER_INSERT_FAILED_FOR_ID_AND_DEPARTMENT":              500,
	"FAILED_TO_UPDATE_OPEN_AUTH_TIER":                            500,
	"FAILED_TO_UPDATE_PIN_AUTH_TIER":                             500,
	"FAILED_TO_UPDATE_OTP_AND_PIN_AUTH_TIER":                     500,
	"FAILED_TO_UPDATE_OTP_AND_OPEN_AUTH_TIER":                    500,
	"FAILED_TO_UPDATE_OPEN_AUTH_TIER_FOR_ID":                     500,
	"MIN_AMOUNT_CANNOT_BE_GREATER_THAN_PIN_MIN_AMOUNT":           400,
	"PIN_AUTHIER_NOT_FOUND":                                      404,
	"PIN_MIN_AMOUNT_CANNOT_BE_LESS_THAN_OPEN_MIN_AMOUNT":         400,
	"MAX_PIN_AMOUNT_SHOULD_BE_GREATER_THAN_PIN_MIN_AMOUNT":       400,
	"OPEN_AUTHIER_NOT_FOUND":                                     404,
	"OPEN_TIER_MAX_AMOUNT_CANNOT_BE_GREATER_THAN_PIN_MAX_AMOUNT": 400,
	"PENDING_CPS_ACTION_PRESENT":                                 409,
	"FAILED_TO_GET_FAYDA_ACCOUNT":                                500,
	"FAILED_TO_GET_CUSTOMER_ACCOUNT":                             500,
	"FAILED_TO_CREATE_CPS_ACTION":                                500,
	"FAILED_TO_UPDATE_CPS_ACTION":                                500,
	"FAILED_TO_UPDATE_CUSTOMER_ACCOUNT":                          500,
	"FAILED_TO_UPDATE_CPS_ACTION_FOR_ID_AND_DEPARTMENT":          500,
	"FAILED_TO_GET_CPS_ACTION":                                   500,
	"FAILED_TO_GET_SERVICE":                                      500,
	"SERVICE_NOT_FOUND":                                          404,
	"FAILED_TO_FIND_CPS_ACTION":                                  500,
	"FAILED_TO_UPDATE_SERVICE_FEE":                               500,
	"FAILED_TO_REJECT_SERVICE_FEE_UPDATE":                        500,

	// Auth
	"AUTH_USER_NOT_FOUND":               404,
	"AUTH_USER_DISABLED":                403,
	"AUTH_INVALID_PASSWORD":             401,
	"AUTH_USER_HAS_NO_PASSWORD":         400,
	"AUTH_TOO_MANY_ATTEMPTS":            429,
	"AUTH_INVALID_OTP":                  400,
	"AUTH_EXPIRED_OTP":                  400,
	"AUTH_USER_RESET_PASSWORD_REQUIRED": 403,
	"AUTH_USER_ALREADY_EXISTS":          409,
	"LOGIN_PROHIBITED_FOR_15_MIN":       429,
	"INCORRECT_PASSWORD":                401,
	"OLD_PASSWORD_SAME_AS_NEW":          400,
	"PIN_OLY_DIG":                       400,
	"PIN_LIMIT":                         400,
	"PIN_REDANDANT":                     400,
	"PIN_SEQ":                           400,
	"FAILD_TO_GEN_TOKEN":                500,
	"FAILD_TO_RESET_PASS":               500,
	"FAILD_VALIDATION":                  400,
	"OLD_PIN_MISMATCH":                  400,
	"INVALID_BEARER":                    401,
	"USE_RIGHT_AUTH":                    401,
	"INVALID_CLAIM":                     401,
	"INVALID_TOKEN_DATA":                401,
	"UNABLE_TO_DYCRYPT_TOKEN":           401,
	"FAILED_LOGIN":                      401,
	"ACCOUNT_LOCKED":                    403,

	// Transaction
	"TRANSACTION_NOT_FOUND": 404,
	"INSUFFICIENT_FUNDS":    402,
	"COMMISSION_NOT_FOUND":  404,
	"CUSTOMER_NOT_FOUND":    404,

	// Account
	"PHONE_LOOKUP_FAILED":           502,
	"API_REQUEST_FAILED":            502,
	"ACCOUNT_NOT_FOUND":             404,
	"USER_PHONE_EXISTS":             409,
	"GENERAL_SIF_GENERATION_FAILED": 500,
	"GENERAL_DB_UPDATE_FAILED":      500,
	"GENERAL_DB_INSERT_FAILED":      500,
	"CORE_ACCOUNT_NOT_FOUND":        404,
	"CORE_ACCOUNT_TYPE_ERROR":       400,
	"ACCOUNT_LOOKUP_FAILED":         502,
	"ACCOUNT_DETAIL_FAILED":         502,
	"INVALID_CORE_RESPONSE":         502,
	"ACCOUNT_ALREADY_LINKED":        409,

	// OTP
	"INVALID_OTP":            400,
	"EXPIRED_OTP":            400,
	"EMAIL_IN_USE":           409,
	"USER_ALREADY_HAS_EMAIL": 409,
	"OTP_GENERATION_FAILED":  500,
	"USER_KYC_LEVEL_ZERO":    400,
	"USER_KYC_LEVEL_WRONG":   400,

	// File
	"INVALID_FORM":      400,
	"NO_FILE":           400,
	"FILE_TOO_LARGE":    413,
	"INVALID_FILE_TYPE": 400,
	"UPLOAD_FAILED":     500,

	// Branch
	"BRANCH_NOT_FOUND":         404,
	"BRANCH_DISABLED":          400,
	"FAILED_TO_FETCH_BRANCHES": 500,
	"INVALID_BRANCH_ID":        400,
	"INVALID_LOCATION_FILTER":  400,
	"BRANCH_ID_REQUIRED":       400,

	// Department
	"DEPARTMENT_NOT_FOUND":      404,
	"DEPARTMENT_ALREADY_EXISTS": 409,
	"DEPARTMENT_CODE_REQUIRED":  400,
	"DEPARTMENT_NAME_REQUIRED":  400,
	"PORTAL_CARDS_INVALID":      400,

	// unidentified key
	"UNIDENTIFIED_KEY": 500,

	// Banks
	"BANKS_NOT_FOUND": 404,
}

func getStatusForErrorKey(key string) int {
	if status, ok := errorKeyToStatus[key]; ok {
		return status
	}
	// Backward compatibility: check DefineError groups
	for _, group := range []common.ErrorGroup{
		common.DefineError.General,
		common.DefineError.Auth,
		common.DefineError.User,
		common.DefineError.Transaction,
		common.DefineError.Account,
		common.DefineError.OTP,
		common.DefineError.File,
		common.DefineError.Branch,
	} {
		if _, ok := group[key]; ok {
			return http.StatusBadRequest
		}
	}
	return http.StatusInternalServerError
}

func SendErrorResponse(w http.ResponseWriter, errorKey string, statusCode int, additionalData map[string]interface{}) {
	var errorDef common.ErrorDefinition
	found := false

	errorGroups := []common.ErrorGroup{
		common.DefineError.General,
		common.DefineError.Auth,
		common.DefineError.User,
		common.DefineError.Transaction,
		common.DefineError.Account,
		common.DefineError.OTP,
		common.DefineError.File,
		common.DefineError.Branch,
	}

	for _, group := range errorGroups {
		if def, ok := group[errorKey]; ok {
			errorDef = def
			found = true
			break
		}
	}

	if !found {
		errorDef = common.ErrorDefinition{
			Code:    errorKey,
			Message: errorKey,
		}
	}

	response := APIErrorResponse{
		Message: errorDef.Message,
		Code:    errorDef.Code,
	}

	if additionalData != nil {
		if accountLocked, ok := additionalData["accountLocked"].(bool); ok {
			response.AccountLocked = accountLocked
		}
		if auth, ok := additionalData["auth"].(bool); ok {
			response.Auth = auth
		}
		if errors, exists := additionalData["errors"]; exists {
			response.Errors = errors
		}
	}

	status := statusCode
	if status == 0 {
		status = getStatusForErrorKey(errorKey)
	}

	returnData := make(map[string]interface{})
	returnData["status"] = status
	returnData["message"] = response.Message
	returnData["data"] = map[string]interface{}{"code": response.Code}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	responseBytes, _ := json.Marshal(returnData)

	w.Write(responseBytes)
}

func HandleServiceError(w http.ResponseWriter, err error) {
	if def, ok := err.(common.ErrorDefinition); ok {
		SendErrorResponse(w, def.Code, http.StatusBadRequest, nil)
		return
	}
	SendErrorResponse(w, "GEN_004", http.StatusInternalServerError, nil)
}

// Department errors
var (
	ErrDepartmentAlreadyExists = errors.New("DEPARTMENT_ALREADY_EXISTS")
	ErrDepartmentNotFound      = errors.New("DEPARTMENT_NOT_FOUND")
	ErrDepartmentCodeRequired  = errors.New("DEPARTMENT_CODE_REQUIRED")
	ErrDepartmentNameRequired  = errors.New("DEPARTMENT_NAME_REQUIRED")
	ErrPortalCardsInvalid      = errors.New("PORTAL_CARDS_INVALID")
)

// Action errors
var (
	ErrActionNotFound       = errors.New("ACTION_NOT_FOUND")
	ErrActionNotAllowed     = errors.New("ACTION_NOT_ALLOWED")
	ErrPendingRequestExists = errors.New("PENDING_REQUEST_EXISTS")
	ErrInvalidActionType    = errors.New("INVALID_ACTION_TYPE")
)

// Input and decoding errors
var (
	ErrInvalidInput       = errors.New("INVALID_INPUT")
	ErrInvalidJSONPayload = errors.New("INVALID_JSON_PAYLOAD")
)

// General errors
var (
	ErrIncompleteUserInfo = errors.New("INCOMPLETE_USER_INFO")
)

// File
var (
	ErrInvalidForm = errors.New("INVALID_FORM")
)
