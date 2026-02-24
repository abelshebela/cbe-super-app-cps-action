package sqlc

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

type VaultAmountTierParams struct {
	VaultCategoryID string
	MinAmount       float64
	MaxAmount       float64
	IsActive        sql.NullBool
	Page            sql.NullInt64 `json:"page"`
	Limit           sql.NullInt64 `json:"limit"`
}

type VaultAmountTier struct {
	ID              string          `json:"id"`
	VaultCategoryID string          `json:"vault_category_id"`
	MinAmount       decimal.Decimal `json:"min_amount"`
	MaxAmount       decimal.Decimal `json:"max_amount"`
	Interest        decimal.Decimal `json:"interest"`
	IsActive        sql.NullBool    `json:"is_active"`
	CreatedAt       time.Time       `json:"created_at,omitempty"`
	UpdatedAt       time.Time       `json:"updated_at,omitempty"`
	DeletedAt       *time.Time      `json:"deleted_at,omitempty"`
	IsDeleted       sql.NullBool    `json:"is_deleted,omitempty"`
	TotalCount      int64           `json:"total_count,omitempty"`
}
