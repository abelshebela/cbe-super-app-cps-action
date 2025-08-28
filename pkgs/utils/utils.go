package utils

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	mathrand "math/rand"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"

	"errors"

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

func GenerateUsername(fullName string) string {
	fullName = strings.TrimSpace(fullName)
	parts := strings.Fields(fullName)
	if len(parts) < 2 {
		return strings.ToLower(strings.ReplaceAll(fullName, " ", ""))
	}
	first := strings.ToLower(parts[0])
	last := strings.ToLower(parts[len(parts)-1])
	clean := func(s string) string {
		var b strings.Builder
		for _, r := range s {
			if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
				b.WriteRune(r)
			}
		}
		return b.String()
	}
	first = clean(first)
	last = clean(last)
	r := mathrand.New(mathrand.NewSource(time.Now().UnixNano()))
	num := r.Intn(900) + 100 // 100-999
	return fmt.Sprintf("%s.%s%d", first, last, num)
}

func GenerateUserCode() string {
	const prefix = "CBEUSR-"

	// Generate a random number between 0 and 999999999999
	r := mathrand.New(mathrand.NewSource(time.Now().UnixNano()))
	number := r.Int63n(1000000000000) // 12 digits

	// Format with leading zeros to ensure 12 digits
	return fmt.Sprintf("%s%012d", prefix, number)
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

	// Ensure the range is valid and non-negative formathrand.Intn
	rangeSize := max - min + 1
	if rangeSize <= 0 {
		return ""
	}

	generatedNumber := min + mathrand.Intn(rangeSize)
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
	_, err := mathrand.Read(bytes)
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
		// remaining := waitDuration - elapsed
		// minutes := int(remaining.Minutes())
		return errors.New(localization.ErrorUserTooManyLoginAttempts.Code)
	}

	return nil
}

func FilterIdFor(id string) (bson.M, error) {
	if strings.HasPrefix(id, "ObjectID(\"") && strings.HasSuffix(id, "\")") {
		id = id[len("ObjectID(\"") : len(id)-2]
	}
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
		return nil, errors.New(localization.ErrorUserUnauthorized.Code)
	}

	fullName, ok := ctx.Value(constants.ContextKey("full_name")).(string)
	if !ok {
		log.Errorf("failed to get full name from context: %v", ok)
		return nil, errors.New(localization.ErrorUserUnauthorized.Code)
	}

	phoneNumber, ok := ctx.Value(constants.ContextKey("phone_number")).(string)
	if !ok {
		log.Errorf("failed to get full name from context", ok)
		return nil, errors.New(localization.ErrorUserUnauthorized.Code)
	}

	action, ok := ctx.Value(constants.ContextKey("action")).(string)
	if !ok {
		log.Errorf("failed to get action from context", ok)
		return nil, errors.New(localization.ErrorUserUnauthorized.Code)
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
		return "", errors.New(localization.ErrorInvalidRequest.Code)
	}
	return step, nil
}

func BuildMongoFilterWithKeys(input map[string]interface{}, allowedKeys []string, handler map[string]func(interface{}) interface{}) bson.M {
	filter := bson.M{}

	allowedMap := make(map[string]bool)
	for _, key := range allowedKeys {
		allowedMap[key] = true
	}

	for key, value := range input {
		if value == nil || value == "" {
			continue
		}

		if !allowedMap[key] {
			continue
		}

		switch v := value.(type) {
		case string:
			if v != "" {
				filter[key] = bson.M{"$regex": v, "$options": "i"}
			}
		case []interface{}:
			if len(v) > 0 {
				filter[key] = bson.M{"$in": v}
			}
		case map[string]interface{}:
			nested := BuildMongoFilterWithKeys(v, allowedKeys, handler)
			for nestedKey, nestedVal := range nested {
				filter[key+"."+nestedKey] = nestedVal
			}
		default:
			filter[key] = v
		}
	}

	return filter
}

func MapSlice[T any, R any](items []T, mapper func(T) R) []R {
	results := make([]R, len(items))
	for i, v := range items {
		results[i] = mapper(v)
	}
	return results
}

// nonEmptyString returns the new value if non-empty, otherwise the old value
func NonEmptyString(new, old string) string {
	if new != "" {
		return new
	}
	return old
}

// nonEmptyAdvertFor returns the new value if non-empty, otherwise the old value
func NonEmptyAdvertFor(new, old constants.AdvertFor) constants.AdvertFor {
	if new != "" {
		return new
	}
	return old
}

// nonEmptyAdvertDate returns the new date if non-zero, otherwise the old date
func NonEmptyAdvertDate(new, old types.AdvertDate) types.AdvertDate {
	result := old
	if !new.StartedAt.IsZero() {
		result.StartedAt = new.StartedAt
	}
	if !new.ExpiredAt.IsZero() {
		result.ExpiredAt = new.ExpiredAt
	}
	return result
}

func BindAction(source any, target any) error {
	bytes, err := json.Marshal(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, target)
}
