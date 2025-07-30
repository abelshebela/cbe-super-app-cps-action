package utils

import (
	"cbe-super-app-member-users/pkgs/entities"
	"cbe-super-app-member-users/pkgs/entities/enums"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	mathrand "math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	common "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"go.mongodb.org/mongo-driver/v2/bson"
	// "cbe-super-app-member-users/pkgs/config"
	// entities "cbe-super-app-member-users/internal/domain/users"
)

func ErrorStatus(errorInfo string) int {
	var status int
	switch errorInfo {
	case "USER_NOT_FOUND":
		status = http.StatusNotFound
	case "INVALID_PIN":
		status = http.StatusUnauthorized
	case "ACCOUNT_BLOCKED":
		status = http.StatusForbidden
	case "ACCOUNT_DELETED":
		status = http.StatusGone
	case "DEVICE_NOT_LINKED":
		status = http.StatusForbidden
	case "TOO_MANY_LOGIN_ATTEMPTS":
		status = http.StatusTooManyRequests
	case "INVALID_PHONE_NUMBER":
		status = http.StatusBadRequest
	case "INVALID_DEVICE_UUID":
		status = http.StatusBadRequest
	default:
		status = http.StatusInternalServerError
	}

	return status
}

func OTPGenerator(length uint8) string {
	numberic := "0123456789"
	r := mathrand.New(mathrand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)
	for i := range result {
		result[i] = numberic[r.Intn(len(numberic))]
	}

	return string(result)
}
func BaseResponseMaker(res map[string]interface{}, w http.ResponseWriter, message string, statusCode interface{}) {

	response := make(map[string]interface{})
	if res == nil {
		res = make(map[string]interface{})
	}

	response["data"] = res
	response["status"] = statusCode
	response["message"] = message

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UntrustedInstallationResponse(w http.ResponseWriter) {
	returndata := make(map[string]interface{})

	returndata["message"] = common.DefineError.General["UNTRUSTED"].Message
	returndata["code"] = common.DefineError.General["UNTRUSTED"].Code
	w.WriteHeader(http.StatusForbidden)
	ResponseMaker(returndata, w)
}
func HeaderRequirement(r *http.Request, additionl []string) (string, string, string, string, string, map[string]interface{}) {
	var headerData = make(map[string]interface{})
	platform := r.Header.Get("platform")
	appVersion := r.Header.Get("app_version")
	devideuuid := r.Header.Get("device_uuid")
	sourceapp := r.Header.Get("source_app")
	installationdate := r.Header.Get("installation_date")
	for _, v := range additionl {
		headerData[v] = r.Header.Get(v)
	}

	return platform, appVersion, devideuuid, sourceapp, installationdate, headerData
}
func GetRandomArbitrary() (string, error) {
	const min = 100000
	const max = 999999
	var n int
	for {
		b := make([]byte, 4)
		if _, err := rand.Read(b); err != nil {
			return "", err
		}
		n = int(b[0])<<24 | int(b[1])<<16 | int(b[2])<<8 | int(b[3])
		n = min + (n % (max - min + 1))
		if n >= min && n <= max {
			break
		}
	}
	return fmt.Sprintf("%06d", n), nil
}

func FormatPhoneNumber(phoneNumber string) string {
	phoneNumber = strings.TrimSpace(phoneNumber)
	if strings.HasPrefix(phoneNumber, "0") {
		return "+251" + phoneNumber[1:]
	} else if strings.HasPrefix(phoneNumber, "9") || strings.HasPrefix(phoneNumber, "7") {
		return "+251" + phoneNumber
	} else if strings.HasPrefix(phoneNumber, "+") {
		return phoneNumber
	} else if strings.HasPrefix(phoneNumber, "251") {
		return "+" + phoneNumber
	}
	return phoneNumber
}
func UserContext(ctx context.Context) (entities.User, error) {
	// ctx = context.WithValue(ctx, constant.ContextKey("user"), userPayload)
	if ctx.Err() != nil {
		return entities.User{}, ctx.Err()
	}
	var userPayload entities.User

	// Helper to get string value from context
	getStr := func(key string) string {
		val := ctx.Value(ContextKey(key))
		if v, ok := val.(string); ok {
			return v
		}
		return ""
	}

	idStr := getStr("user_id")
	fmt.Println("User ID in context: app ", idStr)
	if idStr == "" {
		return entities.User{}, errors.New("user_id not found in context or not a string")
	}
	objID, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		return entities.User{}, errors.New("invalid user_id format")
	}
	userPayload.ID = objID
	userPayload.UserCode = getStr("user_code")
	userPayload.FullName = getStr("full_name")
	userPayload.PhoneNumber = getStr("phone_number")
	userPayload.Email = getStr("user_email")
	userPayload.Realm = enums.Realm(getStr("user_realm"))
	userPayload.MemberType = enums.MemberType(getStr("ifb_member"))
	userPayload.DeviceUUID = getStr("device_uuid")

	return userPayload, nil

}

func CheckLoginThrottle(attempts uint8, lastAttempt time.Time) error {
	elapsed := time.Since(lastAttempt)
	var waiting time.Duration
	switch {
	case attempts >= 5:
		waiting = 10 * time.Minute
	case attempts == 4:
		waiting = 5 * time.Minute
	case attempts == 3:
		waiting = 2 * time.Minute
	default:
		return nil
	}

	if elapsed < waiting {
		remaining := waiting - elapsed
		return fmt.Errorf("Too many login attempts. Please wait %s before trying again.", remaining.Truncate(time.Second))
	}
	return nil
}

func LocalEncryptPassword(password string, dataType string, userSalt string, action string, env *config.VaultConfig) (string, string, error) {

	if env == nil {
		return "", "", errors.New("env is nil")
	}
	var signedPass, salt string
	if dataType == "password" {
		salt, _ = GenerateSalt(20)
		signedPass, _ = SignWithHS256(password, salt)
	} else {
		signedPass = password
	}

	if action == "login" || action == "change" {
		salt = userSalt
		signedPass, _ = SignWithHS256(password, userSalt)
	}

	key := []byte(env.Key)
	iv := []byte(env.IV)
	if len(key) != 32 || len(iv) != aes.BlockSize {
		return "", salt, errors.New("invalid key or IV size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", salt, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	padLen := aes.BlockSize - len(signedPass)%aes.BlockSize
	padding := strings.Repeat(string(byte(padLen)), padLen)
	padded := []byte(signedPass + padding)
	encrypted := make([]byte, len(padded))
	mode.CryptBlocks(encrypted, padded)

	return hex.EncodeToString(encrypted), salt, nil
}

func LocalDecryptPassword(encryptedHex string, env *config.VaultConfig) (string, error) {

	//env, _ := config.Load()

	key := []byte(env.Key)
	iv := []byte(env.IV)

	if len(key) != 32 || len(iv) != aes.BlockSize {
		return "", errors.New("invalid key or IV size")
	}
	encrypted, err := hex.DecodeString(encryptedHex)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	if len(encrypted)%aes.BlockSize != 0 {
		return "", errors.New("invalid encrypted data length")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(encrypted))
	mode.CryptBlocks(decrypted, encrypted)
	// Remove PKCS#7 padding
	padLen := int(decrypted[len(decrypted)-1])
	if padLen > aes.BlockSize || padLen == 0 {
		return "", errors.New("invalid padding")
	}
	return string(decrypted[:len(decrypted)-padLen]), nil
}

func GenerateSalt(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func SignWithHS256(data string, saltHex string) (string, error) {
	key, err := hex.DecodeString(saltHex)
	if err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil)), nil
}

// func CheckOTPExpiration(otpData entities.OTP) time.Duration {

// 	expiresAtStr := otpData.ExpiresAt.Format(time.RFC3339Nano)
// 	expirationTime, _ := time.Parse(time.RFC3339, expiresAtStr)

// 	currentTime := time.Now().UTC()
// 	timeRemaining := expirationTime.Sub(currentTime)

//		return timeRemaining
//	}
func ResponseMaker(res map[string]interface{}, w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func ValidateInputNoSpecialChars(input string) error {
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9 ]*$`)

	if !validPattern.MatchString(input) {
		return errors.New("input contains special characters")
	}

	return nil
}

type StrictValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type StrictValidationResult struct {
	Valid  bool                    `json:"valid"`
	Errors []StrictValidationError `json:"errors"`
}

var validate = validator.New()

func ValidateStrict(input []byte, target interface{}) StrictValidationResult {
	json.Unmarshal(input, target)
	err := validate.Struct(target)
	result := StrictValidationResult{Valid: true}
	if err != nil {
		result.Valid = false
		for _, e := range err.(validator.ValidationErrors) {
			result.Errors = append(result.Errors, StrictValidationError{
				Field: e.Field(),
				Error: e.Tag(),
			})
		}
	}
	return result
}
