package model

import (
	"cbe-super-app-cps-action/internal/constants"
	"time"

	"github.com/shopspring/decimal"
)

type BankVaultProduct struct {
	ID                string                     `json:"id"`
	Name              string                     `json:"name"`        // immutable
	Description       string                     `json:"description"` // updatable
	Currency          string                     `json:"currency"`    // immutable
	RateBps           decimal.Decimal            `json:"rate_bps"`
	Method            constants.AccrualMethod    `json:"method"`
	Frequency         constants.AccrualFrequency `json:"frequency"`
	LockPeriod        time.Duration              `json:"lock_period"` // stored as duration, request/response use days
	MinAmount         decimal.Decimal            `json:"min_amount"`
	MaxAmount         decimal.Decimal            `json:"max_amount"`
	EarlyUnlockFeeBps decimal.Decimal            `json:"early_unlock_fee_bps"`
	IsActive          bool                       `json:"is_active"`
	CreatedAt         time.Time                  `json:"created_at"`
	UpdatedAt         time.Time                  `json:"updated_at"`
	DeletedAt         *time.Time                 `json:"deleted_at,omitempty"`
	CreatedBy         string                     `json:"created_by"`
	UpdatedBy         string                     `json:"updated_by"`
	IsDeleted         bool                       `json:"is_deleted"`
}
type BankVaultProductPatch struct {
	Description *string          `json:"description,omitempty"`
	MinAmount   *decimal.Decimal `json:"min_amount,omitempty"`
	MaxAmount   *decimal.Decimal `json:"max_amount,omitempty"`
	IsActive    *bool            `json:"is_active,omitempty"`
}
type UpdateBankVault struct {
	Description *string   `json:"description"`
	MinAmount   *float64  `json:"min_amount"`
	MaxAmount   *float64  `json:"max_amount"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   *string   `json:"updated_by"`
}
