package amount_based_auth

import (
	"cbe-super-app-cps-action/internal/constants"
)

// ----------- Valid methods set -----------

var ValidMethods = map[constants.Method]bool{
	constants.OPEN:      true,
	constants.PIN:       true,
	constants.OTPANDPIN: true,
}

// ----------- TierInput -----------

type TierInput struct {
	Method    constants.Method `json:"method"`
	MinAmount uint64           `json:"min_amount"`
	MaxAmount uint64           `json:"max_amount"`
}

// ----------- AddCurrencyRequest -----------

type AddCurrencyRequest struct {
	Currency constants.CurrencyType `json:"currency"`
	Methods  []constants.Method     `json:"methods"`
	Tiers    []TierInput            `json:"tiers"`
}

func (r AddCurrencyRequest) Validate() bool {
	// if r.Currency != constants.ETB && r.Currency != constants.USD {
	// 	return false
	// }
	if len(r.Methods) == 0 {
		return false
	}
	if len(r.Tiers) != len(r.Methods) {
		return false
	}
	// check each method is valid and no duplicates
	seen := map[constants.Method]bool{}
	for _, m := range r.Methods {
		if !ValidMethods[m] {
			return false
		}
		if seen[m] {
			return false
		}
		seen[m] = true
	}
	// validate tier amounts
	for i, t := range r.Tiers {
		if t.MinAmount == 0 {
			return false
		}
		isLast := i == len(r.Tiers)-1
		if !isLast && t.MaxAmount == 0 {
			return false
		}
		if !isLast && t.MinAmount >= t.MaxAmount {
			return false
		}
	}
	return true
}

// ----------- ResetConfigRequest -----------

type ResetConfigRequest struct {
	Methods []constants.Method `json:"methods"`
	Tiers   []TierInput        `json:"tiers"`
}

func (r ResetConfigRequest) Validate(currency constants.CurrencyType) bool {
	return AddCurrencyRequest{
		Currency: currency,
		Methods:  r.Methods,
		Tiers:    r.Tiers,
	}.Validate()
}

// ----------- CurrencyGroup (grouped response) -----------

type CurrencyGroup struct {
	Currency constants.CurrencyType `json:"currency"`
	Tiers    []TierResponse         `json:"tiers"`
}

type TierResponse struct {
	ID        string           `json:"id"`
	Method    constants.Method `json:"method"`
	MinAmount uint64           `json:"min_amount"`
	MaxAmount uint64           `json:"max_amount"`
	Enabled   bool             `json:"enabled"`
}

// ----------- UpdateAmountBasedAuthRequest (existing) -----------

// UpdateAmountBasedAuthRequest carries fields for updating any tier type
type UpdateAmountBasedAuthRequest struct {
	MinAmount uint64 `json:"min_amount,omitempty"`
	MaxAmount uint64 `json:"max_amount,omitempty"`
}
