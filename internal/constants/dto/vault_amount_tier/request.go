package vaultamounttier

import "github.com/shopspring/decimal"

type VaultAmountTierRequest struct {
	VaultCategoryID string          `json:"vault_category_id"`
	MinAmount       decimal.Decimal `json:"min_amount"`
	MaxAmount       decimal.Decimal `json:"max_amount"`
	Interest        decimal.Decimal `json:"interest"`
}

type UpdateVaultAmountTierRequest struct {
	MinAmount decimal.Decimal `json:"min_amount"`
	MaxAmount decimal.Decimal `json:"max_amount"`
	Interest  decimal.Decimal `json:"interest"`
}
