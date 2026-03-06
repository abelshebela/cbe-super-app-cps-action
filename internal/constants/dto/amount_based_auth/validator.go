package amount_based_auth

import (
	shared_constants "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/constants"
)

// Validate ensures the request has valid data based on the method.
// MinAmount = 0 is allowed here; contextual cascade validation in the service
// layer will reject invalid combinations (e.g. PIN min=0 when OPEN also exists).
func (r UpdateAmountBasedAuthRequest) Validate(method shared_constants.Method) bool {
	switch method {
	case shared_constants.OPEN:
		// OPEN (free) tier only needs MaxAmount
		return r.MaxAmount > 0
	case shared_constants.PIN:
		// PIN tier needs MaxAmount; MinAmount can be 0 if it's the only method
		return r.MaxAmount > 0
	case shared_constants.OTPANDPIN:
		// OTP_PIN is the last tier; MinAmount validated by cascade in service
		return true
	default:
		return false
	}
}
