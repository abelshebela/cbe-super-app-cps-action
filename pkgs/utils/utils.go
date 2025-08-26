package utils

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	mathrand "math/rand"
	"mime/multipart"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"


	"cbe-super-app-cps-action/internal/constants"
	customErr "cbe-super-app-cps-action/internal/constants/errors"
	"cbe-super-app-cps-action/internal/constants/response"
	"cbe-super-app-cps-action/internal/constants/types"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
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
		return customErr.ErrTooManyAttempt
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

func BuildMongoFilterWithKeys(input map[string]interface{}, allowedKeys []string) bson.M {
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
			nested := BuildMongoFilterWithKeys(v, allowedKeys)
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

// nonEmptyString returns if non-empty, otherwise fallback
func NonEmptyString(s, fallback string) string {
	if s != "" {
		return s
	}
	return fallback
}

var allowedChars = "a-zA-Z0-9\\s._-"

func NoSpecialChars(value any) error {
	str, ok := value.(string)
	if !ok {
		return validation.NewError("validation", "invalid type")
	}
	str = strings.TrimSpace(str)
	if str == "" {
		return nil
	}

	re := regexp.MustCompile("^[" + allowedChars + "]+$")
	if !re.MatchString(str) {
		return validation.NewError("validation", "contains invalid characters")
	}
	return nil
}

func TrimWhiteSpace(value interface{}) error {
	if s, ok := value.(string); ok {
		if strings.TrimSpace(s) == "" {
			return errors.New("value cannot be empty or whitespace")
		}
	}
	return nil
}

func GetParam(r *http.Request, key string) (string, bool) {
	value := chi.URLParam(r, key)
	if value == "" {
		return "", false
	}
	return value, true
}

func ExtractID(w http.ResponseWriter, r *http.Request) (string, error) {
	// TODO: Implement proper parameter extraction when utils.GetParam is available
	// For now, use a simple approach
	id := chi.URLParam(r, "id")
	if id == "" {
		response.SendErrorResponse(w, customErr.ErrIdEmpty)
		return "", customErr.ErrIdEmpty
	}
	return id, nil
}
