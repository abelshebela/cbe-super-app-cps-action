package utils

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/types"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func ParseTime(date string) time.Time {
	parsedTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return time.Time{}
	}
	return parsedTime
}

func ParseUserContext(r *http.Request) (types.UserContext, error) {
	userContext := ExtractUserContext(r)
	if IsIncomplete(userContext) {
		return types.UserContext{}, fmt.Errorf(constants.IncompleteUserInfo)
	}

	return userContext, nil
}

func ExtractUserContext(r *http.Request) types.UserContext {
	get := func(key string) string {
		val, _ := r.Context().Value(constants.ContextKey(key)).(string)
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

	page := constants.DefaultPage
	if v := query.Get("page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			page = n
		}
	}

	perPage := constants.DefaultPerPage
	if v := query.Get("per_page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			perPage = n
		}
	}

	filters := make(map[string]interface{}, len(query))

	for k, v := range query {
		if len(v) == 0 || k == "page" || k == "per_page" || k == "search" {
			continue
		}

		if strings.HasPrefix(k, "filter[") && strings.HasSuffix(k, "]") {
			key := k[7 : len(k)-1]
			filters[key] = parseValue(v[0])
			continue
		}

		if !isReserved(k) {
			filters[k] = parseValue(v[0])
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

	if parsed, err := strconv.Atoi(value); err == nil {
		return parsed
	}

	if parsed, err := strconv.ParseFloat(value, 64); err == nil {
		return parsed
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
		prefix     = "CBE_"
		codeLen    = 20
		charset    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		randomPart = codeLen - len(prefix)
	)
	b := make([]byte, randomPart)
	seed := time.Now().UnixNano()
	for i := range b {
		seed = seed*1664525 + 1013904223 // simple LCG for more randomness
		b[i] = charset[seed%int64(len(charset))]
	}
	return prefix + string(b)
}
