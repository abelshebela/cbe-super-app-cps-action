package utils

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

var counter uint64
var reHex24 = regexp.MustCompile(`(?i)[0-9a-f]{24}`)

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
	return types.UserContext{
		UserCode:     get("user_code"),
		UserID:       get("user_id"),
		FullName:     get("full_name"),
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

	return types.UserContext{
		UserCode:    get("user_code"),
		UserID:      get("user_id"),
		FullName:    get("full_name"),
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
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
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

	if strings.Contains(value, ",") {
		parts := strings.Split(value, ",")
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

// GenerateActionCode generates a unique action code of length 20 with prefix "CBE_"
func GenerateActionCode() string {
	const (
		prefix  = "BANK_"
		codeLen = 20
		charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	)

	randomPart := codeLen - len(prefix)
	b := make([]byte, randomPart)

	// Seed once with high-resolution time
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range b {
		b[i] = charset[rnd.Intn(len(charset))]
	}

	return prefix + string(b)
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

func GenerateCPSUserCode() string {
	const (
		prefix  = "BANKCPSUSER_"
		codeLen = 15
		charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	)

	randomPart := codeLen - len(prefix)
	b := make([]byte, randomPart)

	// Seed once with high-resolution time
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range b {
		b[i] = charset[rnd.Intn(len(charset))]
	}

	return prefix + string(b)
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
