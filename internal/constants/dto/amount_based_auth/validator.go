package amount_based_auth

import (
	"cbe-super-app-cps-action/internal/constants"
	shared_constants "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/constants"
)

// Validate ensures the request has valid data based on the method.
// MinAmount = 0 is allowed here; contextual cascade validation in the service
// layer will reject invalid combinations (e.g. PIN min=0 when OPEN also exists).
func (r UpdateAmountBasedAuthRequest) Validate(method shared_constants.Method) bool {
	cm := CanonicalMethodFromPath(string(method))
	switch cm {
	case constants.OPEN:
		return r.MaxAmount > 0
	case constants.PIN:
		return r.MaxAmount > 0
	case constants.OTPANDPIN:
		return true
	default:
		return false
	}
}
