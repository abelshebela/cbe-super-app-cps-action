package utils

import (
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"net/http"
	"strconv"
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
