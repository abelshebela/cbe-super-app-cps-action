package utils

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"mime/multipart"

	// "math/rand"
	mathrand "math/rand"
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
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	shared_constant "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/constants"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const alphanumberic string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func IsValidImage(fileHeader *multipart.FileHeader) bool {
	var allowedMIMETypes = map[string]bool{
		"image/jpg":  true,
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
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

func IsValidVideo(fileHeader *multipart.FileHeader) bool {
	var allowedMIMETypes = map[string]bool{
		"video/mp4":        true,
		"video/x-msvideo":  true, // avi
		"video/quicktime":  true, // mov
		"video/x-matroska": true, // mkv
		"video/webm":       true,
	}

	if fileHeader.Size > 50*1024*1024 { // 50MB limit for video
		return false
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

		if fn, ok := handler[key]; ok {
			value = fn(value)
		}

		switch v := value.(type) {
		case string:
			if v != "" {
				filter[key] = v
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

func ExtractID(w http.ResponseWriter, r *http.Request) (string, error) {

	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorInvalidID, nil, nil)
		return "", errors.New(localization.ErrorInvalidID.Code)
	}
	return id, nil
}

func NonEmptyBool(newVal, oldVal bool) bool {
	if newVal != oldVal {
		return newVal
	}
	return oldVal
}

// nonEmptyAdvertFor returns the new value if non-empty, otherwise the old value
func NonEmptyAdvertFor(new, old shared_constant.AdvertFor) shared_constant.AdvertFor {
	if new != "" {
		return new
	}
	return old
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

var allowedChars = "a-zA-Z0-9\\s._@-"

func NoSpecialChars(value any) error {
	var str string
	switch v := value.(type) {
	case string:
		str = v
	case *string:
		if v == nil {
			return nil
		}
		str = *v
	default:
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
	// phoneNumber = strings.TrimSpace(phoneNumber)
	phoneNumber = strings.TrimSpace(phoneNumber)

	// Remove all non-digit and non-plus characters
	re := regexp.MustCompile(`[^\d\+]`)
	phoneNumber = re.ReplaceAllString(phoneNumber, "")

	if strings.HasPrefix(phoneNumber, "0") && len(phoneNumber) == 10 {
		phoneNumber = "251" + phoneNumber[1:]
	} else if strings.HasPrefix(phoneNumber, "9") && len(phoneNumber) == 9 {
		phoneNumber = "251" + phoneNumber
	} else if strings.HasPrefix(phoneNumber, "7") && len(phoneNumber) == 9 {
		phoneNumber = "251" + phoneNumber
	}

	if strings.HasPrefix(phoneNumber, "251") && len(phoneNumber) == 13 {
		return phoneNumber
	}

	return ""
	// Acceptable patterns: 2517XXXXXXXX or 2519XXXXXXXX (total 12 digits)
	// validRe := regexp.MustCompile(`^(251[79]\d{8})$`)
	// match := validRe.MatchString(phoneNumber)
	// if match {
	// 	return phoneNumber
	// }
	// return ""
}

func ThreeNamesMinLength(value interface{}) error {
	var name string
	switch v := value.(type) {
	case string:
		name = v
	case *string:
		if v == nil {
			return nil
		}
		name = *v
	default:
		return errors.New("invalid full_name")
	}

	parts := strings.Fields(name)
	if len(parts) != 3 {
		return errors.New("full_name must contain first, middle, and last name")
	}

	for _, p := range parts {
		if len(p) <= 3 {
			return errors.New("each of first, middle, and last name must be longer than 3 characters")
		}
	}

	return nil
}

func TrimWhiteSpace(value interface{}) error {
	var s string
	switch v := value.(type) {
	case string:
		s = v
	case *string:
		if v == nil {
			return nil
		}
		s = *v
	default:
		return nil // Non-string types are not trimmed
	}
	if strings.TrimSpace(s) == "" {
		return errors.New("value cannot be empty or whitespace")
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

// --- FUNCTION 2: CONVERSION ---

// DurationToMonths correctly converts a time.Duration back into the number of months.
// It uses the 365-day year constant from the parser and converts to years first,
// ensuring a duration generated by 'y' correctly rounds back to X*12 months.
func DurationToMonths(d any) int {
	var dur time.Duration

	switch v := d.(type) {
	case int64:
		dur = time.Duration(v)
	case time.Duration:
		dur = v
	default:
		return 0
	}

	hours := dur.Hours()

	// 1. Define the year constant based on the parser's 'y' approximation (365 days * 24 hours).
	const hoursPer365DayYear = 365.0 * 24.0

	// 2. Convert the duration to the number of 365-day years.
	years := hours / hoursPer365DayYear

	// 3. Convert years to months and use math.Round() to handle floating point issues.
	months := int(math.Round(years * 12.0))

	return months
}

func ParseToYears(s string) (float64, error) {
	if len(strings.TrimSpace(s)) < 2 {
		return 0, fmt.Errorf("invalid lock period format")
	}

	s = strings.TrimSpace(s)
	unit := strings.ToLower(s[len(s)-1:])       // last character
	valueStr := strings.TrimSpace(s[:len(s)-1]) // everything except unit

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number in lock period: %v", err)
	}

	const daysPerYear = 365.0

	switch unit {
	case "d":
		return value / daysPerYear, nil
	case "m":
		return value / 12.0, nil
	case "y":
		return value, nil
	default:
		return 0, fmt.Errorf("invalid unit in lock period: %s", unit)
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

func NullTimeToPtr(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
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

func LocalEncryptPassword(password string, dataType string, userSalt string, action string, cfg *config.VaultConfig) (string, string, error) {

	var signedPass, salt string
	if dataType == constants.Password {
		salt, _ = GenerateSalt(20)
		signedPass, _ = SignWithHS256(password, salt)
	} else if dataType == constants.Cred {
		signedPass, _ = SignWithHS256(password, cfg.JwtSecretKey)
	} else {
		signedPass = password
	}

	if action == constants.Login || action == constants.Change {
		salt = userSalt
		signedPass, _ = SignWithHS256(password, userSalt)
	}

	key := []byte(cfg.Key)
	iv := []byte(cfg.IV)
	if len(key) != 32 || len(iv) != aes.BlockSize {
		return "", salt, localization.ErrorInvalidKey
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

func NumbersOnly(value any) error {
	var str string
	switch v := value.(type) {
	case string:
		str = v
	case *string:
		if v == nil {
			return nil
		}
		str = *v
	default:
		return validation.NewError("validation", "unsupported type")
	}
	re := regexp.MustCompile(`^\d+$`)
	if !re.MatchString(str) {
		return validation.NewError("validation", "contains invalid characters")
	}
	return nil
}

func TraceLogger(ctx context.Context, key, spanName, serviceType, serviceName string) (context.Context, trace.Span) {
	tracer := otel.Tracer(key)
	ctx, span := tracer.Start(ctx, spanName)
	span.SetAttributes(attribute.String(serviceType, serviceName))

	return ctx, span

}
