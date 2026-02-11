package bankvault

import "github.com/shopspring/decimal"

type CreateBankVaultProductRequest struct {
	Name                       string          `json:"name"               example:"Diaspora Fixed Deposit"`
	Interest                   decimal.Decimal `json:"interest"           swaggertype:"string" example:"6"`
	Frequency                  int64           `json:"frequency"          example:"365"`
	LockPeriodDays             string          `json:"lock_period_days"   example:"180"`
	MinAmount                  decimal.Decimal `json:"min_amount"         swaggertype:"string" example:"1000"`
	MaxAmount                  decimal.Decimal `json:"max_amount"         swaggertype:"string" example:"500000"`
	ApplyInterestOnEarlyUnlock bool            `json:"apply_interest_on_early_unlock" swaggertype:"string"`
}

type UpdateBankVaultProductRequest struct {
	MinAmount *float64 `json:"min_amount,omitempty"    swaggertype:"string" example:"2000"`
	MaxAmount *float64 `json:"max_amount,omitempty"    swaggertype:"string" example:"750000"`
}
