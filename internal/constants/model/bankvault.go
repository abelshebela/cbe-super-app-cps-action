package model

import (
	"cbe-super-app-cps-action/internal/constants"
	"time"

	"github.com/shopspring/decimal"
)

type BankVaultProduct struct {
	ID                         string                     `json:"id" bson:"id"`
	Name                       string                     `json:"name" bson:"name"`
	Currency                   string                     `json:"currency" bson:"currency"`
	RateBps                    decimal.Decimal            `json:"rate_bps" bson:"rate_bps"`
	Method                     constants.AccrualMethod    `json:"method" bson:"method"`
	Frequency                  constants.AccrualFrequency `json:"frequency" bson:"frequency"`
	LockPeriod                 time.Duration              `json:"lock_period" bson:"lock_period"`
	MinAmount                  decimal.Decimal            `json:"min_amount" bson:"min_amount"`
	MaxAmount                  decimal.Decimal            `json:"max_amount" bson:"max_amount"`
	ApplyInterestOnEarlyUnlock bool                       `json:"apply_interest_on_early_unlock" bson:"apply_interest_on_early_unlock"`
	IsActive                   bool                       `json:"is_active" bson:"is_active"`
	CreatedAt                  time.Time                  `json:"created_at" bson:"created_at"`
	UpdatedAt                  time.Time                  `json:"updated_at" bson:"updated_at"`
	DeletedAt                  *time.Time                 `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
	CreatedBy                  string                     `json:"created_by" bson:"created_by"`
	UpdatedBy                  string                     `json:"updated_by" bson:"updated_by"`
	IsDeleted                  bool                       `json:"is_deleted" bson:"is_deleted"`
}

type BankVaultProductPatch struct {
	// Description *string          `json:"description,omitempty" bson:"description,omitempty"`
	MinAmount *decimal.Decimal `json:"min_amount,omitempty" bson:"min_amount,omitempty"`
	MaxAmount *decimal.Decimal `json:"max_amount,omitempty" bson:"max_amount,omitempty"`
	IsActive  *bool            `json:"is_active,omitempty" bson:"is_active,omitempty"`
}

type UpdateBankVault struct {
	// Description *string   `json:"description" bson:"description"`
	MinAmount *float64  `json:"min_amount" bson:"min_amount"`
	MaxAmount *float64  `json:"max_amount" bson:"max_amount"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	UpdatedBy *string   `json:"updated_by" bson:"updated_by"`
}
