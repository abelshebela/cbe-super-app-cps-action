package bankvault

import (
	"cbe-super-app-cps-action/internal/constants"
	"time"

	"github.com/shopspring/decimal"
)

type BankVaultProductResponse struct {
	ID                 string                     `json:"id"`
	Name               string                     `json:"name"`        // immutable
	Description        string                     `json:"description"` // updatable
	Currency           string                     `json:"currency"`    // immutable
	RateBps            decimal.Decimal            `json:"rate_bps"`
	Method             constants.AccrualMethod    `json:"method"`
	Frequency          constants.AccrualFrequency `json:"frequency"`
	LockPeriod         time.Duration              `json:"lock_period"` // stored as duration, request/response use days
	MinAmount          decimal.Decimal            `json:"min_amount"`
	MaxAmount          decimal.Decimal            `json:"max_amount"`
	EarlyUnlockRateBps bool                       `json:"early_unlock_rate_bps"`
	IsActive           bool                       `json:"is_active"`
	IsDeleted          bool                       `json:"is_deleted"`
	CreatedAt          time.Time                  `json:"created_at"`
	UpdatedAt          time.Time                  `json:"updated_at"`
	DeletedAt          *time.Time                 `json:"deleted_at,omitempty"`
}
