package bankvault

import "github.com/shopspring/decimal"

type CreateBankVaultProductRequest struct {
	Name              string          `json:"name"               example:"Diaspora Fixed Deposit"`
	Description       string          `json:"description"        example:"12-month fixed deposit for diaspora customers"`
	Currency          string          `json:"currency"           example:"ETB"`
	RateBps           decimal.Decimal `json:"rate_bps"           swaggertype:"string" example:"450"`
	Method            string          `json:"method"             example:"SIMPLE"`
	Frequency         string          `json:"frequency"          example:"MONTHLY"`
	LockPeriodDays    string          `json:"lock_period_days"   example:"180"`
	MinAmount         decimal.Decimal `json:"min_amount"         swaggertype:"string" example:"1000"`
	MaxAmount         decimal.Decimal `json:"max_amount"         swaggertype:"string" example:"500000"`
	EarlyUnlockFeeBps decimal.Decimal `json:"early_unlock_fee_bps" swaggertype:"string" example:"100"`
}

type UpdateBankVaultProductRequest struct {
	Description *string  `json:"description,omitempty"  example:"Updated marketing description"`
	MinAmount   *float64 `json:"min_amount,omitempty"    swaggertype:"string" example:"2000"`
	MaxAmount   *float64 `json:"max_amount,omitempty"    swaggertype:"string" example:"750000"`
}
