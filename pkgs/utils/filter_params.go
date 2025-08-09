package utils

import (
	"net/http"
	"strconv"
	"strings"

	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// ExtractFilterParams extracts pagination and filter params from the HTTP request query.
func ExtractFilterParams(r *http.Request) *constant.Filter {
	query := r.URL.Query()

	page := constant.DefaultPage
	if pageInt, err := strconv.Atoi(query.Get("page")); err == nil && pageInt > 0 {
		page = pageInt
	}
	perPage := constant.DefaultPerPage
	if perPageInt, err := strconv.Atoi(query.Get("per_page")); err == nil &&
		perPageInt > 0 {
		perPage = perPageInt
	}

	return &constant.Filter{
		Page:    page,
		PerPage: perPage,
		Search:  query.Get("search"),
		Filters: query.Get("filter"),
	}
}

// BuildMongoFilter constructs a MongoDB filter from the provided input map.
// MongoFilter
func ExtractMongoFilterParams(r *http.Request) *constant.MongoFilter {
	query := r.URL.Query()

	page := constant.DefaultPage
	if v := query.Get("page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			page = n
		}
	}

	perPage := constant.DefaultPerPage
	if v := query.Get("per_page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			perPage = n
		}
	}

	filters := make(map[string]interface{}, len(query))
	for k, v := range query {
		if len(v) == 0 || !strings.HasPrefix(k, "filter[") || !strings.HasSuffix(k, "]") {
			continue
		}
		key := k[7 : len(k)-1]
		filters[key] = v[0]
	}

	return &constant.MongoFilter{
		Page:    page,
		PerPage: perPage,
		Search:  query.Get("search"),
		Filters: filters,
	}
}

func BuildMongoFilter(input map[string]interface{}) bson.M {
	return BuildMongoFilterWithValidation(input, nil)
}

func BuildMongoFilterWithValidation(input map[string]interface{}, validKeys []string) bson.M {
	filter := bson.M{}

	var validKeysMap map[string]bool
	if validKeys != nil {
		validKeysMap = make(map[string]bool, len(validKeys))
		for _, key := range validKeys {
			validKeysMap[key] = true
		}
	}

	for key, value := range input {
		if value == nil || value == "" {
			continue
		}

		if validKeysMap != nil && !validKeysMap[key] {
			continue
		}

		switch v := value.(type) {
		case string:
			filter[key] = bson.M{"$regex": v, "$options": "i"}

		case []interface{}:
			filter[key] = bson.M{"$in": v}

		case map[string]interface{}:
			nested := BuildMongoFilterWithValidation(v, validKeys)
			for nestedKey, nestedVal := range nested {
				filter[key+"."+nestedKey] = nestedVal
			}

		default:
			filter[key] = v
		}
	}

	return filter
}
