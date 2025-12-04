package amount_based_auth

import "cbe-super-app-cps-action/internal/constants"

// UpdateAmountBasedAuthRequest carries fields for updating any tier type
type UpdateAmountBasedAuthRequest struct {
	MinAmount uint64 `json:"min_amount,omitempty"`
	MaxAmount uint64 `json:"max_amount,omitempty"`
}

// Validate ensures the request has valid data based on the method
func (r UpdateAmountBasedAuthRequest) Validate(method constants.Method) bool {
	switch method {
	case constants.OPEN:
		// OPEN tier only needs MaxAmount
		return r.MaxAmount > 0
	case constants.PIN:
		// PIN tier needs both MinAmount and MaxAmount
		return r.MinAmount > 0 && r.MaxAmount > 0
	case constants.OTPANDPIN:
		// OTP_PIN tier only needs MinAmount
		return r.MinAmount > 0
	default:
		return false
	}
}
