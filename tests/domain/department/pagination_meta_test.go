package department_test

import (
	"testing"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	"github.com/stretchr/testify/assert"
)

func TestPaginationMetaCalculation(t *testing.T) {
	tests := []struct {
		name      string
		totalDocs int64
		page      int
		limit     int
		expected  utils.PaginationMeta
	}{
		{
			name:      "First page with multiple pages",
			totalDocs: 25,
			page:      1,
			limit:     10,
			expected: utils.PaginationMeta{
				TotalDocs:     25,
				Limit:         10,
				TotalPages:    3,
				Page:          1,
				PagingCounter: 1,
				HasPrevPage:   false,
				HasNextPage:   true,
				PrevPage:      nil,
				NextPage:      intPtr(2),
			},
		},
		{
			name:      "Middle page",
			totalDocs: 25,
			page:      2,
			limit:     10,
			expected: utils.PaginationMeta{
				TotalDocs:     25,
				Limit:         10,
				TotalPages:    3,
				Page:          2,
				PagingCounter: 11,
				HasPrevPage:   true,
				HasNextPage:   true,
				PrevPage:      intPtr(1),
				NextPage:      intPtr(3),
			},
		},
		{
			name:      "Last page",
			totalDocs: 25,
			page:      3,
			limit:     10,
			expected: utils.PaginationMeta{
				TotalDocs:     25,
				Limit:         10,
				TotalPages:    3,
				Page:          3,
				PagingCounter: 21,
				HasPrevPage:   true,
				HasNextPage:   false,
				PrevPage:      intPtr(2),
				NextPage:      nil,
			},
		},
		{
			name:      "Single page",
			totalDocs: 5,
			page:      1,
			limit:     10,
			expected: utils.PaginationMeta{
				TotalDocs:     5,
				Limit:         10,
				TotalPages:    1,
				Page:          1,
				PagingCounter: 1,
				HasPrevPage:   false,
				HasNextPage:   false,
				PrevPage:      nil,
				NextPage:      nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.BuildPaginationMeta(tt.totalDocs, tt.page, tt.limit)

			assert.Equal(t, tt.expected.TotalDocs, result.TotalDocs)
			assert.Equal(t, tt.expected.Limit, result.Limit)
			assert.Equal(t, tt.expected.TotalPages, result.TotalPages)
			assert.Equal(t, tt.expected.Page, result.Page)
			assert.Equal(t, tt.expected.PagingCounter, result.PagingCounter)
			assert.Equal(t, tt.expected.HasPrevPage, result.HasPrevPage)
			assert.Equal(t, tt.expected.HasNextPage, result.HasNextPage)
			assert.Equal(t, tt.expected.PrevPage, result.PrevPage)
			assert.Equal(t, tt.expected.NextPage, result.NextPage)
		})
	}
}

func intPtr(i int) *int {
	return &i
}
