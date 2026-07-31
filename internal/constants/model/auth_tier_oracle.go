package model

import (
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants"
	"time"
)

// AuthTierOracle represents a single active row in Oracle for amount-based auth.
// This replaces the Mongo model's ObjectID with an Oracle hex string ID.
type AuthTierOracle struct {
	ID string `bson:"id" json:"id"`

	Currency  constants.CurrencyType `bson:"currency" json:"currency"`
	MinAmount uint64                 `bson:"min_amount" json:"min_amount"`
	MaxAmount uint64                 `bson:"max_amount" json:"max_amount"`
	Method    constants.Method       `bson:"method" json:"method"`

	Enabled   int `bson:"enabled" json:"enabled"`
	IsDeleted int `bson:"is_deleted" json:"is_deleted"`

	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
	LastModified time.Time `bson:"last_modified" json:"last_modified"`
}
