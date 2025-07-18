package utils

import (
	"net/http"
	"strconv"
)

type PaginationMeta struct {
	TotalDocs     int64 `json:"totalDocs"`
	Limit         int   `json:"limit"`
	TotalPages    int   `json:"totalPages"`
	Page          int   `json:"page"`
	PagingCounter int   `json:"pagingCounter"`
	HasPrevPage   bool  `json:"hasPrevPage"`
	HasNextPage   bool  `json:"hasNextPage"`
	PrevPage      *int  `json:"prevPage,omitempty"`
	NextPage      *int  `json:"nextPage,omitempty"`
}
type PaginatedResponse[T any] struct {
	Data T              `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

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

// BuildPaginationMeta constructs PaginationMeta based on total documents, current page, and limit
func BuildPaginationMeta(totalDocs int64, page, limit int) PaginationMeta {
	totalPages := int((totalDocs + int64(limit) - 1) / int64(limit)) // ceil(totalDocs / limit)
	hasPrev := page > 1
	hasNext := page < totalPages
	skip := (page - 1) * limit

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

	return PaginationMeta{
		TotalDocs:     totalDocs,
		Limit:         limit,
		TotalPages:    totalPages,
		Page:          page,
		PagingCounter: skip + 1,
		HasPrevPage:   hasPrev,
		HasNextPage:   hasNext,
		PrevPage:      prevPage,
		NextPage:      nextPage,
	}
}
