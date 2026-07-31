package amount_based_auth

import (
	"strings"

	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants"
)

// CanonicalMethodFromPath maps URL path segments (e.g. OTP_PIN from Swagger) to stored Method values (PIN_OTP).
func CanonicalMethodFromPath(path string) constants.Method {
	switch strings.ToUpper(strings.TrimSpace(path)) {
	case string(constants.OPEN):
		return constants.OPEN
	case string(constants.PIN):
		return constants.PIN
	case "OTP_PIN", "PIN_OTP":
		return constants.OTPANDPIN
	default:
		return constants.Method(strings.TrimSpace(path))
	}
}
