package amount_based_auth

import (
	shared_constants "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/constants"
)

// Validate ensures the request has valid data based on the method
func (r UpdateAmountBasedAuthRequest) Validate(method shared_constants.Method) bool {
	switch method {
	case shared_constants.OPEN:
		// OPEN tier only needs MaxAmount
		return r.MaxAmount > 0
	case shared_constants.PIN:
		// PIN tier needs both MinAmount and MaxAmount
		return r.MinAmount > 0 && r.MaxAmount > 0
	case shared_constants.OTPANDPIN:
		// OTP_PIN tier only needs MinAmount
		return r.MinAmount > 0
	default:
		return false
	}
}
