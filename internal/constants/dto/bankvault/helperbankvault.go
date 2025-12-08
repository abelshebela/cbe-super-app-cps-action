package bankvault

import (
	"cbe-super-app-cps-action/internal/constants"
	"strings"
)

func ToDomainMethod(s string) constants.AccrualMethod {
	switch strings.ToUpper(s) {
	case "SIMPLE":
		return constants.AccrualMethodSimple
	case "COMPOUND":
		return constants.AccrualMethodCompound
	default:
		return constants.AccrualMethod(s)
	}
}
