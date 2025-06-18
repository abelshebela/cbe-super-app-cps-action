package utils

import (
	"cbe-super-app-member-auth/pkg/common"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func InvalidInputResponse(w http.ResponseWriter) {

	returndata := make(map[string]interface{})
	returndata["message"] = common.DefineError.General["INVALID_INPUT"].Message
	returndata["code"] = common.DefineError.General["INVALID_INPUT"].Code
	w.WriteHeader(http.StatusBadRequest)
	ResponseMaker(returndata, w)
}

func ResetLockedRespose(w http.ResponseWriter) {
	returndata := make(map[string]interface{})

	returndata["accountLocked"] = true
	returndata["message"] = common.DefineError.Auth["AUTH_USER_RESET_PASSWORD_REQUIRED"].Message
	returndata["code"] = common.DefineError.Auth["AUTH_USER_RESET_PASSWORD_REQUIRED"].Code
	w.WriteHeader(http.StatusUnauthorized)
	ResponseMaker(returndata, w)
}

func UUIDGenerator() *uuid.UUID {
	id := uuid.New()
	return &id
}

func ObjectIDGenerator() (*bson.ObjectID, error) {
	id, err := bson.ObjectIDFromHex(bson.NewObjectID().Hex())
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func ContainsString(data []string, otp string) bool {

	if data == nil {
		return false
	}

	for _, v := range data {
		if v == otp {
			return true
		}
	}
	return false
}

func MaintainPINHistory(currentPIN string, pinHistory []string) []string {
	const maxHistory = 5

	newHistory := append([]string{currentPIN}, pinHistory...)

	unique := []string{}
	seen := map[string]bool{}
	for _, pin := range newHistory {
		if !seen[pin] {
			unique = append(unique, pin)
			seen[pin] = true
		}
	}

	if len(unique) > maxHistory {
		unique = unique[:maxHistory]
	}
	return unique
}
func ParseInstallationDate(installationDate interface{}) (interface{}, error) {
	switch v := installationDate.(type) {
	case string:
		// Try to parse as RFC3339 or common date formats
		t, err := time.Parse(time.RFC3339, v)
		if err == nil {
			return t, nil
		}
		// Try another common format
		t, err = time.Parse("2006-01-02", v)
		if err == nil {
			return t, nil
		}
		return nil, fmt.Errorf("unable to parse date string: %v", v)
	case time.Time:
		return v, nil
	default:
		return nil, fmt.Errorf("unsupported type: %T", installationDate)
	}

}

func SendSMS(msg string, phone string) {
	fmt.Printf("Sending SMS to %s: %s\n", phone, msg)
}
func AccountLockedRespose(w http.ResponseWriter) {
	returndata := make(map[string]interface{})

	returndata["accountLocked"] = true
	returndata["message"] = common.DefineError.Auth["AUTH_USER_RESET_PASSWORD_REQUIRED"].Message
	returndata["code"] = common.DefineError.Auth["AUTH_USER_RESET_PASSWORD_REQUIRED"].Code
	w.WriteHeader(http.StatusUnauthorized)
	ResponseMaker(returndata, w)
}

func InvalidJsonPayloadResponse(w http.ResponseWriter) {
	returndata := make(map[string]interface{})
	returndata["message"] = common.DefineError.General["INVALID_JSON_PAYLOAD"].Message
	returndata["code"] = common.DefineError.General["INVALID_JSON_PAYLOAD"].Code
	w.WriteHeader(http.StatusNotAcceptable)
	ResponseMaker(returndata, w)
}

func ValidationResponse(validation []string, w http.ResponseWriter) {
	returndata := make(map[string]interface{})

	returndata["errors"] = validation
	returndata["message"] = common.DefineError.Auth["FAILD_VALIDATION"]
	w.WriteHeader(http.StatusBadRequest)
	ResponseMaker(returndata, w)
}

func UnauthorizedResponse(w http.ResponseWriter) {
	returndata := make(map[string]interface{})

	returndata["message"] = common.DefineError.General["UNAUTHORIZED"].Message
	returndata["code"] = common.DefineError.General["UNAUTHORIZED"].Code
	w.WriteHeader(http.StatusUnauthorized)
	ResponseMaker(returndata, w)
}

func UntrustedInstallationResponse(w http.ResponseWriter) {
	returndata := make(map[string]interface{})

	returndata["message"] = common.DefineError.General["UNTRUSTED"].Message
	returndata["code"] = common.DefineError.General["UNTRUSTED"].Code
	w.WriteHeader(http.StatusForbidden)
	ResponseMaker(returndata, w)
}

func PinLimitResponse(w http.ResponseWriter) {
	returndata := make(map[string]interface{})

	returndata["message"] = common.DefineError.Auth["PIN_LIMIT"].Message
	returndata["code"] = common.DefineError.Auth["PIN_LIMIT"].Code
	w.WriteHeader(http.StatusForbidden)
	ResponseMaker(returndata, w)
}

func UserNotFoundResponse(w http.ResponseWriter) {
	returndata := make(map[string]interface{})

	returndata["message"] = common.DefineError.Auth["AUTH_USER_NOT_FOUND"].Message
	returndata["code"] = common.DefineError.Auth["AUTH_USER_NOT_FOUND"].Code
	w.WriteHeader(http.StatusUnauthorized)
	ResponseMaker(returndata, w)
}

func OtpCheckFailedResponse(w http.ResponseWriter) {
	returndata := make(map[string]interface{})

	returndata["message"] = common.DefineError.Auth["OTP_CHECK_FAILED"].Message
	returndata["code"] = common.DefineError.Auth["OTP_CHECK_FAILED"].Code
	w.WriteHeader(http.StatusUnauthorized)
	ResponseMaker(returndata, w)
}
func UserDeviceNotFoundResponse(deviceUUID string, w http.ResponseWriter) {
	returndata := make(map[string]interface{})

	returndata["device_uuid"] = deviceUUID
	returndata["message"] = common.DefineError.Auth["AUTH_USER_NOT_FOUND"].Message
	returndata["code"] = common.DefineError.Auth["AUTH_USER_NOT_FOUND"].Code
	w.WriteHeader(http.StatusUnauthorized)
	ResponseMaker(returndata, w)
}
func DeviceNotFoundResponse(w http.ResponseWriter) {
	returndata := make(map[string]interface{})

	returndata["message"] = common.DefineError.Auth["DEVICE_NOT_FOUND"].Message
	returndata["code"] = common.DefineError.Auth["DEVICE_NOT_FOUND"].Code
	w.WriteHeader(http.StatusUnauthorized)
	ResponseMaker(returndata, w)
}

func InvalidOtpResponse(w http.ResponseWriter) {
	returndata := make(map[string]interface{})

	returndata["message"] = common.DefineError.Auth["AUTH_INVALID_OTP"].Message
	returndata["code"] = common.DefineError.Auth["AUTH_INVALID_OTP"].Code
	w.WriteHeader(http.StatusBadRequest)
	ResponseMaker(returndata, w)
}

func FailedToGenerateTokenResponse(w http.ResponseWriter) {
	returndata := make(map[string]interface{})

	returndata["message"] = common.DefineError.Auth["FAILD_TO_GEN_TOKEN"].Message
	returndata["code"] = common.DefineError.Auth["FAILD_TO_GEN_TOKEN"].Code
	w.WriteHeader(http.StatusUnauthorized)
	ResponseMaker(returndata, w)
}

func UnhandledServerErrorResponse(w http.ResponseWriter) {
	returndata := make(map[string]interface{})

	returndata["message"] = common.DefineError.General["UNHANDLED_SERVER_ERROR"].Message
	returndata["code"] = common.DefineError.General["UNHANDLED_SERVER_ERROR"].Code
	w.WriteHeader(http.StatusInternalServerError)
	ResponseMaker(returndata, w)

}
