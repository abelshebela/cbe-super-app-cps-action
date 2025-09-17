package utils

import (
	"net/http"
	"strconv"
	"strings"

	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

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

	return &constant.MongoFilter{
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

func BuildMongoFilter(input map[string]interface{}) bson.M {
	return BuildMongoFilterWithKeys(input, nil)
}

func BuildMongoFilterWithKeys(input map[string]interface{}, allowedKeys []string) bson.M {
	filter := bson.M{}

	allowedMap := make(map[string]bool)
	if allowedKeys != nil {
		for _, key := range allowedKeys {
			allowedMap[key] = true
		}
	}

	for key, value := range input {
		if value == nil || value == "" {
			continue
		}

		if allowedKeys != nil && !allowedMap[key] {
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

func BuildMongoFilterWithHandlers(input map[string]interface{}, allowedKeys []string, handlers map[string]func(interface{}) interface{}) bson.M {
	filter := bson.M{}

	allowedMap := make(map[string]bool)
	if allowedKeys != nil {
		for _, key := range allowedKeys {
			allowedMap[key] = true
		}
	}

	for key, value := range input {
		if value == nil || value == "" {
			continue
		}

		if allowedKeys != nil && !allowedMap[key] {
			continue
		}

		if handler, exists := handlers[key]; exists {
			filter[key] = handler(value)
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
			nested := BuildMongoFilterWithHandlers(v, allowedKeys, handlers)
			for nestedKey, nestedVal := range nested {
				filter[key+"."+nestedKey] = nestedVal
			}
		default:
			filter[key] = v
		}
	}

	return filter
}
