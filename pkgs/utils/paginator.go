package utils

import (
	"net/http"
	"strconv"
)

// ExtractPaginator extracts page and limit from the request
func ExtractPaginator(r *http.Request) (limit, offset int64, err error) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := int64(1)
	limit = int64(10)

	if pageStr != "" {
		page, err = strconv.ParseInt(pageStr, 10, 64)
		if err != nil {
			return 0, 0, err
		}
	}

	if limitStr != "" {
		limit, err = strconv.ParseInt(limitStr, 10, 64)
		if err != nil {
			return 0, 0, err
		}
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset = (page - 1) * limit
	return limit, offset, nil
}
