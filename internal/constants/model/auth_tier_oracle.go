package model

import (
	"cbe-super-app-cps-action/internal/constants"
	"time"
)

// AuthTierOracle represents a single active row in Oracle for amount-based auth.
// This replaces the Mongo model's ObjectID with an Oracle hex string ID.
type AuthTierOracle struct {
	ID string `sqlx:"id" json:"id"`

	Currency  constants.CurrencyType `sqlx:"currency" json:"currency"`
	MinAmount uint64                 `sqlx:"min_amount" json:"min_amount"`
	MaxAmount uint64                 `sqlx:"max_amount" json:"max_amount"`
	Method    constants.Method       `sqlx:"method" json:"method"`

	Enabled   int `sqlx:"enabled" json:"enabled"`
	IsDeleted int `sqlx:"is_deleted" json:"is_deleted"`

	CreatedAt    time.Time `sqlx:"created_at" json:"created_at"`
	LastModified time.Time `sqlx:"last_modified" json:"last_modified"`
}
