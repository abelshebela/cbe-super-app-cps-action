package utils

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/http"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/common"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type APIErrorResponse struct {
	Message       string `json:"message"`
	Code          string `json:"code,omitempty"`
	AccountLocked bool   `json:"accountLocked,omitempty"`
	Auth          bool   `json:"auth,omitempty"`
	Errors        any    `json:"errors,omitempty"`
}

// errorKeyToStatus provides the most appropriate HTTP status code for each error key.
var errorKeyToStatus = map[string]int{
	// General
	"CONFLICT_KEY":                                   http.StatusConflict,
	"INVALID_ID":                                     http.StatusBadRequest,
	"INVALID_JSON_PAYLOAD":                           http.StatusBadRequest,
	"UNHANDLED_SERVER_ERROR":                         http.StatusInternalServerError,
	"INVALID_TOKEN_FORMAT":                           http.StatusBadRequest,
	"INVALID_TOKEN":                                  http.StatusUnauthorized,
	"EXPIRED_TOKEN":                                  http.StatusUnauthorized,
	"ERROR_CHANGING_PIN":                             http.StatusBadRequest,
	"ENCRYPTED_PAYLOAD_REQUIRED":                     http.StatusBadRequest,
	"INVALID_SERVER_CONFIGURATION":                   http.StatusInternalServerError,
	"AUTH_HEADER_MISSING":                            http.StatusUnauthorized,
	"INVALID_KEY_CONFIGURATION":                      http.StatusInternalServerError,
	"INVALID_ENCRYPTED_PAYLOAD":                      http.StatusBadRequest,
	"DECRYPTION_ERROR":                               http.StatusBadRequest,
	"SERVER_KEYS_NOT_CONFIGURED":                     http.StatusInternalServerError,
	"INVALID_INPUT":                                  http.StatusBadRequest,
	"UNAUTHORIZED":                                   http.StatusUnauthorized,
	"USER_REALM_NOT_FOUND":                           http.StatusNotFound,
	"ACTION_NOT_ALLOWED":                             http.StatusForbidden,
	"WAIT_FOR_PREVIOUS_OTP_EXPIRATION":               http.StatusTooManyRequests,
	"EMPTY_ORG":                                      http.StatusBadRequest,
	"INVALID_REQ":                                    http.StatusBadRequest,
	"REG_FRST":                                       http.StatusForbidden,
	"WEAK_PIN":                                       http.StatusBadRequest,
	"NOT_FOUND":                                      http.StatusNotFound,
	"DEVICE_ID_REQUIRED":                             http.StatusBadRequest,
	"NO_LINKED_DEVICES":                              http.StatusNotFound,
	"COULD_NOT_UNLINK_DEVICE":                        http.StatusInternalServerError,
	"INCOMPLETE_USER_INFO":                           http.StatusBadRequest,
	"INVALID_ACTION_TYPE":                            http.StatusBadRequest,
	"ACTION_NOT_FOUND":                               http.StatusNotFound,
	"PENDING_REQUEST_EXISTS":                         http.StatusConflict,
	"MISSING_REQUIRED_HEADERS":                       http.StatusBadRequest,
	"DEVICE_LOOKUP_FAILED":                           http.StatusInternalServerError,
	"MISSING_REQUIRED_FIELDS":                        http.StatusBadRequest,
	"MISSING_OTP":                                    http.StatusBadRequest,
	"INVALID_INPUT_PARAMETERS":                       http.StatusBadRequest,
	"SAME_PIN":                                       http.StatusBadRequest,
	"PIN_IN_HISTORY":                                 http.StatusBadRequest,
	"ERROR_SETTING_PIN":                              http.StatusInternalServerError,
	"OTP_CREATION_FAILED":                            http.StatusInternalServerError,
	"DEVICE_NOT_FOUND":                               http.StatusNotFound,
	"DEVICE_FOUND":                                   http.StatusOK,
	"INVALID_PIN":                                    http.StatusBadRequest,
	"PIN_RESET_SESSION_EXPIRED":                      http.StatusBadRequest,
	"PIN_RESET_SESSION_NOT_FOUND":                    http.StatusNotFound,
	"PIN_RESET_SESSION_ALREADY_VERIFIED":             http.StatusBadRequest,
	"TOKEN_GENERATION_FAILED":                        http.StatusInternalServerError,
	"MAXIMUM_AMOUNT_REQUIRED_FOR_OPEN_METHOD":        http.StatusBadRequest,
	"MINIMUM_AMOUNT_REQUIRED_FOR_OTP_AND_PIN_METHOD": http.StatusBadRequest,
	"EITHER_MINIMUM_OR_MAXIMUM_AMOUNT_REQUIRED_FOR_PIN_METHOD":   http.StatusBadRequest,
	"INVALID_OBJECT_ID_FORMAT":                                   http.StatusBadRequest,
	"AUTH_TIER_NOT_FOUND":                                        http.StatusNotFound,
	"AUTH_TIER_ALREADY_EXISTS":                                   http.StatusBadRequest,
	"AUTH_TIER_VALIDATION_FAILED":                                http.StatusBadRequest,
	"AUTH_TIER_UPDATE_FAILED":                                    http.StatusInternalServerError,
	"AUTH_TIER_INSERT_FAILED":                                    http.StatusInternalServerError,
	"AUTH_TIER_NOT_FOUND_FOR_ID":                                 http.StatusNotFound,
	"AUTH_TIER_FETCH_FAILED":                                     http.StatusInternalServerError,
	"AUTH_TIER_APPROVE_FAILED":                                   http.StatusInternalServerError,
	"AUTH_TIER_REJECT_FAILED":                                    http.StatusInternalServerError,
	"AUTH_TIER_NOT_FOUND_FOR_ID_AND_DEPARTMENT":                  http.StatusNotFound,
	"AUTH_TIER_FETCH_FAILED_FOR_ID_AND_DEPARTMENT":               http.StatusInternalServerError,
	"AUTH_TIER_APPROVE_FAILED_FOR_ID_AND_DEPARTMENT":             http.StatusInternalServerError,
	"AUTH_TIER_REJECT_FAILED_FOR_ID_AND_DEPARTMENT":              http.StatusInternalServerError,
	"AUTH_TIER_UPDATE_FAILED_FOR_ID_AND_DEPARTMENT":              http.StatusInternalServerError,
	"AUTH_TIER_INSERT_FAILED_FOR_ID_AND_DEPARTMENT":              http.StatusInternalServerError,
	"FAILED_TO_UPDATE_OPEN_AUTH_TIER":                            http.StatusInternalServerError,
	"FAILED_TO_UPDATE_PIN_AUTH_TIER":                             http.StatusInternalServerError,
	"FAILED_TO_UPDATE_OTP_AND_PIN_AUTH_TIER":                     http.StatusInternalServerError,
	"FAILED_TO_UPDATE_OTP_AND_OPEN_AUTH_TIER":                    http.StatusInternalServerError,
	"FAILED_TO_UPDATE_OPEN_AUTH_TIER_FOR_ID":                     http.StatusInternalServerError,
	"MIN_AMOUNT_CANNOT_BE_GREATER_THAN_PIN_MIN_AMOUNT":           http.StatusBadRequest,
	"PIN_AUTHIER_NOT_FOUND":                                      http.StatusNotFound,
	"PIN_MIN_AMOUNT_CANNOT_BE_LESS_THAN_OPEN_MIN_AMOUNT":         http.StatusBadRequest,
	"MAX_PIN_AMOUNT_SHOULD_BE_GREATER_THAN_PIN_MIN_AMOUNT":       http.StatusBadRequest,
	"OPEN_AUTHIER_NOT_FOUND":                                     http.StatusNotFound,
	"OPEN_TIER_MAX_AMOUNT_CANNOT_BE_GREATER_THAN_PIN_MAX_AMOUNT": http.StatusBadRequest,
	"PENDING_CPS_ACTION_PRESENT":                                 http.StatusConflict,
	"FAILED_TO_GET_FAYDA_ACCOUNT":                                http.StatusInternalServerError,
	"FAILED_TO_GET_CUSTOMER_ACCOUNT":                             http.StatusInternalServerError,
	"FAILED_TO_CREATE_CPS_ACTION":                                http.StatusInternalServerError,
	"FAILED_TO_UPDATE_CPS_ACTION":                                http.StatusInternalServerError,
	"FAILED_TO_UPDATE_CUSTOMER_ACCOUNT":                          http.StatusInternalServerError,
	"FAILED_TO_UPDATE_CPS_ACTION_FOR_ID_AND_DEPARTMENT":          http.StatusInternalServerError,
	"FAILED_TO_GET_CPS_ACTION":                                   http.StatusInternalServerError,
	"FAILED_TO_GET_SERVICE":                                      http.StatusInternalServerError,
	"SERVICE_NOT_FOUND":                                          http.StatusNotFound,
	"FAILED_TO_FIND_CPS_ACTION":                                  http.StatusInternalServerError,
	"FAILED_TO_UPDATE_SERVICE_FEE":                               http.StatusInternalServerError,
	"FAILED_TO_REJECT_SERVICE_FEE_UPDATE":                        http.StatusInternalServerError,
	"REQUIRED_FIELDS_MISSING":                                    http.StatusBadRequest,
	"TIERS_REQUIRED":                                             http.StatusBadRequest,
	"TIERS_FIRST_MIN_ZERO":                                       http.StatusBadRequest,
	"TIERS_MIN_MUST_EQUAL_PREV_MAX":                              http.StatusBadRequest,
	"TIERS_MAX_MUST_INCREASE":                                    http.StatusBadRequest,
	"TIERS_ABOVE_AMOUNT_MISMATCH":                                http.StatusBadRequest,
	"NOT_IMPLEMENTED":                                            http.StatusNotImplemented,
	"INPUT_TOO_LONG":                                             http.StatusBadRequest,
	"INPUT_INVALID_CHARACTERS":                                   http.StatusBadRequest,
	"PAGE_NOT_FOUND":                                             http.StatusNotFound,
	"MISSING_OR_INVALID_IMAGE":                                   http.StatusBadRequest,
	"OPEN_MIN_GE_OPEN_MAX":                                       http.StatusBadRequest,
	"OPEN_MAX_GE_PIN_MAX":                                        http.StatusBadRequest,
	"PIN_MIN_GE_PIN_MAX":                                         http.StatusBadRequest,
	"PIN_MIN_LE_OPEN_MIN":                                        http.StatusBadRequest,
	"PIN_MAX_LE_PIN_MIN":                                         http.StatusBadRequest,
	"OTP_MIN_GE_PIN_MIN":                                         http.StatusBadRequest,
	"COLOR_ALREADY_EXISTED":                                      http.StatusBadRequest,
	"TIER_AUTH_NOT_FOUND":                                        http.StatusNotFound,
	"FAILED_TO_GET_AUTH_TIER":                                    http.StatusInternalServerError,
	"REQUIRED_TITLE":                                             http.StatusBadRequest,
	"REQUIRED_DESCRIPTION":                                       http.StatusBadRequest,
	"TITLE_TOO_LONG":                                             http.StatusBadRequest,
	"DESCRIPTION_TOO_LONG":                                       http.StatusBadRequest,
	"UNSUPPORTED_REQUEST_ACTION":                                 http.StatusBadRequest,
	"RESOURCE_ALREADY_ENABLED":                                   http.StatusBadRequest,
	"RESOURCE_ALREADY_DISABLED":                                  http.StatusBadRequest,
	"KYC_LEVEL_REQUIRED":                                         http.StatusBadRequest,
	"INVALID_KYC_LEVEL":                                          http.StatusBadRequest,
	"UNEXPECTED_DATABASE_ERROR":                                  http.StatusInternalServerError,
	"USER_CODE_ALREADY_EXIST":                                    http.StatusConflict,
	"PHONE_NUMBER_EXISTS":                                        http.StatusConflict,
	"EMAIL_ALREADY_EXISTS":                                       http.StatusConflict,
	"USERNAME_ALREADY_EXISTS":                                    http.StatusConflict,
	"WALLET_INFORMATION_ALREADY_EXISTST":                         http.StatusBadRequest,
	"NO_DATA_PROVIDED_FOR_UPDATE":                                http.StatusBadRequest,
	"CONTENT_TYPE_MUST_BE_FORM":                                  http.StatusUnsupportedMediaType,
	"CONTENT_TYPE_MUST_BE_JSON":                                  http.StatusUnsupportedMediaType,
	"EXPIRE_DATE_REQUIRED":                                       http.StatusBadRequest,
	"START_DATE_REQUIRED":                                        http.StatusBadRequest,
	"INVALID_PAYLOAD":                                            http.StatusBadRequest,
	"GENERAL_DB_QUERY_FAILED":                                    http.StatusInternalServerError,
	"ERROR_WHILE_CHECKING_PENDING_ACTION":                        http.StatusInternalServerError,

	// Auth
	"AUTH_USER_NOT_FOUND":               http.StatusNotFound,
	"AUTH_USER_DISABLED":                http.StatusForbidden,
	"AUTH_INVALID_PASSWORD":             http.StatusUnauthorized,
	"AUTH_USER_HAS_NO_PASSWORD":         http.StatusBadRequest,
	"AUTH_TOO_MANY_ATTEMPTS":            http.StatusTooManyRequests,
	"AUTH_INVALID_OTP":                  http.StatusBadRequest,
	"AUTH_EXPIRED_OTP":                  http.StatusBadRequest,
	"AUTH_USER_RESET_PASSWORD_REQUIRED": http.StatusForbidden,
	"AUTH_USER_ALREADY_EXISTS":          http.StatusBadRequest,
	"LOGIN_PROHIBITED_FOR_15_MIN":       http.StatusTooManyRequests,
	"INCORRECT_PASSWORD":                http.StatusUnauthorized,
	"OLD_PASSWORD_SAME_AS_NEW":          http.StatusBadRequest,
	"PIN_OLY_DIG":                       http.StatusBadRequest,
	"PIN_LIMIT":                         http.StatusBadRequest,
	"PIN_REDANDANT":                     http.StatusBadRequest,
	"PIN_SEQ":                           http.StatusBadRequest,
	"FAILD_TO_GEN_TOKEN":                http.StatusInternalServerError,
	"FAILD_TO_RESET_PASS":               http.StatusInternalServerError,
	"FAILD_VALIDATION":                  http.StatusBadRequest,
	"OLD_PIN_MISMATCH":                  http.StatusBadRequest,
	"INVALID_BEARER":                    http.StatusUnauthorized,
	"USE_RIGHT_AUTH":                    http.StatusUnauthorized,
	"INVALID_CLAIM":                     http.StatusUnauthorized,
	"INVALID_TOKEN_DATAu":               http.StatusUnauthorized,
	"UNABLE_TO_DYCRYPT_TOKEN":           http.StatusUnauthorized,
	"FAILED_LOGIN":                      http.StatusUnauthorized,
	"ACCOUNT_LOCKED":                    http.StatusForbidden,
	"PENDING_ACTION_ALREADY_EXISTS":     http.StatusBadRequest,
	"NO_DOCUMENT_FOUND":                 http.StatusNotFound,
	"ACCESS_TOKEN_REQUIRED":             http.StatusUnauthorized,

	// Transaction
	"TRANSACTION_NOT_FOUND": http.StatusNotFound,
	"INSUFFICIENT_FUNDS":    http.StatusPaymentRequired,
	"COMMISSION_NOT_FOUND":  http.StatusNotFound,
	"CUSTOMER_NOT_FOUND":    http.StatusNotFound,

	// Account
	"PHONE_LOOKUP_FAILED":           http.StatusBadGateway,
	"API_REQUEST_FAILED":            http.StatusBadGateway,
	"ACCOUNT_NOT_FOUND":             http.StatusNotFound,
	"USER_PHONE_EXISTS":             http.StatusConflict,
	"GENERAL_SIF_GENERATION_FAILED": http.StatusInternalServerError,
	"GENERAL_DB_UPDATE_FAILED":      http.StatusInternalServerError,
	"GENERAL_DB_INSERT_FAILED":      http.StatusInternalServerError,
	"CORE_ACCOUNT_NOT_FOUND":        http.StatusNotFound,
	"CORE_ACCOUNT_TYPE_ERROR":       http.StatusBadRequest,
	"ACCOUNT_LOOKUP_FAILED":         http.StatusBadGateway,
	"ACCOUNT_DETAIL_FAILED":         http.StatusBadGateway,
	"INVALID_CORE_RESPONSE":         http.StatusBadGateway,
	"ACCOUNT_ALREADY_LINKED":        http.StatusBadRequest,
	"ACCOUNT_BLOCKED":               http.StatusForbidden,

	// OTP
	"INVALID_OTP":            http.StatusBadRequest,
	"EXPIRED_OTP":            http.StatusBadRequest,
	"EMAIL_IN_USE":           http.StatusConflict,
	"USER_ALREADY_HAS_EMAIL": http.StatusBadRequest,
	"OTP_GENERATION_FAILED":  http.StatusInternalServerError,
	"USER_KYC_LEVEL_ZERO":    http.StatusBadRequest,
	"USER_KYC_LEVEL_WRONG":   http.StatusBadRequest,

	// File
	"INVALID_FORM":      http.StatusBadRequest,
	"NO_FILE":           http.StatusBadRequest,
	"FILE_TOO_LARGE":    http.StatusRequestEntityTooLarge,
	"INVALID_FILE_TYPE": http.StatusBadRequest,
	"UPLOAD_FAILED":     http.StatusInternalServerError,

	// Branch
	"BRANCH_NOT_FOUND":                           http.StatusNotFound,
	"BRANCH_DISABLED":                            http.StatusBadRequest,
	"FAILED_TO_FETCH_BRANCHES":                   http.StatusInternalServerError,
	"INVALID_BRANCH_ID":                          http.StatusBadRequest,
	"INVALID_LOCATION_FILTER":                    http.StatusBadRequest,
	"BRANCH_ID_REQUIRED":                         http.StatusBadRequest,
	"BRANCH_DISABLE_ACTION_ALREADY_EXISTS":       http.StatusBadRequest,
	"BRANCH_DISABLE_MULTI_ACTION_ALREADY_EXISTS": http.StatusConflict,

	// Department
	"DEPARTMENT_NOT_FOUND":      http.StatusNotFound,
	"DEPARTMENT_ALREADY_EXISTS": http.StatusBadRequest,
	"DEPARTMENT_CODE_REQUIRED":  http.StatusBadRequest,
	"DEPARTMENT_NAME_REQUIRED":  http.StatusBadRequest,
	"PORTAL_CARDS_INVALID":      http.StatusBadRequest,

	// unidentified key
	"UNIDENTIFIED_KEY": http.StatusInternalServerError,

	// Banks
	"BANKS_NOT_FOUND":                   http.StatusNotFound,
	"BANK_FETCH_FAILED":                 http.StatusInternalServerError,
	"INVALID_BANK_NAME":                 http.StatusBadRequest,
	"MISSING_BANK_NAME":                 http.StatusBadRequest,
	"MISSING_BANK_CODE":                 http.StatusBadRequest,
	"MISSING_BANK_BIC":                  http.StatusBadRequest,
	"BANK_ALREADY_EXISTS":               http.StatusBadRequest,
	"BLOCKED_ACTION_USER_ALREADY_EXIST": http.StatusBadRequest,

	// Action
	"FAILED_TO_CREATE_ACTION":  http.StatusInternalServerError,
	"FAILED_TO_UPDATE_ACTION":  http.StatusInternalServerError,
	"FAILED_TO_FETCH_ACTION":   http.StatusInternalServerError,
	"ACTION_NOT_PENDING":       http.StatusBadRequest,
	"FAILED_TO_UPDATE_SERVICE": http.StatusInternalServerError,
	"FAILED_TO_UPDATE_CAP_MIN": http.StatusInternalServerError,

	// CPS
	"CPS_USER_NOT_FOUND":          http.StatusNotFound,
	"ACTION_HAS_ALREADY_APPROVED": http.StatusBadRequest,
	"ACTION_HAS_ALREADY_REJECTED": http.StatusBadRequest,
	"NO_PENDING_ACTION_FOUND":     http.StatusNotFound,

	// AD
	"AD_NOT_FOUND":                             http.StatusNotFound,
	"AD_ALREADY_EXISTS":                        http.StatusBadRequest,
	"AD_CREATION_FAILED":                       http.StatusInternalServerError,
	"AD_UPDATE_FAILED":                         http.StatusInternalServerError,
	"AD_DELETION_FAILED":                       http.StatusInternalServerError,
	"CITY_NOT_FOUND":                           http.StatusNotFound,
	"PENDING_ACTION_ALREADY_EXIST":             http.StatusBadRequest,
	"CITY_ALREADY_BLOCKED":                     http.StatusBadRequest,
	"PENDING_DISTRICT_ACTION_ALREADY_EXIST":    http.StatusBadRequest,
	"DISTRICT_ALREADY_BLOCKED":                 http.StatusBadRequest,
	"DISTRICT_NOT_FOUND":                       http.StatusNotFound,
	"REGION_NOT_FOUND":                         http.StatusNotFound,
	"BRANCH_ALREADY_BLOCKED":                   http.StatusBadRequest,
	"REGION_CODE_AND_NAME_REQUIRED":            http.StatusBadRequest,
	"BLOCK_REGION_ALREADY_PROCESSED":           http.StatusBadRequest,
	"BANK_ID_REQUIRED":                         http.StatusConflict,
	"DEPARTMEN_REQUIRED":                       http.StatusConflict,
	"BANK_ALREADY_DISABLED":                    http.StatusBadRequest,
	"BANK_ALREADY_ENABLE":                      http.StatusBadRequest,
	"FAILED_TO_GET_DISTRICT":                   http.StatusNotFound,
	"BRANCH_ALREADY_DISABLE":                   http.StatusBadRequest,
	"BANK_NAME_ALREADY_EXIST":                  http.StatusBadRequest,
	"BANK_BIC_CODE_ALREADY_EXIST":              http.StatusBadRequest,
	"BIC_CODE_ALREADY_EXIST":                   http.StatusBadRequest,
	"BANK_ALREADY_CREATED_WITH_THIS_PARAMETER": http.StatusConflict,

	// Permission
	"NO_PERMISSION_CATEGORY_FOUND":                http.StatusNotFound,
	"INVALID_PERMISSION_CATEGORY_ID":              http.StatusBadRequest,
	"ONE_OR_MORE_PERMISSION_CATEGORIES_NOT_FOUND": http.StatusNotFound,
	"INVALID_PERMISSION_GROUP_ID":                 http.StatusBadRequest,
	"NO_PERMISSION_GROUP_FOUND":                   http.StatusNotFound,
	"ONE_OR_MORE_PERMISSION_GROUPS_NOT_FOUND":     http.StatusNotFound,
}

func getStatusForErrorKey(key string) int {
	if status, ok := errorKeyToStatus[key]; ok {
		return status
	}

	for _, group := range []common.ErrorGroup{
		common.DefineError.General,
		common.DefineError.Auth,
		common.DefineError.User,
		common.DefineError.Transaction,
		common.DefineError.Account,
		common.DefineError.OTP,
		common.DefineError.File,
		common.DefineError.Branch,
		common.DefineError.Bank,
		common.DefineError.Department,
		common.DefineError.Action,
		common.DefineError.Wallet,
		common.DefineError.AD,
	} {
		if _, ok := group[key]; ok {
			// If the error key is defined in any group but not in the map, default to 400 Bad Request.
			return http.StatusBadRequest
		}
	}

	// If the error key is not found anywhere, return 500 Internal Server Error.
	return http.StatusInternalServerError
}

func formatValidationErrors(ve validation.Errors) map[string]string {
	result := make(map[string]string)
	for field, fieldErr := range ve {
		def := lookupErrorDefinition(fieldErr.Error())
		result[field] = def.Message
	}
	return result
}

func SendErrorResponse(w http.ResponseWriter, errorKey any, statusCode int, additionalData map[string]any) {
	var (
		code              string
		message           string
		errorKeyString    string
		errors            any
		isValidationError bool
	)

	switch v := errorKey.(type) {
	case string:
		def := lookupErrorDefinition(v)
		code = def.Code
		message = def.Message
		errorKeyString = v

	case error:
		if ve, ok := v.(validation.Errors); ok {
			code = "GEN_113"
			message = "One or more required fields are invalid."
			errors = formatValidationErrors(ve)
			errorKeyString = "GEN_113"
			isValidationError = true
		} else {
			def := lookupErrorDefinition(v.Error())
			code = def.Code
			message = def.Message
			errorKeyString = def.Code
		}

	default:
		code = "GEN_UNKNOWN"
		message = fmt.Sprintf("%v", v)
		errorKeyString = "GEN_UNKNOWN"
	}

	status := httpStatusOrDefault(statusCode, errorKeyString, isValidationError)
	response := map[string]any{
		"status":  status,
		"message": message,
		"data": map[string]any{
			"code": code,
		},
	}

	if errors != nil {
		response["errors"] = errors
	}

	maps.Copy(response, additionalData)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Printf("Failed to write error response: %v\n", err)
	}
}

func lookupErrorDefinition(key string) common.ErrorDefinition {
	errorGroups := []common.ErrorGroup{
		common.DefineError.General,
		common.DefineError.Auth,
		common.DefineError.User,
		common.DefineError.Transaction,
		common.DefineError.Account,
		common.DefineError.OTP,
		common.DefineError.File,
		common.DefineError.Branch,
		common.DefineError.Bank,
		common.DefineError.Department,
		common.DefineError.Wallet,
		common.DefineError.Action,
		common.DefineError.AD,
		common.DefineError.Permission,
	}

	for _, group := range errorGroups {
		if def, ok := group[key]; ok {
			return def
		}
	}

	return common.ErrorDefinition{
		Code:    key,
		Message: key,
	}
}

func httpStatusOrDefault(providedStatus int, errorCode string, isValidation bool) int {
	if isValidation {
		return http.StatusBadRequest
	}
	if providedStatus != 0 {
		return providedStatus
	}
	return getStatusForErrorKey(errorCode)
}

// func HandleServiceError(w http.ResponseWriter, err error) {
// 	if def, ok := err.(common.ErrorDefinition); ok {
// 		// Use the most appropriate status code for the error code, not always 400.
// 		status := getStatusForErrorKey(def.Code)
// 		SendErrorResponse(w, def.Code, status, nil)
// 		return
// 	}
// 	SendErrorResponse(w, "GEN_004", http.StatusInternalServerError, nil)
// }
