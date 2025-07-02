package utils

import (
	"encoding/json"
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
	"CONFLICT_KEY":                     409,
	"INVALID_ID":                       400,
	"INVALID_JSON_PAYLOAD":             400,
	"UNHANDLED_SERVER_ERROR":           500,
	"INVALID_TOKEN_FORMAT":             400,
	"INVALID_TOKEN":                    401,
	"EXPIRED_TOKEN":                    401,
	"ERROR_CHANGING_PIN":               400,
	"ENCRYPTED_PAYLOAD_REQUIRED":       400,
	"INVALID_SERVER_CONFIGURATION":     500,
	"AUTH_HEADER_MISSING":              401,
	"INVALID_KEY_CONFIGURATION":        500,
	"INVALID_ENCRYPTED_PAYLOAD":        400,
	"DECRYPTION_ERROR":                 400,
	"SERVER_KEYS_NOT_CONFIGURED":       500,
	"INVALID_INPUT":                    400,
	"UNAUTHORIZED":                     401,
	"USER_REALM_NOT_FOUND":             403,
	"ACTION_NOT_ALLOWED":               403,
	"WAIT_FOR_PREVIOUS_OTP_EXPIRATION": 429,
	"EMPTY_ORG":                        400,
	"INVALID_REQ":                      400,
	"REG_FRST":                         403,
	"WEAK_PIN":                         400,
	"NOT_FOUND":                        404,
	"DEVICE_ID_REQUIRED":               400,
	"NO_LINKED_DEVICES":                404,
	"COULD_NOT_UNLINK_DEVICE":          500,

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
	"SAME_PIN":                          400,
	"PIN_IN_HISTORY":                    400,
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

	// OTP
	"INVALID_OTP":            400,
	"EXPIRED_OTP":            400,
	"EMAIL_IN_USE":           409,
	"USER_ALREADY_HAS_EMAIL": 409,
	// "WAIT_FOR_PREVIOUS_OTP_EXPIRATION": 429,
	"OTP_GENERATION_FAILED": 500,
	"USER_KYC_LEVEL_ZERO":   400,
	"USER_KYC_LEVEL_WRONG":  400,

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

	// if additionalData != nil {
	// 	if accountLocked, ok := additionalData["accountLocked"].(bool); ok {
	// 		response.AccountLocked = accountLocked
	// 	}
	// 	if auth, ok := additionalData["auth"].(bool); ok {
	// 		response.Auth = auth
	// 	}
	// 	if errors, exists := additionalData["errors"]; exists {
	// 		response.Errors = errors
	// 	}
	// }

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

/*
Example usage:

1. Basic error response:
   SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, nil)

2. Account locked error:
   SendErrorResponse(w, "AUTH_USER_RESET_PASSWORD_REQUIRED", http.StatusUnauthorized,
       map[string]interface{}{"accountLocked": true})

3. Validation error with details:
   SendErrorResponse(w, "FAILD_VALIDATION", http.StatusBadRequest,
       map[string]interface{}{"errors": validationErrors})

4. OTP error:
   SendErrorResponse(w, "AUTH_EXPIRED_OTP", http.StatusNotFound,
       map[string]interface{}{"auth": false})
*/
