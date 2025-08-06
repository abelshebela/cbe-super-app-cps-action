package utils

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	mathrand "math/rand"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants"
	customErr "github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/errors"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/types"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func IsWeakPin(pin string) bool {
	// Check for repeated digits
	if strings.Count(pin, string(pin[0])) == len(pin) {
		return true
	}
	// Check for sequential patterns
	isAscending := true
	isDescending := true
	for i := 1; i < len(pin); i++ {
		if pin[i] != pin[i-1]+1 {
			isAscending = false
		}
		if pin[i] != pin[i-1]-1 {
			isDescending = false
		}
	}
	return isAscending || isDescending
}

func ValidateFullName(value interface{}) error {
	fullName := fmt.Sprintf("%s", value)
	name := strings.Split(fullName, " ")

	if len(name) != 2 {
		return fmt.Errorf("invalid full name")
	}
	if len(name[0]) < 3 || len(name[1]) < 3 {
		return fmt.Errorf("invalid full name")
	}
	return nil
}

func IsValidImage(fileHeader *multipart.FileHeader) bool {
	var allowedMIMETypes = map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
	}

	file, err := fileHeader.Open()
	if err != nil {
		return false
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return false
	}
	contentType := http.DetectContentType(buffer)
	return allowedMIMETypes[contentType]
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

func GenerateRandom(digit int) string {
	if digit <= 0 {
		return ""
	}

	min := intPow(10, digit-1)
	max := intPow(10, digit) - 1
	if digit == 1 {
		min = 0
	}

	// Ensure the range is valid and non-negative for rand.Intn
	rangeSize := max - min + 1
	if rangeSize <= 0 {
		return ""
	}

	generatedNumber := min + rand.Intn(rangeSize)
	result := strconv.Itoa(generatedNumber)
	return result
}

func intPow(a, b int) int {
	result := 1
	for i := 0; i < b; i++ {
		result *= a
	}
	return result
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

func CheckLoginThrottle(attempts uint8, lastAttempt time.Time) error {
	elapsed := time.Since(lastAttempt)
	var waitDuration time.Duration

	switch {
	case attempts >= 5:
		waitDuration = 10 * time.Minute
	case attempts == 4:
		waitDuration = 5 * time.Minute
	case attempts == 3:
		waitDuration = 2 * time.Minute
	default:
		return nil
	}

	if elapsed < waitDuration {
		remaining := waitDuration - elapsed
		return fmt.Errorf("Too many login attempts. Please wait %s before trying again.", remaining.Truncate(time.Second))
	}

	return nil
}

func FilterIdFor(id string) (bson.M, error) {
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return bson.M{
		"_id": objId,
	}, nil
}

func ExtractUserInfo(ctx context.Context, log utils.Logger) (*types.UserInfo, error) {
	userID, ok := ctx.Value(constants.ContextKey("user_id")).(string)
	if !ok {
		log.Errorf("failed to fetch user id from context: %v", ok)
		return nil, customErr.ErrBadRequest
	}

	fullName, ok := ctx.Value(constants.ContextKey("full_name")).(string)
	if !ok {
		log.Errorf("failed to get full name from context: %v", ok)
		return nil, customErr.ErrBadRequest
	}

	phoneNumber, ok := ctx.Value(constants.ContextKey("phone_number")).(string)
	if !ok {
		log.Errorf("failed to get full name from context", ok)
		return nil, customErr.ErrBadRequest
	}

	action, ok := ctx.Value(constants.ContextKey("action")).(string)
	if !ok {
		log.Errorf("failed to get action from context", ok)
		return nil, customErr.ErrBadRequest
	}

	return &types.UserInfo{
		UserID:      userID,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
		Action:      action,
	}, nil
}

func ExtractNextStep(ctx context.Context, log utils.Logger) (string, error) {
	step, ok := ctx.Value(constants.ContextKey("next_step")).(string)
	if !ok {
		log.Errorf("faile to get next step from context")
		return "", customErr.ErrBadRequest
	}
	return step, nil
}
