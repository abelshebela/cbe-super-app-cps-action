package sqlc

import (
	"database/sql"
)

type VaultCategory struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	CategoryType string       `json:"category_name"`
	CoverImage   string       `json:"cover_image"`
	IsActive     bool         `json:"is_active"`
	CreatedAt    sql.NullTime `json:"created_at"`
	UpdatedAt    sql.NullTime `json:"updated_at"`
	DeletedAt    sql.NullTime `json:"deleted_at"`
	IsDeleted    bool         `json:"is_deleted"`
	TotalCount   int64        `json:"total_count"`
}
