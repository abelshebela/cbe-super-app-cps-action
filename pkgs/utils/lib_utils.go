package utils

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"context"
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var counter uint64
var reHex24 = regexp.MustCompile(`(?i)[0-9a-f]{24}`)

func ParseOracleBool(v interface{}) (bool, bool) {
	switch b := v.(type) {
	case bool:
		return b, true
	case int:
		return b != 0, true
	case int32:
		return b != 0, true
	case int64:
		return b != 0, true
	case float64:
		return b != 0, true
	case string:
		switch strings.ToLower(strings.TrimSpace(b)) {
		case "true", "1", "yes", "y":
			return true, true
		case "false", "0", "no", "n":
			return false, true
		}
	}
	return false, false
}

func BoolToOracleNumber(v bool) int {
	if v {
		return 1
	}
	return 0
}

func Contains(list []string, v string) bool {
	for _, s := range list {
		if s == v {
			return true
		}
	}
	return false
}

func FormatTime(v any) string {
	switch t := v.(type) {
	case time.Time:
		if t.IsZero() {
			return ""
		}
		return t.Format(time.RFC3339)
	case *time.Time:
		if t == nil || t.IsZero() {
			return ""
		}
		return t.Format(time.RFC3339)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func IsActionInGroup(action constants.RequestAction, group string, RequestActionGroups map[string][]constants.RequestAction) bool {
	actions, exists := RequestActionGroups[group]
	if !exists {
		return false
	}

	for _, a := range actions {
		if a == action {
			return true
		}
	}
	return false
}

func ParseDateInput(s string) (time.Time, error) {
	formats := []string{
		time.RFC3339Nano,                // e.g. export handler FormatDateRangeToUTCStrings
		time.RFC3339,                    // 2026-01-05T07:10:33+00:00
		"2006-01-02T15:04:05.000Z07:00", // 2026-01-05T07:10:33.695+00:00
		"2006-01-02T15:04:05.999Z07:00", // milliseconds variant
		"2006-01-02T15:04:05Z07:00",     // without millis
		"2006-01-02T15:04:05",           // no timezone
		"2006-01-02",                    // date only
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse date: %s", s)
}

func GetRAListForUpdateAction(action constants.RequestAction, group string, RequestActionGroups map[string][]constants.RequestAction) []string {
	var RAUpdateList = []string{}
	actions, exists := RequestActionGroups[group]
	if !exists {
		return []string{}
	}

	for _, a := range actions {
		if strings.Contains(string(a), "UPDATE") || strings.Contains(string(a), "ENABLE") || strings.Contains(string(a), "DISABLE") {
			RAUpdateList = append(RAUpdateList, string(a))
			return RAUpdateList
		}
	}

	return RAUpdateList
}

func RemoveDuplicates(slice []string) []string {
	if slice == nil {
		return nil
	}

	seen := make(map[string]struct{}, len(slice)) // track seen elements
	result := make([]string, 0, len(slice))       // pre-allocate result

	for _, s := range slice {
		if _, ok := seen[s]; ok {
			continue // skip duplicates
		}
		seen[s] = struct{}{}
		result = append(result, s)
	}

	return result
}

func isHex(s string) bool {
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

func FirstHex24(s string) string {
	if s == "" {
		return ""
	}
	if len(s) == 24 && isHex(s) {
		return strings.ToLower(s)
	}
	m := reHex24.FindString(s)
	if m == "" {
		return ""
	}
	return strings.ToLower(m)
}
func NewNotificationID() string {
	// timestamp format: YYYYMMDDHHMMSS
	timestamp := time.Now().Format("20060102150405")
	return fmt.Sprintf("%s", timestamp)
}
func ParseTime(date string) time.Time {
	parsedTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return time.Time{}
	}
	return parsedTime
}

func UniqueIdGenerator() string {
	// timestamp format: YYYYMMDDHHMMSS
	timestamp := time.Now().Format("20060102150405")
	return fmt.Sprintf("%s", timestamp)
}

func ParseUserContext(r *http.Request) (types.UserContext, error) {
	userContext := ExtractUserContext(r)
	if IsIncomplete(userContext) {
		return types.UserContext{}, fmt.Errorf(constants.IncompleteUserInfo)
	}

	return userContext, nil
}

func ExtractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", errors.New("invalid authorization header format")
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
	if token == "" {
		return "", errors.New("empty bearer token")
	}

	return token, nil
}

func ExtractUserContext(r *http.Request) types.UserContext {
	get := func(key string) string {
		val, _ := r.Context().Value(constants.ContextKey(key)).(string)
		return val
	}

	getBool := func(key string) bool {
		val, _ := r.Context().Value(constants.ContextKey(key)).(bool)
		return val
	}

	return types.UserContext{
		IsErp:        getBool("is_erp"),
		UserCode:     get("user_code"),
		UserID:       get("user_id"),
		FullName:     get("full_name"),
		UserName:     get("username"),
		PhoneNumber:  get("phone_number"),
		Department:   get("department"),
		UserRole:     get("user_role"),
		CheckerIndex: get("role_checker_index"),
	}
}

func ExtractUserFromContext(ctx context.Context) types.UserContext {
	get := func(key string) string {
		val, _ := ctx.Value(constants.ContextKey(key)).(string)
		return val
	}
	getBool := func(key string) bool {
		val, _ := ctx.Value(constants.ContextKey(key)).(bool)
		return val
	}

	return types.UserContext{
		IsErp:       getBool("is_erp"),
		UserCode:    get("user_code"),
		UserID:      get("user_id"),
		FullName:    get("full_name"),
		UserName:    get("username"),
		PhoneNumber: get("phone_number"),
		Department:  get("department"),
		UserRole:    get("user_role"),
	}
}

func IsIncomplete(u types.UserContext) bool {
	return strings.TrimSpace(u.UserID) == "" || strings.TrimSpace(u.FullName) == "" || strings.TrimSpace(u.PhoneNumber) == "" || strings.TrimSpace(u.Department) == ""
}

func ExtractFilterParams(r *http.Request) *types.Filter {
	query := r.URL.Query()

	// --- Pagination defaults ---
	page := constants.DefaultPage
	if v := query.Get("page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			page = n
		} else {
			page = 1
		}
	} else {
		page = 1
	}

	perPage := constants.DefaultPerPage
	if v := query.Get("per_page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			perPage = n
		} else {
			perPage = 10
		}
	} else {
		perPage = 10
	}

	filters := make(map[string]interface{})

	for k, v := range query {
		if len(v) == 0 || isReserved(k) {
			continue
		}

		// handle both `filter[key]` and plain query params
		var key string
		if strings.HasPrefix(k, "filter[") && strings.HasSuffix(k, "]") {
			key = k[7 : len(k)-1]
		} else {
			key = k
		}

		// collect all values for this key
		if len(v) == 1 {
			// single value → let parseValue handle comma splits and types
			filters[key] = parseValue(v[0])
		} else {
			// multiple values → combine as []interface{}
			var arr []interface{}
			for _, val := range v {
				arr = append(arr, parseValue(val))
			}
			filters[key] = arr
		}
	}

	return &types.Filter{
		Page:    page,
		PerPage: perPage,
		Search:  query.Get("search"),
		Filters: filters,
	}
}

// StringFromFilterValue returns the first non-empty string from a filter value produced by
// ExtractFilterParams (plain string, or []interface{} when the same query key is repeated).
func StringFromFilterValue(v interface{}) (string, bool) {
	switch t := v.(type) {
	case string:
		s := strings.TrimSpace(t)
		return s, s != ""
	case []interface{}:
		for _, x := range t {
			if s, ok := x.(string); ok {
				s = strings.TrimSpace(s)
				if s != "" {
					return s, true
				}
			}
		}
	}
	return "", false
}

func isReserved(key string) bool {
	reserved := map[string]bool{
		"page": true, "per_page": true, "search": true, "sort": true, "order": true,
		"limit": true, "offset": true, "fields": true, "include": true, "exclude": true,
		"format": true, "callback": true, "pretty": true,
	}
	return reserved[key]
}

func parseValue(value string) interface{} {
	if value == "" {
		return ""
	}

	if value == "true" {
		return true
	}
	if value == "false" {
		return false
	}

	if isDigitString(value) {
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
			return parsed
		} else if errors.Is(err, strconv.ErrRange) {
			return value
		}
	} else if strings.ContainsAny(value, ".eE") {
		if parsed, err := strconv.ParseFloat(value, 64); err == nil {
			return parsed
		}
	}

	// Handle bracket notation: [VAL1,VAL2] — strip outer brackets before splitting
	stripped := value
	if strings.HasPrefix(stripped, "[") && strings.HasSuffix(stripped, "]") {
		stripped = stripped[1 : len(stripped)-1]
	}

	if strings.Contains(stripped, ",") {
		parts := strings.Split(stripped, ",")
		var result []interface{}
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				result = append(result, parseValue(part))
			}
		}
		if len(result) > 0 {
			return result
		}
	}

	return value
}

func isDigitString(value string) bool {
	if value == "" {
		return false
	}
	for _, r := range value {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func BuildPaginationMeta(totalDocs int64, page, limit int) types.PaginationMeta {
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}

	var totalPages int
	if totalDocs == 0 {
		totalPages = 0
		page = 1
	} else {
		totalPages = int((totalDocs + int64(limit) - 1) / int64(limit))
	}

	skip := (page - 1) * limit
	hasPrev := page > 1
	hasNext := page < totalPages

	var prevPage *int
	var nextPage *int

	if hasPrev {
		p := page - 1
		prevPage = &p
	}
	if hasNext {
		n := page + 1
		nextPage = &n
	}

	pagingCounter := 0
	if totalDocs > 0 {
		pagingCounter = skip + 1
	}

	return types.PaginationMeta{
		TotalDocs:     totalDocs,
		Limit:         limit,
		TotalPages:    totalPages,
		Page:          page,
		PagingCounter: pagingCounter,
		HasPrevPage:   hasPrev,
		HasNextPage:   hasNext,
		PrevPage:      prevPage,
		NextPage:      nextPage,
	}
}

// StringToObjectID converts a hex string to a bson.ObjectID.
// Returns the ObjectID and a boolean indicating if the conversion was successful.
func StringToObjectID(id string) (bson.ObjectID, bool) {
	if id == "" {
		return bson.ObjectID{}, false
	}
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return bson.ObjectID{}, false
	}
	return objID, true
}

// GenerateActionCode generates a unique action code in the format: SRM + YY + DDD + HHMMSS
// Example: SRM26216_143025 (year 2026, 216th day, 14:30:25)
func GenerateActionCode() string {
	const prefix = "SRM"

	now := time.Now()
	year := now.Format("06")                        // last 2 digits of year
	dayOfYear := fmt.Sprintf("%03d", now.YearDay()) // day of year zero-padded to 3 digits
	timeStr := now.Format("150405.000000")          // HHMMSSmmm (milliseconds)

	return prefix + year + dayOfYear + "_" + timeStr
}

func HandleMongoError(err error) (string, string) {
	// Defensive: nil error means no error
	if err == nil {
		return "", ""
	}

	if err.Error() == "mongo: no documents in result" || err.Error() == "no documents in result" {
		// Not found error
		return localization.ErrorResourceNotFound.Code, localization.ErrorResourceNotFound.Message
	}

	return localization.ErrorUnexpectedError.Code, localization.ErrorUnexpectedError.Message
}

// func GenerateCPSUserCode() string {
// 	// Numeric time part (last 9 digits of Unix nano for compactness)
// 	now := time.Now().UnixNano()
// 	timePart := fmt.Sprintf("%09d", now%1e9)

// 	// Random part
// 	const length = 6
// 	bytes := make([]byte, length)
// 	_, err := rand.Read(bytes)
// 	if err != nil {
// 		panic(err) // handle properly in production
// 	}

// 	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
// 	for i := range bytes {
// 		bytes[i] = charset[int(bytes[i])%len(charset)]
// 	}

// 	randomPart := string(bytes)

// 	return fmt.Sprintf("%s_%s", timePart, randomPart)
// }

func GenerateCPSUserCode() string {
	// Time part (last 6 digits for shorter length)
	now := time.Now().UnixNano()
	timePart := fmt.Sprintf("%06d", now%1e6)

	// Random part (4 chars instead of 6)
	const length = 4
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}

	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for i := range bytes {
		bytes[i] = charset[int(bytes[i])%len(charset)]
	}
	randomPart := string(bytes)

	// Short UUID (base32 encoded, trimmed)
	u := uuid.New()
	encoded := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(u[:])
	shortUUID := strings.ToLower(encoded[:6]) // take first 6 chars

	return fmt.Sprintf("%s_%s_%s", timePart, randomPart, shortUUID)
}
func GenerateBPSUserCode() string {
	const prefix = "BANKBPSUSER_"

	timestamp := time.Now().Format("20060102150405")

	return prefix + timestamp
}

func GenerateCustomerCode() string {
	const prefix = "CUST-"

	timestamp := time.Now().Format("20060102150405")

	return prefix + timestamp
}

func ParseDateString(dateStr string) (time.Time, error) {
	formats := []string{
		"02/01/2006",                // DD/MM/YYYY
		"01/02/2006",                // MM/DD/YYYY
		"2006-01-02",                // YYYY-MM-DD
		"2006-01-02T15:04:05Z07:00", // ISO format
		"2006-01-02T15:04:05",       // ISO format without timezone
		"02-01-2006",                // DD-MM-YYYY
		"01-02-2006",                // MM-DD-YYYY
	}

	for _, format := range formats {
		if parsed, err := time.Parse(format, dateStr); err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

func ToInterfaceSlice(strs []string) []interface{} {
	res := make([]interface{}, len(strs))
	for i, v := range strs {
		res[i] = v
	}
	return res
}
