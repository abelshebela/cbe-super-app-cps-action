package utils

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	// "math/rand"
	mathrand "math/rand"
	"mime/multipart"
	"net/http"

	"os"

	"regexp"
	"strconv"
	"strings"
	"time"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/types"

	"cbe-super-app-cps-action/internal/constants/localization"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const alphanumberic string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

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

// nonEmptyString returns if non-empty, otherwise fallback
func NonEmptyString(s, fallback string) string {
	if s != "" {
		return s
	}
	return fallback
}
func NonEmptyNotificationFor(newFor string, oldFor constants.NotificationFor) constants.NotificationFor {
	if newFor != "" {
		return constants.NotificationFor(newFor)
	}
	return oldFor
}

func ExtractID(w http.ResponseWriter, r *http.Request) (string, error) {

	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorInvalidID, nil, nil)
		return "", errors.New(localization.ErrorInvalidID.Code)
	}
	return id, nil
}

func NonEmptyBool(newVal, oldVal bool) bool {
	// Handles updates correctly (req can explicitly override old value)
	if newVal != oldVal {
		return newVal
	}
	return oldVal
}

// MergeProductCodes merges by ProductCode value (not index)
func MergeProductCodes(newPCs []types.ProductCode, oldPCs []types.ProductCode) []types.ProductCode {
	if len(newPCs) == 0 {
		return oldPCs
	}

	// Map old codes by ProductCode for quick lookup
	oldMap := make(map[string]types.ProductCode)
	for _, pc := range oldPCs {
		oldMap[pc.ProductCode] = pc
	}

	updated := make([]types.ProductCode, 0, len(newPCs))
	for _, pc := range newPCs {
		if existing, found := oldMap[pc.ProductCode]; found {
			updated = append(updated, types.ProductCode{
				ID:             existing.ID,
				BranchType:     constants.BranchType(pc.BranchType),
				ProductCode:    NonEmptyString(pc.ProductCode, existing.ProductCode),
				VATCode:        NonEmptyString(pc.VATCode, existing.VATCode),
				ServiceFeeCode: NonEmptyString(pc.ServiceFeeCode, existing.ServiceFeeCode),
			})
		} else {
			// New ProductCode → assign new ID
			updated = append(updated, types.ProductCode{
				ID:             utils.RandomGenerator(20),
				BranchType:     constants.BranchType(pc.BranchType),
				ProductCode:    pc.ProductCode,
				VATCode:        pc.VATCode,
				ServiceFeeCode: pc.ServiceFeeCode,
			})
		}
	}

	return updated
}

func NonZeroTime(t, fallback time.Time) time.Time {
	if !t.IsZero() {
		return t
	}
	return fallback
}

func NonZeroUint64(n, fallback uint64) uint64 {
	if n != 0 {
		return n
	}
	return fallback
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

func RandomGenerator(length uint8) string {
	if length <= 0 {
		panic("length must be greater than 0")
	}

	entropy := fmt.Sprintf("%d-%d", time.Now().UnixNano(), os.Getpid())
	hash := sha256.Sum256([]byte(entropy))
	seed := int64(binary.LittleEndian.Uint64(hash[:8]))
	r := mathrand.New(mathrand.NewSource(seed))

	result := make([]byte, length)
	for i := range result {
		result[i] = alphanumberic[r.Intn(len(alphanumberic))]
	}

	return string(result)
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
func FormatPhoneNumber(phoneNumber string) string {
	phoneNumber = strings.TrimSpace(phoneNumber)

	// Remove all non-digit and non-plus characters
	re := regexp.MustCompile(`[^\d\+]`)
	phoneNumber = re.ReplaceAllString(phoneNumber, "")

	if strings.HasPrefix(phoneNumber, "+2510") {
		phoneNumber = "+251" + phoneNumber[5:]
	} else if strings.HasPrefix(phoneNumber, "2510") {
		phoneNumber = "+251" + phoneNumber[4:]
	} else if strings.HasPrefix(phoneNumber, "0") && len(phoneNumber) == 10 {
		phoneNumber = "+251" + phoneNumber[1:]
	} else if strings.HasPrefix(phoneNumber, "9") && len(phoneNumber) == 9 {
		phoneNumber = "+251" + phoneNumber
	} else if strings.HasPrefix(phoneNumber, "7") && len(phoneNumber) == 9 {
		phoneNumber = "+251" + phoneNumber
	} else if strings.HasPrefix(phoneNumber, "251") {
		phoneNumber = "+" + phoneNumber
	}

	if strings.HasPrefix(phoneNumber, "+251") && len(phoneNumber) == 13 {
		return phoneNumber
	}

	return ""
}

func TrimWhiteSpace(value interface{}) error {
	if s, ok := value.(string); ok {
		if strings.TrimSpace(s) == "" {
			return errors.New("value cannot be empty or whitespace")
		}
	}
	return nil

}

func JsonUnmarshal[T any](data any) (*T, error) {

	var jsonData *T
	byte, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(byte, &jsonData); err != nil {
		return nil, err
	}

	return jsonData, nil
}

func ExtraSpaceRemover(s string) string {
	return strings.TrimSpace(strings.Join(strings.Split(s, " "), " "))
}

func ParseLockPeriod(s string) (time.Duration, error) {
	if len(s) < 2 {
		return 0, fmt.Errorf("invalid lock period format")
	}

	unit := s[len(s)-1]      // last character: 'd', 'm', 'y'
	valueStr := s[:len(s)-1] // number part
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("invalid number in lock period: %v", err)
	}

	switch strings.ToLower(string(unit)) {
	case "d":
		return time.Duration(value) * 24 * time.Hour, nil
	case "m":
		return time.Duration(value*30) * 24 * time.Hour, nil // approximate 1 month = 30 days
	case "y":
		return time.Duration(value*365) * 24 * time.Hour, nil // approximate 1 year = 365 days
	default:
		return 0, fmt.Errorf("invalid unit in lock period: %s", string(unit))
	}
}
func NullStringToPtrLike(ns sql.NullString) *string {
	if !ns.Valid {
		return (*string)(nil)
	}
	s := "%" + strings.ToUpper(ns.String) + "%"
	return &s
}

func NullStringToPtr(v any) *string {
	switch x := v.(type) {
	case string:
		if x == "" {
			return nil
		}
		return &x
	case sql.NullString:
		if !x.Valid {
			return nil
		}
		return &x.String
	default:
		return nil
	}
}

func NullInt64ToPtr(ni sql.NullInt64) *int64 {
	if !ni.Valid {
		return (*int64)(nil)
	}
	i := ni.Int64
	return &i
}
func NullBoolToInt(nb sql.NullBool) interface{} {
	if nb.Valid {
		if nb.Bool {
			return 1
		}
		return 0
	}
	return nil
}
func NullBoolToIntPtr(nb sql.NullBool) *int {
	if !nb.Valid {
		return nil
	}
	if nb.Bool {
		i := 1
		return &i
	}
	i := 0
	return &i
}

func NullBoolToBool(nb sql.NullBool) bool {
	return nb.Valid && nb.Bool
}

type PaginatedResponse[T any] struct {
	Items       []T   `json:"items"`
	Page        int64 `json:"page"`
	Limit       int64 `json:"limit"`
	Total       int64 `json:"total"`
	TotalPages  int64 `json:"total_pages"`
	HasNextPage bool  `json:"has_next_page"`
	HasPrevPage bool  `json:"has_prev_page"`
}

func NewPaginatedResponse[T any](data []T, page, limit, total int64) PaginatedResponse[T] {
	if limit <= 0 {
		limit = 50
	}
	if page <= 0 {
		page = 1
	}
	totalPages := (total + limit - 1) / limit

	return PaginatedResponse[T]{
		Items:       data,
		Page:        page,
		Limit:       limit,
		Total:       total,
		TotalPages:  totalPages,
		HasNextPage: page < totalPages,
		HasPrevPage: page > 1,
	}
}
